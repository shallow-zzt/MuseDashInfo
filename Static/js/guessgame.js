const standardAnswer = document.getElementById('answer-song');
const userAnswer = document.getElementById('answer');

const submitAnswer = document.getElementById('submit-answer');
const resetQuest = document.getElementById('reset-quest');

function setResult(content){
  document.getElementById("answer").value = content;
}

submitAnswer.addEventListener('click', function() {
    var answer = userAnswer.value;
    fetch('/submit/guessgame/answer', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ "answer": answer ,"standard-answer":standardAnswer.innerText})
      })
      .then(response => response.json())
      .then(data => {
        let result = data["result"];
        let possibleResult = data["possible-result"];
        document.getElementById("possible-answer").style.display="none";
        document.getElementById("possible-list").innerHTML = "";
        if(parseInt(result)==0){
            document.getElementById("answer-result").innerHTML = "回答错误";
        } else if(parseInt(result)==1){
            document.getElementById("answer-result").innerHTML = "回答正确";
            document.getElementById("quest-pic").style.display = "none";
            document.getElementById("answer-pic").style.display = "block";
            document.getElementById("submit-answer").style.display = "none";
        } else if(parseInt(result)==2){
          document.getElementById("possible-answer").style.display="block";
            document.getElementById("answer-result").innerHTML =  "";
            for(i=0;i<possibleResult.length;i++){
              content = `<div><a href="#" onclick="setResult('${possibleResult[i]}')">${possibleResult[i]}</a></div>`;
              document.getElementById("possible-list").innerHTML += content;
            }
        }
    });
});
resetQuest.addEventListener('click', function() {
    location.reload();
});

