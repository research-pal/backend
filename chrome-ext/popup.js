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


function getNotesData(searchTerm, errorCallback) {
  
  var searchUrl = 'http://research-pal.appspot.com/notes/' + searchTerm;
	 
	//renderStatus(searchTerm);
  var x = new XMLHttpRequest();

  x.open('GET', searchUrl);
  
  x.responseType = 'json';
  x.onload = function() { // Parse and process the response 
    var response = x.response;

    if (!response || !response.responseData || !response.responseData.results ||
        response.responseData.results.length === 0) {
      errorCallback('No response from Google Image search!');
      return;
    }
    var firstResult = response.responseData.results[0];
    renderStatus(firstResult);
    
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
});
