const diffRegex = /(\d+)-(\d+)-(\d+)/

function detectDevice() {
    var userAgent = navigator.userAgent;
    if (/Mobile|Android|iP(hone|od|ad)|IEMobile|BlackBerry|Opera Mini/i.test(userAgent)) {
        return 1;
    } else {
        return 0;
    }
}

document.getElementById("song-table").addEventListener('click', (event) => {
    var platform = detectDevice();
    if (diffRegex.test(event.target.id)) {
        diffcode = diffRegex.exec(event.target.id)
        window.location.href = `rank/${diffcode[1]}-${diffcode[2]}?diff=${diffcode[3]-1}&platform=${platform}`
    }
}
);