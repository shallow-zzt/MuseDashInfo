package main

import (
	MDWebPageUtils "db/Shiori3WebUtils"
	"db/SongCatcher"
	"db/SongDataUpdater"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var db = SongDataUpdater.DBConnector()

func main() {

	go SongCatcher.Catcher()

	http.HandleFunc("/css/{anything}", ServeStaticFile("static/css", "css"))
	http.HandleFunc("/js/{anything}", ServeStaticFile("static/js", "js"))
	http.HandleFunc("/alias/pic/{anything}", ServeStaticFile("SongAssets/Songpic", "png"))

	http.HandleFunc("/guessgame", miniGameIndex)
	http.HandleFunc("/alias", songAliasIndex)
	http.HandleFunc("/alias/{anything}", songAliasSettingIndex)

	http.HandleFunc("/submit/guessgame/answer", guessGmaeJudgeSubmit)
	http.HandleFunc("/submit/aliassong/alias", songAliasSubmit)

	http.ListenAndServe(":8080", nil)

}

func ServeStaticFile(basepath string, fileext string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := strings.Split(r.URL.Path, "/")
		switch fileext {
		case "css":
			w.Header().Set("Content-Type", "text/css")
		case "js":
			w.Header().Set("Content-Type", "application/javascript")
		case "png":
			w.Header().Set("Content-Type", "image/png")
		}
		http.ServeFile(w, r, filepath.Join(basepath, filename[len(filename)-1]))
	}
}

func miniGameIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	t, err := template.ParseFiles("Static/guessgame.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	miniGameData := MDWebPageUtils.GenerateMiniGameContent(db, false)

	t.Execute(w, MDWebPageUtils.MiniGameContent{
		Rotateb64:      miniGameData.Rotateb64,
		Cuttedb64:      miniGameData.Cuttedb64,
		Answerb64:      miniGameData.Answerb64,
		SongAnswerName: miniGameData.SongAnswerName})
}

func songAliasIndex(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("Static/alias.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	AllSongInfos := MDWebPageUtils.GetAllSongInfo(db)
	t.Execute(w, AllSongInfos)
}

func songAliasSettingIndex(w http.ResponseWriter, r *http.Request) {

	originURL := strings.Split(r.URL.Path, "/")
	fullSongCode := strings.Split(originURL[len(originURL)-1], "-")
	AlbumCode, err := strconv.Atoi(fullSongCode[0])
	SongCode, err := strconv.Atoi(fullSongCode[1])

	songAliasInfo := MDWebPageUtils.GetSongInfoFromCode(db, AlbumCode, SongCode).GetAlias()

	t, err := template.ParseFiles("Static/alias-song.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if songAliasInfo.SongName == "" {
		http.Error(w, "No Music Founded", http.StatusNotFound)
		return
	}

	t.Execute(w, songAliasInfo)
}

func songAliasSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submitContent map[string]string
	var submitResult MDWebPageUtils.AliasSubmitResult

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &submitContent)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(submitContent)

	albumCode, err := strconv.Atoi(submitContent["album-code"])
	songCode, err := strconv.Atoi(submitContent["song-code"])
	inputAlias := submitContent["input-alias"]

	if inputAlias == "" {
		submitResult.Result = 2
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(submitResult)
		return
	}
	if MDWebPageUtils.GetSongAliasIsUsed(db, inputAlias) {
		submitResult.Result = 3
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(submitResult)
		return
	}

	err = SongDataUpdater.InsertAliasData(db, SongDataUpdater.AliasBasicData{MusicAlbum: albumCode, MusicAlbumNumber: songCode, MusicAlias: inputAlias})
	if err != nil {
		submitResult.Result = 0
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(submitResult)
		return
	}

	submitResult.Result = 1
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(submitResult)

}

func guessGmaeJudgeSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]string
	var MiniGameResult MDWebPageUtils.MiniGameResult

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(requestData)
	answer := requestData["answer"]
	answerName := requestData["standard-answer"]

	answerList := MDWebPageUtils.GetPossibleAnswer(db, answer)

	if len(answerList) == 0 || len(answerList) >= 6 {
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
	} else if len(answerList) >= 2 && len(answerList) <= 5 {
		MiniGameResult.Result = 2
		MiniGameResult.ResultNum = len(answerList)
		MiniGameResult.PossibleResult = answerList
	}

	fmt.Println(MiniGameResult)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MiniGameResult)
}
