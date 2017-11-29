chrome.runtime.onConnect.addListener(function(port) {
  //var tab = port.sender.tab;

  port.onMessage.addListener(function(info) {
    var xhr = new XMLHttpRequest();
    var body = info.file
   
    console.log(body)

    xhr.open('POST', 'http://localhost:8080/convert', false)
    xhr.setRequestHeader('Content-Type', "application/yaml")
    xhr.setRequestHeader('Access-Control-Allow-Origin', 'http://localhost:8080/')
    xhr.send(JSON.stringify(body))
    xhr.onload = function() {
      port.postMessage(xhr.ResponseText)
    }
  });
});
