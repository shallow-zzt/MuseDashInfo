
const diffUserRegex = /(\d+)-(\d+)-(\d+)-(\d+)/

const PlatformList=["pc","mobile"];

document.getElementById("user-song-table").addEventListener('click', (event) => {
    if(event.target.closest(".song-data-frame") == null){
        return;
    }
    var songInfo = event.target.closest(".song-data-frame").id
    if (diffUserRegex.test(songInfo)) {
        diffcode = diffUserRegex.exec(songInfo)
        window.location.href = `../rank/${diffcode[1]}-${diffcode[2]}?diff=${diffcode[3]-1}&platform=${diffcode[4]}`
    }
}
);