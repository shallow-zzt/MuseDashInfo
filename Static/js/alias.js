document.getElementById("song-table").addEventListener('click', (event) => {
    var clickedSong = event.target.id
    if(clickedSong != "song-table"){
        window.location.href = `/alias/${clickedSong}`;
    }
});