package SongCatcher

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	MAX_ALBUM         = 100
	MAX_ALBUM_NUMBER  = 100
	MAX_DIFF          = 4
	DB_NAME           = "MDRankData.db"
	PLAY_TABLE        = "mdplay"
	USER_TABLE        = "mduser"
	RKS_HISTORY_TABLE = "mdrks"
)

func DBConnector(filepath ...string) *sql.DB {
	var SQLpath string
	if len(filepath) == 0 {
		SQLpath = "mdsong.db"
	} else {
		SQLpath = filepath[0]
	}
	db, err := sql.Open("sqlite3", SQLpath)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return nil
	}
	return db
}

func GetSongValue(db *sql.DB, musicAlbum int, musicAlbumNumber int, musicDiff int) float64 {
	var songValue float64
	SQLCode := `SELECT music_value FROM mdvalue WHERE music_album = ? AND music_album_number = ? AND music_diff_tier = ?`
	result := db.QueryRow(SQLCode, musicAlbum, musicAlbumNumber, musicDiff-1)
	result.Scan(&songValue)
	return songValue
}

func calRks(acc float64, rank int, musicAlbum int, musicAlbumNumber int, musicDiff int) float64 {
	db := DBConnector()
	songRKS := GetSongValue(db, musicAlbum, musicAlbumNumber, musicDiff)
	defer db.Close()
	if songRKS < 0 {
		return -1
	}
	index := float64(1001-rank)/1000*0.04 + 1
	return acc * index * songRKS / 100
}

func Catcher(albumStart ...int) {
	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}
	defer db.Close()
	createTable(db)

	var albumStartNum int
	if len(albumStart) == 0 {
		albumStartNum = 0
	} else {
		albumStartNum = albumStart[0]
	}

	for songAlbum := albumStartNum; songAlbum <= MAX_ALBUM; songAlbum++ {
		for songAlbumNumber := 0; songAlbumNumber <= MAX_ALBUM; songAlbumNumber++ {
			fmt.Println(songAlbum, songAlbumNumber)
			for songDiff := 1; songDiff <= MAX_DIFF; songDiff++ {
				for platform := 0; platform <= 1; platform++ {
					for offset := 0; offset <= 17; offset++ {
						rand.Seed(time.Now().UnixNano())
						randCoolDown := rand.Intn(3) + 2
						apiData, err := GetAPIData(platform, songAlbum, songAlbumNumber, songDiff, offset)
						if err != nil {
							continue
						}
						if offset%2 == 0 {
							t := time.After(time.Duration(randCoolDown-1) * time.Second)
							<-t
						}
						mdStatus := apiData.(map[string]interface{})["total"]
						fmt.Println(mdStatus)
						if mdStatus == nil || (mdStatus.(float64)) != 2000 {
							platform = 2
							songDiff = MAX_DIFF + 1
							break
						}

						mdDatas := apiData.(map[string]interface{})["result"].([]interface{})
						for rankOffset, mdData := range mdDatas {
							var dbData MDdbData
							playData := mdData.(map[string]interface{})["play"].(map[string]interface{})
							userData := mdData.(map[string]interface{})["user"].(map[string]interface{})
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
							dbData.rks = calRks(dbData.acc, dbData.rank, dbData.musicAlbum, dbData.musicAlbumNumber, dbData.musicDiff)
							fmt.Println(dbData)
							err := insertMDData(db, dbData)
							if err != nil {
								fmt.Println("数据插入失败", err)
							}
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

func GetAPIData(platform int, songAlbum int, songAlbumNumber int, songDiff int, offset int) (interface{}, error) {
	var platformPrefix string
	if platform == 0 {
		platformPrefix = "pcleaderboard"
	} else {
		platformPrefix = "leaderboard"
	}
	apiUrl := "https://prpr-muse-dash.peropero.net/musedash/v1/" + platformPrefix + "/top?music_uid=" + strconv.Itoa(songAlbum) + "-" + strconv.Itoa(songAlbumNumber) + "&music_difficulty=" + strconv.Itoa(songDiff) + "&limit=100&offset=" + strconv.Itoa(offset) + "&version=1.5.0&platform=1"
	var apiData interface{}

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &apiData)
	if err != nil {
		return nil, err
	}
	return apiData, nil
}

func insertMDData(db *sql.DB, dbData MDdbData) error {
	db.Exec("PRAGMA journal_mode=WAL;")
	queryMDData := `INSERT INTO 
	mdplay (platform, music_album_id,music_album_song_id,music_difficulty,rank,user_id,acc,miss,score,judge, character_uid,elfin_uid,updated_at,rks) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(platform, music_album_id,music_album_song_id,music_difficulty,rank)
	DO UPDATE SET user_id=excluded.user_id,acc=excluded.acc,miss=excluded.miss,score=excluded.score,judge=excluded.judge,character_uid=excluded.character_uid,elfin_uid=excluded.elfin_uid,updated_at=excluded.updated_at,rks=excluded.rks`
	queryUser := `INSERT OR REPLACE INTO 
	mduser (user_id,nickname,user_rks) 
	VALUES (?, ?, ?)
	ON CONFLICT(user_id)
	DO UPDATE SET nickname=excluded.nickname`
	statement, err := db.Prepare(queryMDData)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(dbData.platform, dbData.musicAlbum, dbData.musicAlbumNumber, dbData.musicDiff, dbData.rank, dbData.userId, dbData.acc, dbData.miss, dbData.score, dbData.judge, dbData.characterId, dbData.elfinId, dbData.playTime, dbData.rks)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	statement, err = db.Prepare(queryUser)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(dbData.userId, dbData.nickname, -1)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	return nil
}
