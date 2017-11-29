import githubInjection from 'github-injection';
import ghPageType from 'github-page-type';
import $ from "jquery";

githubInjection(window, err => {
    !err
    && ghPageType(window.location.href, ghPageType.REPOSITORY_BLOB)
    && $(".final-path").text().endsWith(".yaml")
    && $("#toggle-yaml").length <= 0
    && renderYAMLToggleButton();
});

chrome.runtime.onConnect.addListener(function(port) {
    port.onMessage.addListener(function(msg) {
        updateFileWithInfo(msg.file-info, msg.file-data);
    });
});

function updateFileWithInfo(fileInfo, fileData) {
    $("#pretty-yaml-info").text(fileInfo);
    $("#pretty-yaml-data").text(fileData);
}

function renderYAMLToggleButton() {
    let fileElement = $(".file");
    let prettyYAMLFileInfo = $("<div />", {id: "pretty-yaml-info"}).hide().appendTo(fileElement.find(".file-header"));
    let prettyYAMLFileData = $("<table />", {id: "pretty-yaml-data"}).hide().appendTo(fileElement.find(".blob-wrapper"));
    let codeElementFileInfo = fileElement.find(".file-info");
    let codeElementFileData = fileElement.find(".blob-wrapper");
    let rawButton = fileElement.find("#raw-url");
    let toggleButton = rawButton
        .clone()
        .text("Pretty YAML")
        .attr("href", "#")
        .attr("id", "toggle-yaml")
        .insertBefore(rawButton)
        .click(() => {
            if (toggleButton.text() === "Pretty YAML") {
                codeElementFileInfo.hide();
                codeElementFileData.hide();
                prettyYAMLFileInfo.show();
                prettyYAMLFileData.show();
                chrome.runtime.connect().
                    postMessage(codeElementFileData.text());
                toggleButton.text("Kube YAML");
            } else {
                prettyYAMLFileInfo.hide();
                prettyYAMLFileData.hide();
                codeElementFileInfo.show();
                codeElementFileData.show();
                toggleButton.text("Pretty YAML");
            }
        });
}
