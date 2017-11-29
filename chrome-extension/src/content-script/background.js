chrome.runtime.onConnect.addListener(function(port) {
    //var tab = port.sender.tab;
    function sendRequest(body) {
        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://localhost:8080/convert', true);
        xhr.setRequestHeader('Content-Type', 'application/yaml');
        xhr.setRequestHeader('Access-Control-Allow-Origin', 'http://localhost:8080/');
        xhr.onload = (e) => {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    port.postMessage({
                        body: xhr.responseText
                    });
                } else {
                    port.postMessage({
                        error: {
                            status: xhr.status,
                            statusText: xhr.statusText,
                            text: xhr.responseText
                        }
                    });
                }
            }
        };
        xhr.onerror = (e) => {
            console.error({
                status: xhr.status,
                statusText: xhr.statusText,
                text: xhr.responseText
            });
        };

        xhr.send(body);
    }

    port.onMessage.addListener(function(info) {
        if (info.fileLines) {
            let body = info.fileLines.join('\n');
            sendRequest(body);
        }
    });
});
