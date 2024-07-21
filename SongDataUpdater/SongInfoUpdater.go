package SongDataUpdater

import (
	"database/sql"
	"db/SongCatcher"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Use SQLite driver
)

const (
	ASSETS_PATH = "SongAssets/Json/"
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

func SongUpdater() {

	db := DBConnector()
	songTableInit(db)

	albumDatas := getJsonRawFile("albums.json").([]interface{})
	albumDatasCN := getJsonRawFile("albums_ChineseS.json").([]interface{})
	albumDatasCNT := getJsonRawFile("albums_ChineseT.json").([]interface{})
	albumDatasEN := getJsonRawFile("albums_English.json").([]interface{})
	albumDatasJA := getJsonRawFile("albums_Japanese.json").([]interface{})
	albumDatasKR := getJsonRawFile("albums_Korean.json").([]interface{})
	for i, albumData := range albumDatas {
		var albumDetail AlbumData
		albumDetail.albumCode = albumData.(map[string]interface{})["jsonName"].(string)
		albumDetail.albumNameCN = albumDatasCN[i].(map[string]interface{})["title"].(string)
		albumDetail.albumNameCNT = albumDatasCNT[i].(map[string]interface{})["title"].(string)
		albumDetail.albumNameEN = albumDatasEN[i].(map[string]interface{})["title"].(string)
		albumDetail.albumNameJA = albumDatasJA[i].(map[string]interface{})["title"].(string)
		albumDetail.albumNameKR = albumDatasKR[i].(map[string]interface{})["title"].(string)
		insertAlbumData(db, albumDetail)
		fmt.Println(albumDetail)
		if albumDetail.albumCode != "" {
			albumSongDetails := getJsonRawFile(albumDetail.albumCode + ".json").([]interface{})
			for _, albumSongDetail := range albumSongDetails {
				var musicData SongData
				var err error
				var ok bool
				musicAlbumCode := strings.Split(albumSongDetail.(map[string]interface{})["uid"].(string), "-")
				musicData.musicAlbum, err = strconv.Atoi(musicAlbumCode[0])
				musicData.musicAlbumNumber, err = strconv.Atoi(musicAlbumCode[1])
				musicData.musicName = albumSongDetail.(map[string]interface{})["name"].(string)
				musicData.musicAuthor = albumSongDetail.(map[string]interface{})["author"].(string)
				musicData.musicBPM = albumSongDetail.(map[string]interface{})["bpm"].(string)
				musicData.musicPicName = albumSongDetail.(map[string]interface{})["cover"].(string)
				musicData.musicSceneName = albumSongDetail.(map[string]interface{})["scene"].(string)
				levelDesigner, ok := albumSongDetail.(map[string]interface{})["levelDesigner"].(string)
				if !ok {
					for j := 1; j <= 4; j++ {
						levelDesigner, ok = albumSongDetail.(map[string]interface{})["levelDesigner"+strconv.Itoa(j)].(string)
						if ok {
							musicData.musicSheetAuthor = append(musicData.musicSheetAuthor, levelDesigner)
						}

					}
				} else {
					musicData.musicSheetAuthor = append(musicData.musicSheetAuthor, levelDesigner)
				}

				for j := 1; j <= 4; j++ {
					t := time.After(1 * time.Second)
					<-t
					Difficulty, ok := albumSongDetail.(map[string]interface{})["difficulty"+strconv.Itoa(j)].(string)
					if ok && Difficulty != "0" {

						musicData.musicDiff = append(musicData.musicDiff, Difficulty)
						notes, score := GetSongNotesAndScore(musicData.musicAlbum, musicData.musicAlbumNumber, j)
						musicData.musicNotes = append(musicData.musicNotes, notes)
						musicData.musicScore = append(musicData.musicScore, score)
					} else {
						musicData.musicDiff = append(musicData.musicDiff, "0")
						musicData.musicNotes = append(musicData.musicNotes, "0")
						musicData.musicScore = append(musicData.musicScore, "0")
					}
				}
				musicData.musicSpecialDiff, ok = albumSongDetail.(map[string]interface{})["difficulty5"].(string)
				if !ok {
					musicData.musicSpecialDiff = "0"

				}
				if err != nil {
					fmt.Println(err)
				}
				insertMDData(db, musicData)

			}
		}
	}

}

func GetSongNotesAndScore(songAlbum int, songAlbumNumber int, songDiff int) (string, string) {
	apiData, err := SongCatcher.GetAPIData(0, songAlbum, songAlbumNumber, songDiff, 0)
	if err != nil {
		return "", ""
	}
	mdData := apiData.(map[string]interface{})["result"]
	if mdData == nil {
		fmt.Println("No Notes Data Found")
		return "0", "0"
	}
	fmt.Println(songAlbum, songAlbumNumber)
	NotesData := strconv.FormatFloat(mdData.([]interface{})[5].(map[string]interface{})["play"].(map[string]interface{})["combo"].(float64), 'f', 0, 64)
	ScoreData := strconv.FormatFloat(mdData.([]interface{})[5].(map[string]interface{})["play"].(map[string]interface{})["score"].(float64), 'f', 0, 64)
	return NotesData, ScoreData
}

func GetBasicAilas() {
	db := DBConnector()
	ailasTableInit(db)

	aliasDatas := getJsonRawFile("music_search_tag.json").([]interface{})
	for _, aliasData := range aliasDatas {
		var aliasBasicData AliasBasicData
		var err error
		uid := strings.Split(aliasData.(map[string]interface{})["uid"].(string), "-")
		tag := aliasData.(map[string]interface{})["tag"]
		if tag != nil {
			aliases := tag.([]interface{})
			aliasBasicData.MusicAlbum, err = strconv.Atoi(uid[0])
			aliasBasicData.MusicAlbumNumber, err = strconv.Atoi(uid[1])
			for _, alias := range aliases {
				aliasBasicData.MusicAlias = alias.(string)
				fmt.Println(aliasBasicData)
				InsertAliasData(db, aliasBasicData)
			}
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func getJsonRawFile(filename string) interface{} {
	var songRawData interface{}
	data, err := ioutil.ReadFile(ASSETS_PATH + filename)
	err = json.Unmarshal(data, &songRawData)
	if err != nil {
		fmt.Println(err)
	}
	return songRawData
}

func ailasTableInit(db *sql.DB) {
	play_table_sql := `
	CREATE TABLE  IF NOT EXISTS mdalias(
		music_album INTEGER NOT NULL,
		music_album_number INTEGER NOT NULL,
		music_alias TEXT,
		PRIMARY KEY (music_album, music_album_number,music_alias)
		FOREIGN KEY (music_album, music_album_number) REFERENCES mdsong(music_album, music_album_number)
	);
    `
	_, err := db.Exec(play_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func songTableInit(db *sql.DB) {
	play_table_sql := `
	CREATE TABLE  IF NOT EXISTS mdsong(
		music_album INTEGER NOT NULL,
		music_album_number INTEGER NOT NULL,
		music_name TEXT,
		music_author TEXT,
		music_bPM TEXT,
		music_pic_name TEXT,
		music_scene_name TEXT,
		music_sheet_author TEXT, 
		music_diff TEXT,        
		music_diff_special TEXT,
		music_total_note TEXT,
		music_total_score TEXT,
		PRIMARY KEY (music_album, music_album_number)
	);
    `
	album_table_sql := `
	CREATE TABLE  IF NOT EXISTS mdalbum(
		music_album TEXT NOT NULL,
		music_album_name_CN TEXT,
		music_album_name_CNT TEXT,
		music_album_name_EN TEXT,
		music_album_name_JA TEXT,
		music_album_name_KR TEXT,
		PRIMARY KEY (music_album)
	);
    `
	_, err := db.Exec(play_table_sql)
	_, err = db.Exec(album_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func insertMDData(db *sql.DB, musicData SongData) error {
	queryMDData := `INSERT INTO 
	mdsong (music_album,music_album_number,music_name,music_author,music_bpm,music_pic_name,music_scene_name,music_sheet_author,music_diff,music_diff_special,music_total_note,music_total_score) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(music_album,music_album_number)
	DO UPDATE SET music_name=excluded.music_name,music_author=excluded.music_author,
	music_bpm=excluded.music_bpm,music_pic_name=excluded.music_pic_name,
	music_scene_name=excluded.music_scene_name,music_sheet_author=excluded.music_sheet_author,
	music_diff=excluded.music_diff,music_diff_special=excluded.music_diff_special,
	music_total_note=excluded.music_total_note,music_total_score=excluded.music_total_score`
	statement, err := db.Prepare(queryMDData)
	_, err = statement.Exec(musicData.musicAlbum, musicData.musicAlbumNumber, musicData.musicName, musicData.musicAuthor, musicData.musicBPM, musicData.musicPicName, musicData.musicSceneName, strings.Join(musicData.musicSheetAuthor, " | "), strings.Join(musicData.musicDiff, ","), musicData.musicSpecialDiff, strings.Join(musicData.musicNotes, ","), strings.Join(musicData.musicScore, ","))
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}

func insertAlbumData(db *sql.DB, AlbumData AlbumData) error {
	queryMDData := `INSERT INTO 
	mdalbum (music_album,music_album_name_CN,music_album_name_CNT,music_album_name_EN,music_album_name_JA,music_album_name_KR) 
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(music_album)
	DO UPDATE SET music_album_name_CN=excluded.music_album_name_CN,music_album_name_CNT=excluded.music_album_name_CNT,
	music_album_name_EN=excluded.music_album_name_EN,music_album_name_JA=excluded.music_album_name_JA,
	music_album_name_KR=excluded.music_album_name_KR`
	statement, err := db.Prepare(queryMDData)
	_, err = statement.Exec(AlbumData.albumCode, AlbumData.albumNameCN, AlbumData.albumNameCNT, AlbumData.albumNameEN, AlbumData.albumNameJA, AlbumData.albumNameKR)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}

func InsertAliasData(db *sql.DB, aliasData AliasBasicData) error {
	queryMDData := `INSERT INTO 
	mdalias (music_album,music_album_number,music_alias) 
	VALUES (?, ?, ?)`
	statement, err := db.Prepare(queryMDData)
	_, err = statement.Exec(aliasData.MusicAlbum, aliasData.MusicAlbumNumber, aliasData.MusicAlias)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}
