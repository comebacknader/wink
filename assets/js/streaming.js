var twitter = document.getElementById('dash-stream-twitter');
var siteone = document.getElementById('dash-stream-siteone');
var twitterBox = document.querySelector('#twitter-present');
var siteoneBox = document.querySelector('#siteone-present');

// Error message that appears when invalid info submitted
var streamErrs = document.getElementById('dash-stream-errors');

var msgBox = document.querySelector(".chatbox-msg-box");


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

// VideoJS Player
var streamerName = document.getElementById("streaming-receiver-input").value;
var srcUrl = 'https://winkgg.com/hls/' + streamerName + '.m3u8';
var player = videojs('video-player');
player.play();
// Order of errText below 'player' is important
var errText = document.querySelector('.vjs-modal-dialog-content');
player.on('error', function(e) {
	errText.textContent = "Streamer is offline. Try refreshing your " +  
	"browser around the time of their stream";
});


// CHATBOX

// Sending a Tip 
var sendTipBtn = document.querySelector('#sendTipBtn');
var amtSelAmt = document.getElementById('chat-coin-sel-amt');
var chatErrMsg = document.getElementById('chat-err-msg');

sendTipBtn.onclick = function() {
	// Get values of inputs
	var senderVal = document.querySelector('#streaming-sender-input').value;
	var receiverVal = document.querySelector('#streaming-receiver-input').value;
	var tipVal = parseInt(amtSelAmt.textContent);

	if (tipVal === 0) {
		return;
	}

	xhr = new XMLHttpRequest();
	if (!xhr) {
		alert('Cannot make request at this time.');
		return false;
	}
	xhr.onreadystatechange = function() {
		if (xhr.readyState === XMLHttpRequest.DONE) {
			if (xhr.status === 200) {
				var jsonRes = JSON.parse(xhr.responseText);
				// Send message on socket.
				var coins = 'coins';
				var tipAmt = jsonRes.amt; 
				if (tipAmt == 1) {
					coins = 'coin';
				}
				var tipMsg = jsonRes.sender + ' gave ' + tipAmt + ' ' + coins + '!';
				var jsonTipMsg = {"mtype":"SEND-TIP", "msg":tipMsg, "amt": tipAmt};  
				socket.send(JSON.stringify(jsonTipMsg)); 
				amtSelAmt.textContent = '0';
			} else {
				chatErrMsg.style.display = 'block';							
				setTimeout(function() {chatErrMsg.style.display = 'none';}, 3000);					
			}
		}
	};

	xhr.open('POST', '/tip', true);
	xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhr.send("sender="+senderVal+"&receiver="+receiverVal+"&amount="+tipVal);
};

// Initialize Socket

var socket; 
initChat();

