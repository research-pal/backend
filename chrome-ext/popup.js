

var apiUrl = 'http://research-pal.appspot.com/notes'//'http://research-pal.appspot.com/notes' //'http://localhost:8080/notes' //


function getCurrentTabUrl(callback) { //Question: what does callback hear mean?
  // Query filter to be passed to chrome.tabs.query - see
  // https://developer.chrome.com/extensions/tabs#method-query
  var queryInfo = {
    active: true,
    currentWindow: true
  };

  chrome.tabs.query(queryInfo, function(tabs) {

    var tab = tabs[0];
    var url = tab.url;

    console.assert(typeof url == 'string', 'tab.url should be a string'); //Question: does console.assert() going to work in chrome extensions?

    callback(url);
  });

}

function dosubmit(){
  var notes = document.getElementById('notes').value
  
  getCurrentTabUrl(function(url) {
    putNotesData(encodeURIComponent(url),notes, function(errorMessage) {
      //TODO: need to handle error or send back to ui
    });
  });

}

function putNotesData(url, notes, errorCallback){

var putUrl = apiUrl;
   
  
  var x = new XMLHttpRequest();

  x.open('PUT', putUrl);
  
  x.responseType = 'json';
  x.onload = function() { // Parse and process the response 
    var response = x.response;

    if (!response) { 
      errorCallback('No response from API!');
      return;
    }
    var firstResult = response.Notes;
    document.getElementById('notes').value = firstResult;
    
  };
  x.onerror = function() {
    errorCallback('Network error.');
  };
  var body = '{"URL":"'+url+'","Notes":"'+notes+'"}'
  x.send(body);

}


function getNotesData(searchTerm, errorCallback) {


  var searchUrl = apiUrl + '/'+searchTerm;
	
  var x = new XMLHttpRequest();

  x.open('GET', searchUrl);
  //x.setRequestHeader( 'Access-Control-Allow-Origin', '*'); 
  //x.setRequestHeader( 'Content-Type', 'application/json' );
  
  x.responseType = 'json';
  x.onload = function() { // Parse and process the response 
    var response = x.response;

    if (!response) { 
      errorCallback('No response from API!');
      return;
    }
    var firstResult = response.Notes;
    document.getElementById('notes').value = firstResult;
    
  };
  x.onerror = function() {
    errorCallback('Network error.');
  };
  x.send();
  
}



/*function renderStatus(statusText) {
  document.getElementById('status').textContent = statusText;
}*/

document.addEventListener('DOMContentLoaded', function() { //Question: what is the significance of this pattern of calling function in nested manner in javascript?
  getCurrentTabUrl(function(url) {
    getNotesData(encodeURIComponent(url), function(errorMessage) {
      //TODO: need to handle error
    });
  });

  var postButton = document.getElementById('postButton');
  postButton.addEventListener('click', function() {
    dosubmit();
  }, false); //Question: what is this false doing here?
});
