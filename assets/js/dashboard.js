var editBtn = document.querySelector('#dash-edit-strm-btn');
var strmInfo = document.querySelector('#dash-stream-info');
var editStrmInfo = document.querySelector('#dash-edit-stream-info');
var cancelEditBtn = document.querySelector('#dash-cancel-edit');
var updateStreamForm = document.querySelector('#dash-update-stream-form');

var title = document.getElementById('dash-stream-title');
var game = document.getElementById('dash-stream-game');
var twitter = document.getElementById('dash-stream-twitter');
var siteone = document.getElementById('dash-stream-siteone');
var twitterBox = document.querySelector('#twitter-present');
var siteoneBox = document.querySelector('#siteone-present');

// Error message that appears when invalid info submitted
var streamErrs = document.getElementById('dash-stream-errors');

var msgBox = document.querySelector(".chatbox-msg-box");


editBtn.onclick = function() {
	strmInfo.style.display = "none";
	editStrmInfo.style.display = "block";
	streamErrs.classList.remove("alert");
	streamErrs.classList.remove("alert-danger");
	streamErrs.textContent = "";	
};

cancelEditBtn.onclick = function() {
	strmInfo.style.display = "block";
	editStrmInfo.style.display = "none";	
};

// Check whether to show twitter logo/text. 
function showTwitter() {
	var text = twitter.textContent;
	if (text ===  "" || text === 'Enter Twitter') {
		twitterBox.style.display = "none"
	} else {
		twitterBox.style.display = "flex";
	}
};

// Check whether to show personal website (siteone) logo/text.
function showSiteOne() {
	var text = siteone.textContent;
	if (text === "" || text == 'Enter Personal Website'){
		siteoneBox.style.display = "none";
	} else {
		siteoneBox.style.display = "flex";
	}
}

showTwitter();
showSiteOne();

// Mobile View 

var viewChatBtn = document.getElementById('viewChatBtn');
var viewVidBtn = document.getElementById('viewVidBtn');
var chatbox = document.getElementById('chatbox');
var vidCol = document.getElementById('video-column');

viewChatBtn.onclick = function() {
	chatbox.style.display = 'flex';
	vidCol.style.display = 'none';
}

viewVidBtn.onclick = function() {
	chatbox.style.display = 'none';
	vidCol.style.display = 'flex';	
}

window.addEventListener("resize", function() {
	var cliWidth = document.documentElement.clientWidth;
	if (cliWidth > 600 && cliWidth < 900) {
		chatbox.style.display = 'none';
		vidCol.style.display = 'flex';
	}
	if (cliWidth > 900) {
		chatbox.style.display = 'flex';
	}
});


// VideoJS Player 
var streamerName = document.getElementById("get-curr-streamer").textContent;
var srcUrl = 'https://winkgg.com/hls/' + streamerName + '.m3u8';
var player = videojs('video-player');
player.play();
// Order of errText below 'player' is important
var errText = document.querySelector('.vjs-modal-dialog-content');
player.on('error', function(e) {
	errText.textContent = "Streamer is offline. Try refreshing your " +  
	"browser around the time of their stream";
});

updateStreamForm.addEventListener('submit', updateStreamRequest);

function updateStreamRequest(e) {
	e.preventDefault();
	var streamTitle = document.getElementById('dash-edit-stream-title').value;
	var streamGame = document.getElementById('dash-edit-stream-game').value;
	var streamTwit = document.getElementById('dash-edit-stream-twitter').value;
	var streamSiteOne = document.getElementById('dash-edit-stream-siteone').value;

	xhr = new XMLHttpRequest();
	if (!xhr) {
		alert('Cannot make request at this time.');
		return false;
	}

	xhr.onreadystatechange = function() {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			if (xhr.status === 200) {
				strmInfo.style.display = "block";
				editStrmInfo.style.display = "none";
				var response = JSON.parse(xhr.responseText);
				title.textContent = response.title;
				game.textContent = response.game;
				twitter.textContent = response.twitter;
				showTwitter();
				siteone.textContent = response.siteone;
				showSiteOne(); 
			} else {
				streamErrs.classList.add("alert");
				streamErrs.classList.remove("alert-danger");				
				streamErrs.classList.add("alert");
				streamErrs.classList.add("alert-danger");
				streamErrs.textContent = xhr.responseText;				
			}
		}
	}

	xhr.open('POST', '/updatestream', true);
	xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhr.send("title="+streamTitle+"&game="+streamGame+"&twitter="+streamTwit+"&siteone="+streamSiteOne);
};

var updateStatusForm = document.querySelector('#dash-online-form');

updateStatusForm.addEventListener('submit', updateStatusReq);

