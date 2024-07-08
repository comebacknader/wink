// Signup Modal 
var signupModal = document.getElementById('signupModal');

var signupLink = document.getElementById('header-signup-link');

var signupClose = document.getElementById('closeSignupModal');

if (signupLink !== null) {
	signupLink.onclick = function() {
		signupModal.style.display = 'block';
	}	
}

signupClose.onclick = function() {
	signupModal.style.display = 'none';
}

// Login Modal 

var loginModal = document.getElementById('loginModal');

var loginLink = document.getElementById('header-login-link');

var loginClose = document.getElementById('closeLoginModal');

if (loginLink !== null) {
	loginLink.onclick = function() {
		loginModal.style.display = 'block';
	}
}

loginClose.onclick = function() {
	loginModal.style.display = 'none';
}