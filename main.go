package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Use SQLite driver
)

const (
	ASSETS_PATH = "SongAssets/Json/"
)

type AlbumData struct {
	albumCode    string
	albumNameCN  string
	albumNameCNT string
	albumNameEN  string
	albumNameJA  string
	albumNameKR  string
}

type SongData struct {
	musicAlbum       int
	musicAlbumNumber int
	musicName        string
	musicAuthor      string
	musicBPM         string
	musicPicName     string
	musicSceneName   string
	musicSheetAuthor []string
	musicDiff        []string
	musicSpecialDiff string
}

func main() {

	db, err := sql.Open("sqlite3", "mdsong.db")
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}
	defer db.Close()
	createTable(db)

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
					Difficulty, ok := albumSongDetail.(map[string]interface{})["difficulty"+strconv.Itoa(j)].(string)
					if ok {
						musicData.musicDiff = append(musicData.musicDiff, Difficulty)
					} else {
						musicData.musicDiff = append(musicData.musicDiff, "0")
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

func getJsonRawFile(filename string) interface{} {
	var songRawData interface{}
	data, err := ioutil.ReadFile(ASSETS_PATH + filename)
	err = json.Unmarshal(data, &songRawData)
	if err != nil {
		fmt.Println(err)
	}
	return songRawData
}

func createTable(db *sql.DB) {
	play_table_sql := `
	CREATE TABLE  IF NOT EXISTS mdsong(
		music_album INTEGER NOT NULL,
		music_album_number INTEGER NOT NULL,
		music_name TEXT,
		music_author TEXT,
		music_bPM TEXT,
		music_pic_name TEXT,
		music_scene_name TEXT,
		music_sheet_author TEXT, -- Store as JSON string
		music_diff TEXT,        -- Store as JSON string
		music_diff_special TEXT,
		PRIMARY KEY (music_album, music_album_number)
	);
    `
	_, err := db.Exec(play_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func insertMDData(db *sql.DB, musicData SongData) error {
	fmt.Println(musicData)
	queryMDData := `INSERT INTO 
	mdsong (music_album,music_album_number,music_name,music_author,music_bpm,music_pic_name,music_scene_name,music_sheet_author,music_diff,music_diff_special) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(music_album,music_album_number)
	DO UPDATE SET music_name=excluded.music_name,music_author=excluded.music_author,
	music_bpm=excluded.music_bpm,music_pic_name=excluded.music_pic_name,
	music_scene_name=excluded.music_scene_name,music_sheet_author=excluded.music_sheet_author,
	music_diff=excluded.music_diff,music_diff_special=excluded.music_diff_special`
	statement, err := db.Prepare(queryMDData)
	_, err = statement.Exec(musicData.musicAlbum, musicData.musicAlbumNumber, musicData.musicName, musicData.musicAuthor, musicData.musicBPM, musicData.musicPicName, musicData.musicSceneName, strings.Join(musicData.musicSheetAuthor, " | "), strings.Join(musicData.musicDiff, ","), musicData.musicSpecialDiff)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}
