function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 0.99)) + min;
}

function displayTryFunction(){
    const functionTip = document.getElementById("try-functions");
    functionTipListChilds = functionTip.getElementsByTagName("div");
    functionTipCount = functionTipListChilds.length;

    functionTipListChilds[getRandomInt(0,functionTipCount-1)].style.display = "block";    
    
}

window.addEventListener('load', function() {
    displayTryFunction();
});
