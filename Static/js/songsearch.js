const urlParams = new URLSearchParams(window.location.search);

const searchButton = document.getElementById('submit-search')
const fliterButton = document.getElementById('submit-fliter')
const searchBoxInput = document.getElementById('user-search-content')
const minDiffFilter = document.getElementById('diff-floor')
const maxDiffFilter = document.getElementById('diff-ceil')


function searchUser(userInput, minDiff, maxDiff){
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    url.searchParams.set('search', userInput);
    if (minDiff != "" && maxDiff != "") {
        url.searchParams.set('minDiff', minDiff);
        url.searchParams.set('maxDiff', maxDiff);
    }
    window.location.replace(url.toString());
}

function getMinDIffValue() {
    if (minDiffFilter.value == "" || parseInt(minDiffFilter.value,10) < 1) {
        return 1;
    }
    return parseInt(minDiffFilter.value,10);
}

function getMaxDIffValue() {
    if (maxDiffFilter.value == "" || parseInt(maxDiffFilter.value,10) > 12) {
        return 12;
    }
    return parseInt(maxDiffFilter.value,10);
}


searchButton.addEventListener('click', function() {
    userSearchInput = searchBoxInput.value;
    minDiffFliterInput = minDiffFilter.value;
    maxDiffFilterInput = maxDiffFilter.value;
    searchUser(userSearchInput, minDiffFliterInput, maxDiffFilterInput);
})

fliterButton.addEventListener('click', function() {
    userSearchInput = searchBoxInput.value;
    minDiffFliterInput = minDiffFilter.value;
    maxDiffFilterInput = maxDiffFilter.value;
    if(userSearchInput == "" &&(minDiffFliterInput == "" || maxDiffFilterInput == "")) {
        return;
    }
    searchUser(userSearchInput, minDiffFliterInput, maxDiffFilterInput);
})


searchBoxInput.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {  
        userSearchInput = searchBoxInput.value;
        minDiffFliterInput = minDiffFilter.value;
        maxDiffFilterInput = maxDiffFilter.value;
        searchUser(userSearchInput, minDiffFliterInput, maxDiffFilterInput);
        searchBoxInput.blur();
    }
});

minDiffFilter.addEventListener('keydown', function(event) {
    if (event.key === 'Enter' || event.key === 'ArrowRight') {  
        maxDiffFilter.focus();
    }
});

maxDiffFilter.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {  
        userSearchInput = searchBoxInput.value;
        minDiffFliterInput = minDiffFilter.value;
        maxDiffFilterInput = maxDiffFilter.value;
        if(userSearchInput == "" &&(minDiffFliterInput == "" || maxDiffFilterInput == "")) {
            return;
        }
        searchUser(userSearchInput, minDiffFliterInput, maxDiffFilterInput);
        maxDiffFilter.blur();
    } else if (event.key === "ArrowLeft") {
        minDiffFilter.focus();
    }
});

minDiffFilter.addEventListener('input', function(event) {
    this.value = this.value.replace(/\D/g, '');
    if (this.value !== "") {
        let num = parseInt(this.value, 10);
        let minDiff = getMinDIffValue();
        let maxDiff = getMaxDIffValue();
        if (num < minDiff) this.value = minDiff;
        if (num > maxDiff) this.value = maxDiff;
    }
});

maxDiffFilter.addEventListener('input', function(event) {
    this.value = this.value.replace(/\D/g, '');
    if (this.value !== "") {
        let num = parseInt(this.value, 10);
        let minDiff = getMinDIffValue();
        let maxDiff = getMaxDIffValue();
        console.log(minDiff, maxDiff,num);
        if (num < minDiff) this.value = minDiff;
        if (num > maxDiff) this.value = maxDiff;
    }
});

window.addEventListener('load', function() {
    let url = new URL(window.location);
    searchBoxInput.value = url.searchParams.get('search');
    minDiffFilter.value = url.searchParams.get('minDiff');
    maxDiffFilter.value = url.searchParams.get('maxDiff');
    url.searchParams.delete('search');
    url.searchParams.delete('minDiff');
    url.searchParams.delete('maxDiff');
    window.history.replaceState(null, '', url);
});

