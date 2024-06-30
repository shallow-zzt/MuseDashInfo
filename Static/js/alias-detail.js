if(document.getElementById('show-alias')){
    const showAlias = document.getElementById('show-alias');
    showAlias.addEventListener('click', function() {
        let displaying = document.getElementById("alias-list").style.display;
    
        if ( displaying == "none" || displaying=="" ){
            document.getElementById("alias-list").style.display = "block";
            showAlias.innerHTML = "[收起]";
        } else {
            document.getElementById("alias-list").style.display = "none";
            showAlias.innerHTML = "[展开]";
        }
    });
}

const returnAliasList = document.getElementById("return-alias-list")
const submitAlias = document.getElementById("submit-alias")
const userAliasInput = document.getElementById("user-alias")


returnAliasList.addEventListener("click",function() {
    history.back();
});

submitAlias.addEventListener("click",function(){
    aliasInput = userAliasInput.value;
    fullSongCode = window.location.href.split('/')[4].split("-");
    albumCode = fullSongCode[0].replace(/[^0-9]/g, "");
    songCode = fullSongCode[1].replace(/[^0-9]/g, "");
    songName = document.getElementById("song-name").innerText;

    fetch('/submit/aliassong/alias', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "album-code": albumCode ,"song-code":songCode,"input-alias": aliasInput})
      })
      .then(response => response.json())
      .then(data => {
        if(data["result"] == 1){
            function recoveryContent() {
                document.getElementById("submit-reminder").innerHTML = "";
                location.reload();
              }
              document.getElementById("submit-reminder").innerHTML = "提交成功!";
              document.getElementById("submit-reminder").style.color = "green";   
              
              setTimeout(recoveryContent, 1000);

        } else {
            function recoveryContent() {
                document.getElementById("submit-reminder").innerHTML = "";
              }
              errorContent = ["该别名已存在","","别名不能为空","别名不能和原名相同","未知错误"]
              document.getElementById("submit-reminder").innerHTML = errorContent[data["result"]];
              document.getElementById("submit-reminder").style.color = "red";   
              
              setTimeout(recoveryContent, 1000);
        }   
    });
});