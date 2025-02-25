const diffSongRegex = /(\d+)-(\d+)/

document.getElementById("song-table").addEventListener('click', (event) => {
    var clickedSong = event.target.id
    if(diffSongRegex.test(clickedSong)){
        window.location.href = `/alias/${clickedSong}`;
    }
});