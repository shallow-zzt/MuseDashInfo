package MDWebPageUtils

import (
	"database/sql"
	"fmt"
	"log"
)

type AliasSubmitResult struct {
	Result int `json:"result"`
}

type FullSongAliasInfoWithCount struct {
	SongCount int
	SongList  []FullSongAliasInfo
}

type FullSongAliasInfo struct {
	db             *sql.DB
	AlbumCode      int
	SongCode       int
	SongPic        string
	SongName       string
	SongAliasCount int
	SongAliasList  []string
}

func GetAllSongInfo(db *sql.DB) FullSongAliasInfoWithCount {
	var aliasInfos []FullSongAliasInfo
	SQLCode := `SELECT music_album,music_album_number,music_pic_name,music_name FROM mdsong`
	result, err := db.Query(SQLCode)
	for result.Next() {
		var aliasInfo FullSongAliasInfo
		err = result.Scan(&aliasInfo.AlbumCode, &aliasInfo.SongCode, &aliasInfo.SongPic, &aliasInfo.SongName)
		if err != nil {
			log.Fatal(err)
		}
		getSQLAliasNum := `SELECT COUNT(music_alias) FROM mdalias WHERE music_album = ? AND music_album_number = ?`
		result2 := db.QueryRow(getSQLAliasNum, aliasInfo.AlbumCode, aliasInfo.SongCode)
		result2.Scan(&aliasInfo.SongAliasCount)
		aliasInfos = append(aliasInfos, aliasInfo)
	}
	return FullSongAliasInfoWithCount{SongCount: len(aliasInfos), SongList: aliasInfos}
}

func GetPartialSongInfo(db *sql.DB, userInput string) FullSongAliasInfoWithCount {
	var aliasInfos []FullSongAliasInfo
	SQLCode := `SELECT DISTINCT music_album,music_album_number,music_pic_name,music_name
	FROM(SELECT ms.music_album,ms.music_album_number,ms.music_pic_name,music_name,music_alias
	FROM mdsong ms
	LEFT OUTER JOIN mdalias ma
	ON ms.music_album = ma.music_album AND ms.music_album_number = ma.music_album_number)
	WHERE music_name LIKE "%` + userInput + `%" OR music_alias LIKE "%` + userInput + `%"`
	result, err := db.Query(SQLCode)
	for result.Next() {
		var aliasInfo FullSongAliasInfo
		err = result.Scan(&aliasInfo.AlbumCode, &aliasInfo.SongCode, &aliasInfo.SongPic, &aliasInfo.SongName)
		if err != nil {
			log.Fatal(err)
		}
		getSQLAliasNum := `SELECT COUNT(music_alias) FROM mdalias WHERE music_album = ? AND music_album_number = ?`
		result2 := db.QueryRow(getSQLAliasNum, aliasInfo.AlbumCode, aliasInfo.SongCode)
		result2.Scan(&aliasInfo.SongAliasCount)
		aliasInfos = append(aliasInfos, aliasInfo)
	}
	return FullSongAliasInfoWithCount{SongCount: len(aliasInfos), SongList: aliasInfos}
}

func GetSongAliasIsUsed(db *sql.DB, songName string) bool {
	SQLCode := `SELECT * FROM mdsong WHERE music_name = ?`
	result, err := db.Query(SQLCode, songName)
	if !result.Next() && err == nil {
		return false
	}
	return true
}

func (c *FullSongAliasInfo) GetAlias() FullSongAliasInfo {
	SQLCode := `SELECT music_alias FROM mdalias WHERE music_album = ? AND music_album_number = ?`
	result, err := c.db.Query(SQLCode, c.AlbumCode, c.SongCode)
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var buf string
		err = result.Scan(&buf)
		if err != nil {
			log.Fatal(err)
		}
		c.SongAliasList = append(c.SongAliasList, buf)
	}
	c.SongAliasCount = len(c.SongAliasList)
	return *c
}
