const urlParams = new URLSearchParams(window.location.search);

const searchButton = document.getElementById('submit-search')
const searchBoxInput = document.getElementById('user-search-content')

function searchUser(userInput){
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    url.searchParams.set('search', userInput);
    window.location.replace(url.toString());
}


searchButton.addEventListener('click', function() {
    userSearchInput = searchBoxInput.value;
    if(userSearchInput == "") {
        return;
    }
    searchUser(userSearchInput);
})

searchBoxInput.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {  
        userSearchInput = searchBoxInput.value;
        if (userSearchInput == "") {
            return;
        }
        searchUser(userSearchInput);
    }
});

document.getElementById("user-result-list").addEventListener('click', (event) => {
    if(event.target.closest(".userdata-title") == null){
        return;
    }
    var userId = event.target.closest(".userdata-title").id
    if (userId.length == 32) {
        window.location.href = `../user/${userId}`
    }
}
);
