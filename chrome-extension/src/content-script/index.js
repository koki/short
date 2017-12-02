import githubInjection from 'github-injection';
import ghPageType from 'github-page-type';
import $ from 'jquery';
import _ from 'lodash';
import hljs from 'highlight.js';

githubInjection(window, err => {
    if (!err &&
        ghPageType(window.location.href, ghPageType.REPOSITORY_BLOB) &&
        $('.final-path').text().endsWith('.yaml') &&
        $('#toggle-yaml').length <= 0 &&
        isKubeYaml(getLinesFromRows(getFileRows()))) {
        main();
    }
});

function main() {
    const toggleButton = renderYAMLToggleButton();
    const kubeRows = getFileRows();
    const kubeLines = getLinesFromRows(kubeRows);

    const kokiRows = new Promise((resolve, reject) => {
        sendFileContents(kubeLines, (content) => {
            resolve(buildRowsForContent(content));
        }, reject);
    });

    kokiRows.catch((error) => {
        // TODO: Find a good way to present the error to the user.
        toggleButton.detach();
        console.error(error);
    });

    let toggleDisabled = false;
    toggleButton.click(() => {
        if (toggleDisabled) {
            return;
        }
        if (toggleButton.text() === 'Pretty YAML') {
            toggleButton.text('Working...');
            toggleDisabled = true;
            kokiRows.then((kokiRows) => {
                toggleDisabled = false;
                replaceFileRows(kokiRows);
                toggleButton.text('Kube YAML');
            });
        } else {
            replaceFileRows(kubeRows);
            toggleButton.text('Pretty YAML');
        }
    });
}

function isKubeYaml(lines) {
    let hasApiVersion, hasKind;
    _.forEach(lines, (line) => {
        hasApiVersion |= line.startsWith('apiVersion:');
        hasKind |= line.startsWith('kind:');
    });
    return hasApiVersion && hasKind;
}

function buildRowForLine(lineNumber, line) {
    return `
        <tr>
            <td id="L${lineNumber}" class="blob-num js-line-number" data-line-number="${lineNumber}"></td>
            <td id="LC${lineNumber}" class="blob-code blob-code-inner js-file-line">${line}</td>
        </tr>`;
}

function buildRowsForContent(content) {
    const lines = content.split('\n');

    // Strip the last line if it's empty.
    if (!_.last(lines)) {
        lines.pop();
    }

    let hlState; // Syntax highlighter state, preserved between lines.
    return _.map(lines, (line, lineIndex) => {
        // Highlight each line and build a row for it.
        let hlResult = hljs.highlight('yaml', line, true, hlState);
        hlState = hlResult.top;
        return buildRowForLine(lineIndex + 1, hlResult.value);
    });
}

function getFileRows() {
    return $('.file .blob-wrapper tr');
}

function replaceFileRows(rows) {
    const tbody = $('.file .blob-wrapper table tbody');
    const oldRows = tbody.find('tr');
    oldRows.detach();
    tbody.append(rows);

    return oldRows;
}

function getLinesFromRows(rows) {
    return _.map(rows.find('td.blob-code.blob-code-inner.js-file-line'), (line) => line.textContent);
}

function updateFileWithInfo(fileInfo, fileData) {
    $('#pretty-yaml-info').text(fileInfo);
    $('#pretty-yaml-data').text(fileData);
}

function sendFileContents(lines, onSuccess, onError) {
    const port = chrome.runtime.connect();
    port.postMessage({
        fileLines: lines
    });
    port.onMessage.addListener(function(msg) {
        if (msg.error) {
            onError(msg.error);
            return;
        }

        onSuccess(msg.body);
    });
}

function renderYAMLToggleButton() {
    const rawButton = $('#raw-url');

    const toggleButton = rawButton
        .clone()
        .text('Pretty YAML')
        .attr('href', '#')
        .attr('id', 'toggle-yaml')
        .insertBefore(rawButton);

    return toggleButton;
}
