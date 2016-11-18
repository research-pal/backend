// Copyright (c) 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/**
 * Get the current URL.
 *
 * @param {function(string)} callback - called when the URL of the current tab
 *   is found.
 */
function getCurrentTabUrl(callback) {
  // Query filter to be passed to chrome.tabs.query - see
  // https://developer.chrome.com/extensions/tabs#method-query
  var queryInfo = {
    active: true,
    currentWindow: true
  };

  chrome.tabs.query(queryInfo, function(tabs) {

    var tab = tabs[0];
    var url = tab.url;

    console.assert(typeof url == 'string', 'tab.url should be a string');

    callback(url);
  });

}

function dosubmit(){
  var notes = 'my dummy notes'
  renderStatus(notes);
  getCurrentTabUrl(function(url) {
    //renderStatus(encodeURIComponent(url)+':Performing Google Image search');

    putNotesData(encodeURIComponent(url),notes, function(errorMessage) {
      //renderStatus('Error: ' + errorMessage);
    });
  });

}

function putNotesData(url, notes, errorCallback){

var apiUrl = 'http://localhost:8080/notes'//'http://research-pal.appspot.com/notes/' + url;
   
  
  var x = new XMLHttpRequest();

  x.open('PUT', apiUrl);
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
    //renderStatus(firstResult);
    
  };
  x.onerror = function() {
    errorCallback('Network error.');
  };
  var body = '{"URL":"'+url+'","Notes":"'+notes+'"}'
  x.send(body);

}


function getNotesData(searchTerm, errorCallback) {
  searchTerm= 'www.google.com';

  var searchUrl = 'http://localhost:8080/notes/'+ searchTerm;//'http://research-pal.appspot.com/notes/' + searchTerm;
	
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
    //renderStatus(firstResult);
    
  };
  x.onerror = function() {
    errorCallback('Network error.');
  };
  x.send();
  
}



function renderStatus(statusText) {
  document.getElementById('status').textContent = statusText;
}

document.addEventListener('DOMContentLoaded', function() {
  getCurrentTabUrl(function(url) {
    //renderStatus(encodeURIComponent(url)+':Performing Google Image search');

    getNotesData(encodeURIComponent(url), function(errorMessage) {
      //renderStatus('Error: ' + errorMessage);
    });
  });

   var postButton = document.getElementById('postButton');
   postButton.addEventListener('click', function() {
    dosubmit();
}, false);

});