function updateStatusReq(e) {
	e.preventDefault();

	xhr = new XMLHttpRequest();
	if (!xhr) {
		alert('Cannot make request at this time.');
		return false;
	}

	var onlineInput = document.querySelector('#dash-input-online');
	var onlineStatus = onlineInput.value;
	var statusBtn = document.querySelector('#dash-change-status-btn');
	var statusText = document.querySelector('#dash-online-status');
	var statusCirc = document.querySelector('#dash-status-circ');

	xhr.onreadystatechange = function() {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			if (xhr.status === 200) {
				var jsonRes = JSON.parse(xhr.responseText);
				if (jsonRes.status === "offline") {
					onlineInput.value = "online";
					statusBtn.classList.remove('btn-outline-red');
					statusBtn.classList.add('form-outline-green');
					statusBtn.textContent = "Go Online";
					statusCirc.classList.remove('dash-online-circ');
					statusCirc.classList.add('dash-offline-circ');
					statusText.textContent = "Offline";
					socket.send('{"mtype":"STATUS","msg":"offline"}');
				}
				if (jsonRes.status === "online") {
					onlineInput.value = "offline";
					statusBtn.classList.remove('form-outline-green');
					statusBtn.classList.add('btn-outline-red');
					statusBtn.textContent = "Go Offline";
					statusCirc.classList.remove('dash-offline-circ');
					statusCirc.classList.add('dash-online-circ');
					statusText.textContent = "Online";
					socket.send('{"mtype":"STATUS","msg":"online"}');					
				} 
			}
		}
	};

	xhr.open('POST', '/updateonline', true);
	xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhr.send("online="+onlineStatus);
};

