const serverURL = 'KOKI_SHORT_SERVER_URL';

function sendRequest(body, onSuccess, onError) {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', `${serverURL}/convert`, true);
    xhr.setRequestHeader('Content-Type', 'application/yaml');
    xhr.setRequestHeader('Access-Control-Allow-Origin', `${serverURL}/`);
    xhr.onload = (e) => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                onSuccess(xhr.responseText);
            } else {
                onError({
                    status: xhr.status,
                    statusText: xhr.statusText,
                    text: xhr.responseText
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

chrome.runtime.onConnect.addListener(function(port) {
    //let tab = port.sender.tab;

    function sendSuccessResponse(body) {
        port.postMessage({
            body: body
        });
    }

    function sendErrorResponse(error) {
        port.postMessage({
            error: error
        });
    }

    port.onMessage.addListener(function(info) {
        if (info.fileLines) {
            doConversion(info.fileLines);
        }
    });

    let attemptedSignIn = false; // TODO: this is messy
    function doConversion(fileLines) {
        checkCookies(doSendRequest, () => githubSignIn(doSendRequest));

        function doSendRequest() {
            sendRequest(fileLines.join('\n'), sendSuccessResponse, (error) => {
                if (error.status == 401 && !attemptedSignIn) {
                    // Attempt to sign in and try again.
                    githubSignIn(doSendRequest);
                } else {
                    sendErrorResponse(error);
                }
            });
        }
    }

    function checkCookies(onSuccess, onFailure) {
        chrome.cookies.get({
            url: serverURL,
            name: 'user'
        }, (result) => {
            console.log(result);
            if (result) {
                onSuccess();
            } else {
                onFailure && onFailure();
            }
        });
    }

    function githubSignIn(onSuccess) {
        attemptedSignIn = true;
        const loginWindow = window.open(`${serverURL}/login`);
        const checkWindowClosed = setInterval(() => {
            if (!loginWindow || loginWindow.closed) {
                clearInterval(checkWindowClosed);
                checkCookies(onSuccess);
            }
        }, 200);
    }

});
