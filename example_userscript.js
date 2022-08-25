// ==UserScript==
// @name         Example Autop OTP
// @namespace    http://tampermonkey.net/
// @version      0.1
// @description  Solves a grid challange and fills out a custom company portal fields
// @author       You
// @match        <your url>
// @icon
// @grant        GM_xmlhttpRequest
// @require https://www.example.com/some/js/GM_fetch.js
// ==/UserScript==


//adapt to your values before use
var user = "your site login user";
var pass = "your site login pass";
var otpSmsUrl = "https://yoursever.com/messages/?user=dau&pass=1234";

//also adapt the match url above and the selectors ids for user and password below to match your login site 
//all document.getElementById() and document.getElementById() calls.

check4login();
function check4login(){
    console.log("Tampering: check4login");
    var recheck = 0;
    //var d = document.querySelector("login");
    var d = document.getElementById("login")
    if (d){ // found destination, send user and pw
        console.log("Tampering: found id login");
        document.getElementById("login").value = user;
        document.getElementById("passwd").value = pass;
        document.querySelector('#nsg-x1-logon-button').click()
        check4response()
    }
    else { // destination not found, try again in 1 sec
        if (recheck > 5) return;
        recheck++;
        console.log("Tampering: recheck=" + recheck);
        window.setTimeout(check4login, 1000);
    }
}

function check4response(){
    console.log("Tampering: check4otp");
    var recheck = 0;
    var d = document.getElementById("response")
    if (d){ // found destination, send user and pw
        console.log("Tampering: found id response");
        fetchResponse()
    }
    else { // destination not found, try again in 1 sec
        if (recheck > 5) return;
        recheck++;
        console.log("Tampering: recheck=" + recheck);
        window.setTimeout(check4response, 1000);

    }
}

function fetchResponse(){
    console.log("Tampering: fetching response...");
    fetch(otpSmsUrl)
          .then((response) => response.text().then(responseCallback));
}

function responseCallback( retrievedText ) {
    var recheck = 0;
    console.log("Tampering: got response...");
    if (retrievedText) {
        console.log("Tampering: Message" + retrievedText);
        var sms = retrievedText.split("\n");
        if (sms.length > 2) {
            var timeout = parseDate(sms[2]);
            var now = new Date();
            var passIsStillValid = now < timeout;

            console.log("passIsStillValid" + passIsStillValid);
            if (passIsStillValid) {
                let pass = sms[1];
                console.log(pass);
                document.getElementById("response").value = pass;
                document.querySelector('#ns-dialogue-submit').click()
            }
            else {
                if (recheck > 5) return;
                recheck++;
                console.log("Tampering: got old message that is no longer valid");
                window.setTimeout(fetchResponse, 1000);
            }
        }
    }
    else {
        if (recheck > 5) return;
        recheck++;
        console.log("Tampering: got empty response fetch again...");
        window.setTimeout(fetchResponse, 1000);
    }
}

function parseDate(dateText) {
    console.log("Tampering: paring date: " + dateText);
    let position = dateText.search(/\Date:/i);
    dateText = dateText.substring(17,17+29)
	  // format "Thu Aug 04 03:52:27 +0100 2022" tested in Safari.
    // Use new Date(year, month, day, hour, minute) for better browser compatibility.
	  //CEST isn't a valid time zone in ISO 8601,
	  var replacements = {
     "CET": "+0100",
     "CEST": "+0200",
	  };
    for (var key in replacements) dateText = dateText.replace(key, replacements[key]);
    var date = new Date(Date.parse(dateText))
    console.log("Tampering: message valid until: " + date);
    return date;
}
