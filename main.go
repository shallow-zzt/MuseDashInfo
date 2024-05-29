package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/mattn/go-sqlite3" // Use SQLite driver
)

const (
	DB_NAME    = "test.db"
	PLAY_TABLE = "mdplay"
	USER_TABLE = "mduser"
)

func main() {

	apiData, err := getAPIData()
	mdDatas := apiData.(map[string]interface{})["result"].([]interface{})
	for _, mdData := range mdDatas {
		playData := mdData.(map[string]interface{})["play"].(map[string]interface{})
		userData := mdData.(map[string]interface{})["user"].(map[string]interface{})
		fmt.Println(playData["score"])
		fmt.Println(userData["nickname"])

	}
	if err != nil {
		fmt.Println("获取API数据失败:", err)
		return
	}

	// 连接数据库
	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		return
	}
	defer db.Close()

	// 创建表（如果不存在）
	createTable(db)
}

func getAPIData() (interface{}, error) {
	apiUrl := "https://prpr-muse-dash.peropero.net/musedash/v1/pcleaderboard/top?music_uid=0-44&music_difficulty=3&limit=100&offset=1&version=1.5.0&platform=1"
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
            id INTEGER PRIMARY KEY,
			platform INTEGER,
			user_id TEXT,
			music_uid TEXT,
			music_difficulty INTEGER,
			rank INTEGER,
			acc REAL,
			miss INTEGER,
            score INTEGER,
			judge TEXT,
			character_uid INTEGER,
			elfin_uid INTEGER,
			updated_at INTEGER
        );
    `
	user_table_sql := `
        CREATE TABLE IF NOT EXISTS ` + USER_TABLE + ` (
            id INTEGER PRIMARY KEY,
            nickname TEXT,
            user_id TEXT
        );
    `
	_, err := db.Exec(play_table_sql)
	_, err = db.Exec(user_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}
