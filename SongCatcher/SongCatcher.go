package SongCatcher

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3" // Use SQLite driver
)

const (
	MAX_ALBUM        = 100
	MAX_ALBUM_NUMBER = 100
	MAX_DIFF         = 4
	DB_NAME          = "MDRankData.db"
	PLAY_TABLE       = "mdplay"
	USER_TABLE       = "mduser"
)

func Catcher() {

	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}
	defer db.Close()
	createTable(db)

	for songAlbum := 0; songAlbum <= MAX_ALBUM; songAlbum++ {
		for songAlbumNumber := 0; songAlbumNumber <= MAX_ALBUM; songAlbumNumber++ {
			for songDiff := 1; songDiff <= MAX_DIFF; songDiff++ {
				for platform := 0; platform <= 1; platform++ {
					for offset := 0; offset <= 17; offset++ {
						if offset%5 == 0 {
							t := time.After(1 * time.Second)
							<-t
						}
						apiData, err := getAPIData(platform, songAlbum, songAlbumNumber, songDiff, offset)

						mdStatus := apiData.(map[string]interface{})["total"]
						fmt.Println(mdStatus)
						if mdStatus != nil && (mdStatus.(float64)) != 2000 {
							break
						}
						mdDatas := apiData.(map[string]interface{})["result"].([]interface{})
						for rankOffset, mdData := range mdDatas {
							var dbData MDdbData
							playData := mdData.(map[string]interface{})["play"].(map[string]interface{})
							userData := mdData.(map[string]interface{})["user"].(map[string]interface{})
							fmt.Println(playData)
							dbData.platform = platform
							dbData.musicAlbum = songAlbum
							dbData.musicAlbumNumber = songAlbumNumber
							dbData.musicDiff = songDiff
							dbData.rank = offset*100 + rankOffset + 1
							dbData.userId = playData["user_id"].(string)
							dbData.acc = playData["acc"].(float64)
							dbData.miss = int(playData["miss"].(float64))
							dbData.score = int(playData["score"].(float64))
							dbData.judge = playData["judge"].(string)
							dbData.characterId, err = strconv.Atoi(playData["character_uid"].(string))
							dbData.elfinId, err = strconv.Atoi(playData["elfin_uid"].(string))
							dbData.playTime = playData["updated_at"].(string)
							dbData.nickname = userData["nickname"].(string)
							fmt.Println(dbData.rank)
							insertMDData(db, dbData)
						}
						if err != nil {
							fmt.Println("获取API数据失败:", err)
						}
					}
				}
			}
		}
	}

}

func getAPIData(platform int, songAlbum int, songAlbumNumber int, songDiff int, offset int) (interface{}, error) {
	var platformPrefix string
	if platform == 0 {
		platformPrefix = "pcleaderboard"
	} else {
		platformPrefix = "leaderboard"
	}
	apiUrl := "https://prpr-muse-dash.peropero.net/musedash/v1/" + platformPrefix + "/top?music_uid=" + strconv.Itoa(songAlbum) + "-" + strconv.Itoa(songAlbumNumber) + "&music_difficulty=" + strconv.Itoa(songDiff) + "&limit=100&offset=" + strconv.Itoa(offset) + "&version=1.5.0&platform=1"
	var apiData interface{}

	resp, err := http.Get(apiUrl)
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &apiData)
	if err != nil {
		return nil, err
	}

	return apiData, nil
}

func createTable(db *sql.DB) {
	play_table_sql := `
        CREATE TABLE IF NOT EXISTS ` + PLAY_TABLE + ` (
			platform INTEGER NOT NULL,
			music_album_id INTEGER NOT NULL,
			music_album_song_id INTEGER NOT NULL,
			music_difficulty INTEGER NOT NULL,
			rank INTEGER NOT NULL,
			user_id TEXT,
			acc REAL,
			miss INTEGER,
            score INTEGER,
			judge TEXT,
			character_uid INTEGER,
			elfin_uid INTEGER,
			updated_at TEXT,
			PRIMARY KEY (platform, music_album_id,music_album_song_id, music_difficulty, rank)
        );
    `
	user_table_sql := `
        CREATE TABLE IF NOT EXISTS ` + USER_TABLE + ` (
			user_id TEXT PRIMARY KEY NOT NULL,
            nickname TEXT
        );
    `
	_, err := db.Exec(play_table_sql)
	_, err = db.Exec(user_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func insertMDData(db *sql.DB, dbData MDdbData) error {
	queryMDData := `INSERT INTO 
	mdplay (platform, music_album_id,music_album_song_id,music_difficulty,rank,user_id,acc,miss,score,judge, character_uid,elfin_uid,updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(platform, music_album_id,music_album_song_id,music_difficulty,rank)
	DO UPDATE SET user_id=excluded.user_id,acc=excluded.acc,miss=excluded.miss,score=excluded.score,judge=excluded.judge,character_uid=excluded.character_uid,elfin_uid=excluded.elfin_uid,updated_at=excluded.updated_at`
	queryUser := `INSERT OR REPLACE INTO 
	mduser (user_id,nickname) 
	VALUES (?, ?)
	ON CONFLICT(user_id)
	DO UPDATE SET nickname=excluded.nickname`
	statement, err := db.Prepare(queryMDData)
	_, err = statement.Exec(dbData.platform, dbData.musicAlbum, dbData.musicAlbumNumber, dbData.musicDiff, dbData.rank, dbData.userId, dbData.acc, dbData.miss, dbData.score, dbData.judge, dbData.characterId, dbData.elfinId, dbData.playTime)
	statement, err = db.Prepare(queryUser)
	_, err = statement.Exec(dbData.userId, dbData.nickname)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}
