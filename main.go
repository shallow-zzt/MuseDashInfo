package main

import (
	"bytes"
	"database/sql"
	MiniGame "db/Minigame"
	"db/SongDataUpdater"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"image/png"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
)

type MiniGameContent struct {
	Rotateb64      string
	Answerb64      string
	Cuttedb64      string
	SongAnswerName string
}

type MiniGameResult struct {
	Result         int      `json:"result"`
	ResultNum      int      `json:"result-num"`
	PossibleResult []string `json:"possible-result"`
}

var store = sessions.NewCookieStore([]byte("^_^"))
var db = SongDataUpdater.DBConnector()

func main() {

	//SongDataUpdater.SongUpdater()
	//go SongCatcher.Catcher()
	http.HandleFunc("/css/style.css", ServeStaticFile("style.css", "css"))
	http.HandleFunc("/js/index.js", ServeStaticFile("index.js", "js"))

	http.HandleFunc("/guessgame", miniGameIndex)
	http.HandleFunc("/submit/answer", JudgeSubmit)
	http.ListenAndServe(":8080", nil)

}

func ServeStaticFile(filename string, fileext string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if fileext == "js" {
			w.Header().Set("Content-Type", "application/javascript")
			http.ServeFile(w, r, filepath.Join("static/js", filename))
		} else if fileext == "css" {
			w.Header().Set("Content-Type", "text/css")
			http.ServeFile(w, r, filepath.Join("static/css", filename))
		}
	}
}

func miniGameIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	t, err := template.ParseFiles("Static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	miniGameData := generateMiniGameContent(false)

	t.Execute(w, MiniGameContent{
		Rotateb64:      miniGameData.Rotateb64,
		Cuttedb64:      miniGameData.Cuttedb64,
		Answerb64:      miniGameData.Answerb64,
		SongAnswerName: miniGameData.SongAnswerName})
}

func JudgeSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]string
	var MiniGameResult MiniGameResult

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(requestData)
	answer := requestData["answer"]
	answerName := requestData["standard-answer"]

	answerList := getPossibleAnswer(db, answer)

	if len(answerList) == 0 || len(answerList) >= 5 {
		MiniGameResult.Result = 0
		MiniGameResult.ResultNum = 0
		MiniGameResult.PossibleResult = nil
	} else if len(answerList) == 1 {
		MiniGameResult.PossibleResult = answerList
		MiniGameResult.ResultNum = 1
		if answer == answerName {
			MiniGameResult.Result = 1
		} else if answer == answerList[0] && answer != answerName {
			MiniGameResult.Result = 0
			MiniGameResult.ResultNum = 0
			MiniGameResult.PossibleResult = nil
		} else {
			MiniGameResult.Result = 2
		}
	} else if len(answerList) >= 2 && len(answerList) <= 4 {
		MiniGameResult.Result = 2
		MiniGameResult.ResultNum = len(answerList)
		MiniGameResult.PossibleResult = answerList
	}

	fmt.Println(MiniGameResult)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MiniGameResult)
}

func generateMiniGameContent(isRotate bool) MiniGameContent {

	songPicName, songAnswerName := getRandomSongPicName(db)
	songPicPath := "SongAssets/Songpic/" + songPicName

	rotatePic, answerPic, cuttedPic := MiniGame.ProductingSongPic(songPicPath, isRotate)

	rotateBuf := new(bytes.Buffer)
	answerBuf := new(bytes.Buffer)
	cuttedBuf := new(bytes.Buffer)

	err := png.Encode(rotateBuf, rotatePic)
	err = png.Encode(answerBuf, answerPic)
	err = png.Encode(cuttedBuf, cuttedPic)
	if err != nil {
		fmt.Println(err)
	}

	rotateb64 := base64.StdEncoding.EncodeToString(rotateBuf.Bytes())
	answerb64 := base64.StdEncoding.EncodeToString(answerBuf.Bytes())
	cuttedb64 := base64.StdEncoding.EncodeToString(cuttedBuf.Bytes())

	return MiniGameContent{Rotateb64: rotateb64, Answerb64: answerb64, Cuttedb64: cuttedb64, SongAnswerName: songAnswerName}
}

func getRandomSongPicName(db *sql.DB) (string, string) {
	var songPicName string
	var songName string
	getRandomSongPic := `SELECT music_pic_name,music_name
	FROM mdsong
	ORDER BY RANDOM()
	LIMIT 1`
	result := db.QueryRow(getRandomSongPic)
	result.Scan(&songPicName, &songName)
	return songPicName, songName
}

func getPossibleAnswer(db *sql.DB, answer string) []string {
	var possibleAnswers []string
	getPossibleAnswer := `SELECT music_name
	FROM mdsong
	WHERE music_name LIKE "%` + answer + `%"`
	result, err := db.Query(getPossibleAnswer)

	for result.Next() {
		var name string
		err = result.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		possibleAnswers = append(possibleAnswers, name)
	}

	return possibleAnswers
}
