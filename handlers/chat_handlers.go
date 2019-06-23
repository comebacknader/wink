package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/microcosm-cc/bluemonday"
	uuid "github.com/satori/go.uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Message is a message that the user sends.
type Message struct {
	data  []byte
	mtype string
	room  string
}

// CheckMsg converts the JSON sent by user to check mtype.
type CheckMsg struct {
	Mtype  string `json:"mtype,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Amt    int    `json:"amt,omitempty"`
	Sender string `json:"sender, omitempty"`
}

// Client represents the websocket connection.
// The send channel holds the message to be sent to the client's ws.
// The presence channel monitors the Users in the Streamer's room.
// The tip channel deals with coins sent from Users to Streamers.
// The status channel notifies whether Streamer is online/offline.
// The isban channel notifies User if they are banned.
type Client struct {
	id       string
	username string
	ws       *websocket.Conn
	send     chan []byte
	presence chan PresenceMsg
	tip      chan CoinMsg
	status   chan []byte
	isban    chan []byte
}

// subscription is an intermediary between a Client (connection) and a Hub.
type subscription struct {
	conn *Client
	room string
}

// BanCred is the credentials for a banned user.
type BanCred struct {
	username string
	room     string
}

// Hub is a single manager of rooms and the clients within those rooms.
type Hub struct {
	rooms        map[string]map[*Client]bool
	broadcast    chan Message
	addClient    chan subscription
	removeClient chan subscription
	banned       map[string]map[string]bool // BanCred --> true/false
	ban          chan BanCred
	unban        chan BanCred
	sendCoin     chan CoinMsg
	status       chan Message
}

var Hoob = Hub{
	broadcast:    make(chan Message),
	addClient:    make(chan subscription),
	removeClient: make(chan subscription),
	rooms:        make(map[string]map[*Client]bool),
	banned:       make(map[string]map[string]bool),
	ban:          make(chan BanCred),
	unban:        make(chan BanCred),
	sendCoin:     make(chan CoinMsg),
	status:       make(chan Message),
}

type PresenceMsg struct {
	Mtype string   `json:"mtype,omitempty"`
	List  []string `json:"list,omitempty"`
}

type MessageMsg struct {
	Mtype string `json:"mtype, omitempty"`
	Msg   string `json:"msg, omitempty"`
}

// CoinMsg is the message that is sent back to the user.
type CoinMsg struct {
	Mtype string `json:"mtype, omitempty"`
	Msg   string `json:"msg, omitempty"`
	Amt   int    `json:"amt, omitempty"`
	Room  string `json:"room, omitempty"`
}

// Start starts the whole process of listening for added/removed clients and
// broadcasting of messages.
func (hub *Hub) Start() {
	for {
		select {
		case sub := <-Hoob.addClient:
			connections := Hoob.rooms[sub.room]
			if connections == nil {
				connections = make(map[*Client]bool)
				Hoob.rooms[sub.room] = connections
			}
			connections[sub.conn] = true
		case sub := <-Hoob.removeClient:
			connections := Hoob.rooms[sub.room]
			if connections != nil {
				if _, ok := connections[sub.conn]; ok {
					delete(connections, sub.conn)
					close(sub.conn.send)
					if len(connections) == 0 {
						delete(Hoob.rooms, sub.room)
					}
				}
			}
		case message := <-Hoob.broadcast:
			connections := Hoob.rooms[message.room]
			banned := Hoob.banned[message.room]
			if message.mtype == "MSG" {
				// I could save the message here.
				for conn := range connections {
					select {
					case conn.send <- message.data:
					default:
						close(conn.send)
						delete(connections, conn)
						if len(connections) == 0 {
							delete(Hoob.rooms, message.room)
						}
					}
				}
			}
			if message.mtype == "USERS-IN-ROOM" {
				var listOfNames []string
				var streamer *Client
				for conn := range connections {
					if conn.username != message.room {
						if conn.username != "" {
							banned := Hoob.banned[message.room][conn.username]
							if banned != true {
								listOfNames = append(listOfNames, conn.username)
							}
						}
					} else {
						streamer = conn
					}
				}
				list := PresenceMsg{"USERS-IN-ROOM", listOfNames}
				streamer.presence <- list
			}
			if message.mtype == "BANNED-LIST" {
				var bannedUsers []string
				var streamer *Client
				for usr := range banned {
					if usr != message.room {
						if usr != "" {
							bannedUsers = append(bannedUsers, usr)
						}
					}
				}
				list := PresenceMsg{"BANNED-LIST", bannedUsers}
				// Get the streamer client
				for conn := range connections {
					if conn.username == message.room {
						streamer = conn
					}
				}
				streamer.presence <- list
			}
		case bCred := <-Hoob.ban:
			// Ban the user
			banned := Hoob.banned[bCred.room]
			if banned == nil {
				banned = make(map[string]bool)
				Hoob.banned[bCred.room] = banned
			}
			banned[bCred.username] = true
		case bCred := <-Hoob.unban:
			// Unban the user
			delete(Hoob.banned[bCred.room], bCred.username)
		case coinMsg := <-Hoob.sendCoin:
			connections := Hoob.rooms[coinMsg.Room]
			for conn := range connections {
				// If streamer - send integer of how much updated
				select {
				case conn.tip <- coinMsg:
				default:
					close(conn.tip)
					delete(connections, conn)
					if len(connections) == 0 {
						delete(Hoob.rooms, coinMsg.Room)
					}
				}
			}
		case stat := <-Hoob.status:
			connections := Hoob.rooms[stat.room]
			for conn := range connections {
				select {
				case conn.status <- stat.data:
				default:
					close(conn.status)
					delete(connections, conn)
					if len(connections) == 0 {
						delete(Hoob.rooms, stat.room)
					}
				}
			}
		}
	}
}

// write takes a message that is in the Client's send channel and
// 'writes' it to the Client's websocket connection.
func (s *subscription) writePump() {
	ticker := time.NewTicker(pingPeriod)
	c := s.conn
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			jsonMsg := MessageMsg{"MSG", string(message[:])}
			c.writeJ(jsonMsg)
		case listOfUsers, ok := <-c.presence:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			c.ws.WriteJSON(listOfUsers)
		case tipMsg, ok := <-c.tip:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			c.writeJ(tipMsg)
		case stat, ok := <-c.status:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			jsonMsg := MessageMsg{"STATUS", string(stat[:])}
			c.writeJ(jsonMsg)
		case msg, ok := <-c.isban:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			jsonMsg := MessageMsg{"IS-BAN", string(msg[:])}
			c.writeJ(jsonMsg)
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Client) write(messageType int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, message)
}

func (c *Client) writeJ(jsonMsg interface{}) {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	c.ws.WriteJSON(jsonMsg)
}

// read reads the message sent by Client and sends it to Hub.broadcast channel.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		Hoob.removeClient <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var mess CheckMsg
		err := c.ws.ReadJSON(&mess)
		if err != nil {
			//Hoob.removeClient <- c
			//c.ws.Close()
			break
		}
		banned := Hoob.banned[s.room][s.conn.username]
		pol := bluemonday.StrictPolicy()
		if banned != true {
			var m Message
			if mess.Mtype == "MSG" {
				sanMsg := pol.Sanitize(mess.Msg)
				sanSndr := pol.Sanitize(mess.Sender)
				if sanSndr == c.username {
					fullMsg := sanSndr + ": " + sanMsg
					if sanMsg != "" && len(sanMsg) < 150 {
						m = Message{[]byte(fullMsg), "MSG", s.room}
						Hoob.broadcast <- m
					}
				}
			}
			if mess.Mtype == "USERS-IN-ROOM" {
				sanMsg := pol.Sanitize(mess.Msg)
				m = Message{[]byte(sanMsg), "USERS-IN-ROOM", s.room}
				Hoob.broadcast <- m
			}
			if mess.Mtype == "BANNED-LIST" {
				sanMsg := pol.Sanitize(mess.Msg)
				m = Message{[]byte(sanMsg), "BANNED-LIST", s.room}
				Hoob.broadcast <- m
			}
			if mess.Mtype == "BAN" {
				if s.room == s.conn.username {
					usr := BanCred{mess.Msg, s.room}
					Hoob.ban <- usr
				}
			}
			if mess.Mtype == "UNBAN" {
				if s.room == s.conn.username {
					usr := BanCred{mess.Msg, s.room}
					Hoob.unban <- usr
				}
			}
			if mess.Mtype == "SEND-TIP" {
				sanMsg := pol.Sanitize(mess.Msg)
				cm := CoinMsg{"SEND-TIP", sanMsg, mess.Amt, s.room}
				Hoob.sendCoin <- cm
			}
			if mess.Mtype == "STATUS" {
				sanMsg := pol.Sanitize(mess.Msg)
				m = Message{[]byte(sanMsg), "STATUS", s.room}
				Hoob.status <- m
			}
		} else {
			// Client is Banned
			if mess.Mtype == "SEND-TIP" {
				sanMsg := pol.Sanitize(mess.Msg)
				cm := CoinMsg{"SEND-TIP", sanMsg, mess.Amt, s.room}
				Hoob.sendCoin <- cm
			} else {
				noChatMsg := "You have been banned."
				c.isban <- []byte(noChatMsg)
			}
		}
	}
}

// WebSockPage is the handler for websocket connections.
func WebSockPage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	usr, exist := GetCurrentUser(req)
	conn, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		http.NotFound(w, req)
		return
	}

	client := &Client{
		id:       uuid.Must(uuid.NewV4()).String(),
		ws:       conn,
		send:     make(chan []byte),
		presence: make(chan PresenceMsg),
		tip:      make(chan CoinMsg),
		status:   make(chan []byte),
		isban:    make(chan []byte),
	}

	if exist == true {
		client.username = usr.Username
	}

	roomName := req.FormValue("room")

	sub := subscription{client, roomName}
	Hoob.addClient <- sub

	if exist == true {
		go sub.writePump()
		go sub.readPump()
	} else {
		go sub.writePump()
	}

}
