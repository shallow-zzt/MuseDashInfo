const urlParams = new URLSearchParams(window.location.search);

const previousButton = document.getElementById('previous-page')
const nextButton = document.getElementById('next-page')
const totalPage = parseInt(document.getElementById('total-page').textContent)
const nowPage = parseInt(document.getElementById('now-page').textContent)
const pagecounter = document.getElementsByClassName('page-counter')
const pagebox = document.getElementsByClassName('pagebox')


previousButton.addEventListener('click', function() {
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    if (nowPage == 1 || nowPage == null) return;
    url.searchParams.set('page', parseInt(nowPage)- 1);
    window.location.replace(url.toString());   
})

nextButton.addEventListener('click', function() {
    var currentUrl = window.location.href;
    var url = new URL(currentUrl);
    if(nowPage == null) nowPage = 1;
    if (nowPage == totalPage) return;
    url.searchParams.set('page',  parseInt(nowPage) + 1);
    window.location.replace(url.toString());
})

window.addEventListener('load', function() {
    if (parseInt(nowPage) == 1) {
        previousButton.style.borderColor = 'gray';
    }
    if(parseInt(nowPage) == totalPage) {
        nextButton.style.borderColor = 'gray';
    }
    if(nowPage.toString().length + totalPage.toString().length > 4) {
        for (let i = 0; i < pagecounter.length; i++) {
            pagecounter[i].style.fontSize = '16px';
        }
    }
    if(nowPage.toString().length + totalPage.toString().length > 6) {
        for (let i = 0; i < pagebox.length; i++) {
            pagebox[i].style.width = '350px';
        }
    }
})