// CHATBOX 
var socket;
initChat();
function initChat(){
	var userName = streamerName;
	if (window.location.hostname === 'localhost') {
		socket = new WebSocket("wss://localhost:10443/ws" + '?' + 'room=' + streamerName);		
	} else {
		socket = new WebSocket("wss://www.wink.gg:443/ws" + '?' + 'room=' + streamerName);		
	}

	var sendMsgBtn = document.getElementById("sendMsgBtn");
	var chatTextArea = document.getElementById("chatbox-msg-textarea");

	// Chatbox Navbar 

	var burgerBtn = document.getElementById("dash-chat-menu-btn");
	var usersInRoomMenu = document.getElementById("dash-chat-dropdown-menu"); 
	var usersInRoomLink = document.getElementById("dash-chat-users-item");
	var usersInRoomBox = document.getElementById("chatbox-users-room-box");
	var chatRoomLink = document.getElementById("dash-chatroom-item");
	var chatUsersUL = document.getElementById("chatbox-users-ul");
	var bannedUsersLink = document.getElementById("dash-chat-banned-item");
	var bannedUserBox = document.getElementById("dash-banned-box");
	var bannedUserUL = document.getElementById("chatbox-banned-ul");

	usersInRoomMenu.style.display = 'none';
	burgerBtn.onclick = function() {
		if (usersInRoomMenu.style.display === 'none') {
			usersInRoomMenu.style.display = 'flex';
		} else {
			usersInRoomMenu.style.display = 'none';		
		}
	}

	// Navigating to the Users In Room 
	// Upon clicking, should get list of users. 
	usersInRoomLink.onclick = function() {
		chatUsersUL.innerHTML = "";		
		usersInRoomMenu.style.display = 'none';
		msgBox.style.display = 'none';
		bannedUserBox.style.display = 'none';	
		usersInRoomBox.style.display = 'block';
		var jsonMsg = {"mtype":"USERS-IN-ROOM", "msg":""};
		socket.send(JSON.stringify(jsonMsg));
		bannedUserUL.innerHTML = "";		
	}

	// Navigating (usually back) to the Chatroom
	chatRoomLink.onclick = function() {
		bannedUserBox.style.display = 'none';	
		usersInRoomMenu.style.display = 'none';
		usersInRoomBox.style.display = 'none';
		msgBox.style.display = 'block';
		// Delete all the <li>'s in UsersInRoom list
		chatUsersUL.innerHTML = "";	
		bannedUserUL.innerHTML = ""; 
	}

	// Navigate to the Banned Users List
	bannedUsersLink.onclick = function() {
		bannedUserUL.innerHTML = ""; 	
		usersInRoomMenu.style.display = 'none';
		usersInRoomBox.style.display = 'none';
		msgBox.style.display = 'none';
		bannedUserBox.style.display = 'block';
		var jsonMsg = {"mtype":"BANNED-LIST", "msg":""};
		socket.send(JSON.stringify(jsonMsg));
		chatUsersUL.innerHTML = "";			
	}

	// Chatbox -- Submitting Message

	chatTextArea.addEventListener("keypress", function(event) {
		if (event.keyCode == 13) {
			event.preventDefault();
			var msgText = document.getElementById("chatbox-msg-textarea").value;
			if (msgText === '') {
				return;
			}
			if (msgText.length > 150) {
				return;
			}		
			var jsonMsg = {"mtype":"MSG", "sender": userName, "msg":msgText};
			socket.send(JSON.stringify(jsonMsg));
			document.getElementById("chatbox-msg-textarea").value = "";		
		} 
	});

	sendMsgBtn.addEventListener("click", function(event) {
		var msgText = document.getElementById("chatbox-msg-textarea").value;
		if (msgText === '') {
			return;
		}		
		if (msgText.length > 150) {
			return;
		}		
		var jsonMsg = {"mtype":"MSG", "sender": userName, "msg":msgText};
		socket.send(JSON.stringify(jsonMsg));	
		document.getElementById("chatbox-msg-textarea").value = "";		
	});

	socket.onopen = function(event) {
		var li = document.createElement("li");
		li.innerHTML = "Connected to Chatroom!";
		li.className = "list-group-item chatbox-li-item chatbox-connect-msg";
		var ul = document.getElementById("chatbox-ul");
		ul.appendChild(li);
		msgBox.scrollTop = msgBox.scrollHeight;			
	};

	socket.onmessage = function(event) {
		// Need to check what type of message it is 'msg' or 'presence'
		var jsonData = JSON.parse(event.data);
		if (jsonData.mtype === 'MSG') {
			var li = document.createElement("li");
			li.innerHTML = jsonData.msg;
			li.className = "list-group-item chatbox-li-item";
			var ul = document.getElementById("chatbox-ul");
			ul.appendChild(li);
			msgBox.scrollTop = msgBox.scrollHeight;		
		}
		if (jsonData.mtype === 'SEND-TIP') {
			var li = document.createElement("li");
			li.innerHTML = jsonData.msg;
			li.className = "list-group-item chatbox-li-item chat-send-tip-msg";
			var ul = document.getElementById("chatbox-ul");
			ul.appendChild(li);
			msgBox.scrollTop = msgBox.scrollHeight;
			// Update the Total # of Coins
			var sentAmt = jsonData.amt;
			var totalAmt = document.getElementById('dash-total-coin-amt');
			var totalAmtNum = parseInt(totalAmt.textContent);
			totalAmtNum += sentAmt;
			totalAmt.textContent = totalAmtNum + '';

		}		
		if (jsonData.mtype === "USERS-IN-ROOM") {
			var inRoomArray = jsonData.list;
			if (typeof inRoomArray === 'undefined') {
				var li = document.createElement("li");
				li.innerHTML = "No users in room.";
				li.className = "list-group-item chatbox-li-item";
				chatUsersUL.appendChild(li);			
				return;
			}		
			for (var i = 0; i < inRoomArray.length; i++) {
				var li = document.createElement("li");
				li.innerHTML = inRoomArray[i];
				li.className = "list-group-item chatbox-li-item";
				chatUsersUL.appendChild(li);
				var img = document.createElement("img");
				img.className = "dash-ban-btn";
				img.src = "/assets/img/ban.svg";
				li.appendChild(img);
				img.onclick = function() {
					// Fire off message to Ban this user
					var jsonMsg = {"mtype":"BAN", "msg":this.parentElement.textContent};
					socket.send(JSON.stringify(jsonMsg));
					chatUsersUL.removeChild(this.parentElement);
				}					
			}
		}
		if (jsonData.mtype === "BANNED-LIST") {
			var bannedList = jsonData.list;
			if (typeof bannedList === 'undefined') {
				var li = document.createElement("li");
				li.innerHTML = "No banned users.";
				li.className = "list-group-item chatbox-li-item";
				bannedUserUL.appendChild(li);			
				return;
			}
			var listLength = bannedList.length;
			for (var i = 0; i < listLength; i++) {
				var li = document.createElement("li");
				li.innerHTML = bannedList[i];
				li.className = "list-group-item chatbox-li-item";
				bannedUserUL.appendChild(li);
				var img = document.createElement("img");
				img.className = "dash-unban-btn";
				img.src = "/assets/img/unban.svg";
				li.appendChild(img);
				img.onclick = function() {
					// Fire off message to Ban this user
					var jsonMsg = {"mtype":"UNBAN", "msg":this.parentElement.textContent};
					socket.send(JSON.stringify(jsonMsg));
					bannedUserUL.removeChild(this.parentElement);
				}				
			}
		} 
	}

	var chatConnBtn = document.querySelector('#chatbox-connect-btn');

	socket.onclose = function(event) {
		// Display message as an LI with it's own style - color red.
		var li = document.createElement("li");
		li.innerHTML = "DISCONNECTED FROM ROOM";
		li.className = "list-group-item chatbox-li-item chatbox-disconnect-msg";
		var ul = document.getElementById("chatbox-ul");
		ul.appendChild(li);
		li = document.createElement("li");
		li.innerHTML = "Reconnecting...";
		li.className = "list-group-item chatbox-li-item";
		ul = document.getElementById("chatbox-ul");
		ul.appendChild(li);
		msgBox.scrollTop = msgBox.scrollHeight;		
		setTimeout(function(){
			var last = ul.lastElementChild;
			ul.removeChild(last);
			initChat()
		}, 3000);	 
	}

	window.onbeforeunload = function(event){
		socket.close();
	};
};

