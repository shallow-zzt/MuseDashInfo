package main

import (
	MDWebPageUtils "db/MDWebUtils"
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
var rankdb = SongDataUpdater.DBConnector("MDRankData.db")

func main() {

	// go func() {
	// 	for {
	// 		SongCatcher.Catcher()
	// 		time.Sleep(1 * time.Hour)
	// 	}
	// }()
	// go func() {
	// 	time.Sleep(15 * time.Minute)
	// 	for {
	// 		SongCatcher.SetUserTotalRks(rankdb)
	// 	}
	// }()

	http.HandleFunc("/css/{anything}", ServeStaticFile("static/css", "css"))
	http.HandleFunc("/js/{anything}", ServeStaticFile("static/js", "js"))
	http.HandleFunc("/pic/{anything}", ServeStaticFile("static/image", "webp"))
	http.HandleFunc("/alias/pic/{anything}", ServeStaticFile("SongAssets/Songpic", "png"))

	http.HandleFunc("/", homePageIndex)
	http.HandleFunc("/guessgame", miniGameIndex)
	http.HandleFunc("/alias", songAliasIndex)
	http.HandleFunc("/alias/{anything}", songAliasSettingIndex)
	http.HandleFunc("/value", songValueIndex)
	http.HandleFunc("/rank/{anything}", songRankShowIndex)
	http.HandleFunc("/user/{anything}", songUserData)
	http.HandleFunc("/search", songUserSearch)
	http.HandleFunc("/about", aboutPageIndex)

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
		case "webp":
			w.Header().Set("Content-Type", "image/webp")
		}
		http.ServeFile(w, r, filepath.Join(basepath, filename[len(filename)-1]))
	}
}

func homePageIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/home" {
		t, err := template.ParseFiles("Static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	} else {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
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
	query := r.URL.Query()
	searchInput := query.Get("search")
	var AllSongInfos MDWebPageUtils.FullSongAliasInfoWithCount
	if searchInput != "" {
		AllSongInfos = MDWebPageUtils.GetPartialSongInfo(db, searchInput)
	} else {
		AllSongInfos = MDWebPageUtils.GetAllSongInfo(db)
	}
	t.Execute(w, AllSongInfos)
}

func songUserSearch(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("Static/user-search.html")
	var songSearchResult MDWebPageUtils.SongUserSearch
	userPage := 1
	query := r.URL.Query()
	if pageStr := query.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil {
			userPage = parsedPage
		}
	}
	if userPage <= 0 || err != nil {
		userPage = 1
	}
	userInput := query.Get("search")
	if userInput == "" {
		songSearchResult.UserNum = -1
	} else {
		songSearchResult = MDWebPageUtils.GetSongUserSearchResult(rankdb, userInput, 10, (userPage-1)*10)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, songSearchResult)
}

func songUserData(w http.ResponseWriter, r *http.Request) {

	originURL := strings.Split(r.URL.Path, "/")
	userId := originURL[len(originURL)-1]

	t, err := template.ParseFiles("Static/userb50.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userPage := 1
	query := r.URL.Query()
	if pageStr := query.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil {
			userPage = parsedPage
		}
	}
	if userPage <= 0 || err != nil {
		userPage = 1
	}
	songUserDataList := MDWebPageUtils.GetUserSongList(rankdb, db, userId, 50, (userPage-1)*50)

	t.Execute(w, songUserDataList)
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

func songValueIndex(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("Static/value.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var allValueInfos MDWebPageUtils.FullSongValueInfoWithNum

	useOriginDiff := "1"
	searchInput := ""
	diffMaxValue := 13.0
	diffMinValue := -1.0
	hasMaxValue := true

	query := r.URL.Query()
	useOriginDiff = query.Get("originDiff")
	searchInput = query.Get("search")
	diffMaxValue, err = strconv.ParseFloat(query.Get("maxDiff"), 64)
	if err != nil {
		diffMaxValue = 13.0
		hasMaxValue = false
	}
	diffMinValue, err = strconv.ParseFloat(query.Get("minDiff"), 64)
	if err != nil {
		if hasMaxValue {
			diffMinValue = 1.0
		} else {
			diffMinValue = -1.0
		}
	}
	if useOriginDiff == "1" {
		allValueInfos = MDWebPageUtils.GetPartialSongValueInfo(db, searchInput, diffMinValue, diffMaxValue, 1)
	} else {
		allValueInfos = MDWebPageUtils.GetPartialSongValueInfo(db, searchInput, diffMinValue, diffMaxValue)
	}
	t.Execute(w, allValueInfos)
}

func songRankShowIndex(w http.ResponseWriter, r *http.Request) {
	originURL := strings.Split(r.URL.Path, "/")
	fullSongCode := strings.Split(originURL[len(originURL)-1], "-")
	AlbumCode, err := strconv.Atoi(fullSongCode[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	SongCode, err := strconv.Atoi(fullSongCode[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	diff, err := strconv.Atoi(r.URL.Query().Get("diff"))
	if err != nil {
		diff = 0
	}
	platform, err := strconv.Atoi(r.URL.Query().Get("platform"))
	if err != nil {
		platform = 0
	}

	t, err := template.ParseFiles("Static/rank-song.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rankData := SongDataUpdater.GetSongRankData(AlbumCode, SongCode, diff+1, platform)
	t.Execute(w, rankData)
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

func aboutPageIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("Static/about.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
