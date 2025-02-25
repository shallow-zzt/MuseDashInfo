const urlParams = new URLSearchParams(window.location.search);
const diff = urlParams.get('diff') || 0;
const platform = urlParams.get('platform') || 0;

const DiffTier=["easy","hard","master","hidden"];
const PlatformList=["pc","mobile"];
const PlatformColor = ["#b2a500","#a16101"];
const DiffColor=["#12af47","#1aacac","#8b28bd","#463d83"];

document.getElementById(DiffTier[diff]).style.backgroundColor = DiffColor[diff];
document.getElementById(DiffTier[diff]).style.color = "white";

document.getElementById(PlatformList[platform]).style.backgroundColor = PlatformColor[platform];
document.getElementById(PlatformList[platform]).style.color = "white";

function changeOption(diff,platform){
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    url.searchParams.set('diff', diff);
    url.searchParams.set('platform', platform);
    window.location.replace(url.toString());
}

function displaySongChartAuthor(){
    const chartAuthorList = document.getElementById("chart-author-list");
    chartAuthorListChilds =chartAuthorList.getElementsByTagName("div");
    chartAuthorCount = chartAuthorListChilds.length;
    if(chartAuthorCount != 1) {
        chartAuthorListChilds[diff].style.display = "block";    
    }else{
        chartAuthorListChilds[0].style.display = "block";    
    }
}

document.getElementById("diff-area").addEventListener('click', (event) => {
    if(event.target.matches("#easy")){
        changeOption(0,platform);
    }
    else if(event.target.matches("#hard")){	
        changeOption(1,platform);
    }
    else if(event.target.matches("#master")){	
        changeOption(2,platform);
    }
    else if(event.target.matches("#hidden")){	
        changeOption(3,platform);
    }
}
);

document.getElementById("platform-area").addEventListener('click', (event) => {
    if(event.target.matches("#pc")){
        changeOption(diff,0);
    }
    else if(event.target.matches("#mobile")){	
        changeOption(diff,1);	
    }
}
);

document.getElementById("rank-table").addEventListener('click', (event) => {
    if(event.target.id.length == 32){
        window.location.href = `../user/${event.target.id}`
    }
}
);

document.addEventListener("DOMContentLoaded", function () {
    const input = document.getElementById("rank-user-search");
    const tableRows = document.querySelectorAll("tbody tr");

    input.addEventListener("input", function () {
        const filterText = input.value.toLowerCase();

        tableRows.forEach(row => {
            const usernameCell = row.querySelector("td:nth-child(2) .user-underline");
            if (usernameCell) {
                const username = usernameCell.textContent.toLowerCase();
                if (username.includes(filterText)) {
                    row.style.display = "";
                } else {
                    row.style.display = "none";
                }
            }
        });
    });
});

window.addEventListener('load', function() {
    displaySongChartAuthor();
});


