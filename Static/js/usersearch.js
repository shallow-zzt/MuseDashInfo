

const searchButton = document.getElementById('submit-search')
const searchBoxInput = document.getElementById('user-search-content')

const statusPic = document.getElementById('status-pic')
const statusTip = [document.getElementById('status-tip-1'), document.getElementById('status-tip-2')]

function isPageRefresh() {
    if (performance.getEntriesByType("navigation")[0]) {
        return performance.getEntriesByType("navigation")[0].type === "reload";
    } else {
        return performance.navigation.type === 1;
    }
}

function searchUser(userInput){
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    url.searchParams.delete('page');
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

window.addEventListener("beforeunload", function () {
    statusPic.innerHTML = '<img src="pic/loading.webp">';
    statusTip[0].innerText = "正在搜索用户";
    statusTip[1].innerText = "请稍等^_^";
});

window.addEventListener("load", function() {
    if (isPageRefresh()){
        let url = new URL(window.location);
        url.searchParams.delete('search');
        url.searchParams.delete('page');
        window.history.replaceState(null, '', url);
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