function initChat() {
	var userName = document.getElementById("streaming-sender-input").value;
	if (window.location.hostname === 'localhost') {
		socket = new WebSocket("wss://localhost:10443/ws" + '?' + 'room=' + streamerName);		
	} else {
		socket = new WebSocket("wss://www.wink.gg/ws" + '?' + 'room=' + streamerName);		
	}


	var sendMsgBtn = document.getElementById("sendMsgBtn");
	var chatTextArea = document.getElementById("chatbox-msg-textarea");
	var chatBanMsg = document.getElementById("chatbox-banned-msg");

	var statusText = document.querySelector('#dash-online-status');
	var statusCirc = document.querySelector('#dash-status-circ');


	// User presses enter on chatbox
	var signupModal = document.getElementById('signupModal');
	var signupClose = document.getElementById('closeSignupModal');
	signupClose.onclick = function() {
		signupModal.style.display = 'none';
	}

	// Check if logged in - if not, display signup modal. 
	var loggedIn = document.getElementById("header-logout-form");
	chatTextArea.addEventListener("keypress", function(event) {
		if (event.keyCode == 13) {
			event.preventDefault();
			if (loggedIn === null) {
				signupModal.style.display = 'block';
			} else {
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
		} 
	});

	sendMsgBtn.addEventListener("click", function(event) {
		if (loggedIn === null) {
			signupModal.style.display = 'block';
		} else {
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

	socket.onopen = function(event) {
		var li = document.createElement("li");
		li.innerHTML = "Connected to Chatroom!";
		li.className = "list-group-item chatbox-li-item chatbox-connect-msg";
		var ul = document.getElementById("chatbox-ul");
		ul.appendChild(li);
		msgBox.scrollTop = msgBox.scrollHeight;	
	};

	socket.onmessage = function(event) {
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
		}
		if (jsonData.mtype === 'STATUS') {
			if (jsonData.msg === 'offline') {
				statusCirc.classList.remove('dash-online-circ');
				statusCirc.classList.add('dash-offline-circ');
				statusText.textContent = "Offline";			
			}
			if (jsonData.msg === 'online') {
				statusCirc.classList.remove('dash-offline-circ');
				statusCirc.classList.add('dash-online-circ');
				statusText.textContent = "Online";			
			}
		}
		if (jsonData.mtype === 'IS-BAN') {
			chatBanMsg.style.display = 'block';
			setTimeout(function() {
				chatBanMsg.style.display = 'none';
			}, 5000);
		}	
	}

	socket.onclose = function(event) {
		// Display message as an LI with it's own style - color red.
		var li = document.createElement("li");
		li.innerHTML = "DISCONNECTED FROM CHAT!";
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

}

window.onbeforeunload = function(event){
	socket.close();
};

var sendCoinBtn = document.getElementById("sendCoinBtn");
var sendMsgBox = document.getElementById("chatbox-send-chat-box");
var sendCoinBox = document.getElementById("chatbox-send-coins-box");


// Go to the Send Coin panel 
sendCoinBtn.onclick = function() {
	sendMsgBox.style.display = 'none';
	sendCoinBox.style.display = 'block';
}

// onclick listeners on the different amount of coins
var plusOne = document.getElementById('coin-one-btn');
var plusFive = document.getElementById('coin-five-btn');
var plusTen = document.getElementById('coin-ten-btn');
var plusFifty = document.getElementById('coin-fitty-btn');
var plusHundred = document.getElementById('coin-hundo-btn');

var totalCoin = document.getElementById('chat-coin-tot-amt');

plusOne.onclick = function() {
	// See if they have 1 or more
	var totalNum = parseInt(totalCoin.textContent);
	if (totalNum < 1) {
		return;
	}
	// Subtract from total 
	totalNum = totalNum - 1;
	totalCoin.textContent = totalNum.toString();
	// Increment by 1 
	var getNum = parseInt(amtSelAmt.textContent);
	getNum += 1;
	amtSelAmt.textContent = getNum.toString();
}

plusFive.onclick = function() {
	// See if they have 5 or more
	var totalNum = parseInt(totalCoin.textContent);
	if (totalNum < 5) {
		return;
	}
	// Subtract from total 
	totalNum = totalNum - 5;
	totalCoin.textContent = totalNum.toString();
	// Increment by 5 	
	var getNum = parseInt(amtSelAmt.textContent);
	getNum += 5;
	amtSelAmt.textContent = getNum.toString();
}

plusTen.onclick = function() {
	// See if they have 10 or more
	var totalNum = parseInt(totalCoin.textContent);
	if (totalNum < 10) {
		return;
	}
	// Subtract from total 
	totalNum = totalNum - 10;
	totalCoin.textContent = totalNum.toString();
	// Increment by 10 	
	var getNum = parseInt(amtSelAmt.textContent);
	getNum += 10;
	amtSelAmt.textContent = getNum.toString();
}

plusFifty.onclick = function() {
	// See if they have 50 or more
	var totalNum = parseInt(totalCoin.textContent);
	if (totalNum < 50) {
		return;
	}
	// Subtract from total 
	totalNum = totalNum - 50;
	totalCoin.textContent = totalNum.toString();
	// Increment by 50 	
	var getNum = parseInt(amtSelAmt.textContent);
	getNum += 50;
	amtSelAmt.textContent = getNum.toString();
}

plusHundred.onclick = function() {
	// See if they have 100 or more
	var totalNum = parseInt(totalCoin.textContent);
	if (totalNum < 100) {
		return;
	}
	// Subtract from total 
	totalNum = totalNum - 100;
	totalCoin.textContent = totalNum.toString();
	// Increment by 100 	
	var getNum = parseInt(amtSelAmt.textContent);
	getNum += 100;
	amtSelAmt.textContent = getNum.toString();
}

var backToChatBtn = document.getElementById('backToChatBtn');

backToChatBtn.onclick = function() {
	sendMsgBox.style.display = 'block';
	sendCoinBox.style.display = 'none';
	var amtToReset = parseInt(amtSelAmt.textContent);
	amtSelAmt.textContent = '0';
	totalNum = parseInt(totalCoin.textContent);
	totalNum += amtToReset;
	totalCoin.textContent = totalNum.toString(); 
}


