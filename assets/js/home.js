var becomeBtn = document.querySelector("#home-become-btn");
var becomeInfo = document.querySelector("#home-become-info");

becomeBtn.onclick = function() {
	this.style.display = "none";
	becomeInfo.style.display = "block";
}


var eighteenModal = document.getElementById('eighteenModal');
var eighteenEnterBtn = document.getElementById('eighteenEnterBtn');

var eighteen = localStorage.getItem('eighteen');
if (eighteen === null) {
	eighteenModal.style.display = 'block';
}

eighteenEnterBtn.onclick = function() {
	eighteenModal.style.display = 'none';
	localStorage.setItem('eighteen', 'agreed');
}

var player = videojs('video-player');
player.play();
// Order of errText below 'player' is important
var errText = document.querySelector('.vjs-modal-dialog-content');
player.on('error', function(e) {
	errText.textContent = "Streamer is offline. Try refreshing your " +  
	"browser around the time of their stream";
});