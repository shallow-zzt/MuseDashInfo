package MDWebPageUtils

import (
	"database/sql"
	"strings"
	"unicode/utf8"
)

func GetTotalRKS(db *sql.DB, userID string, bestSongs int) float64 {
	var averageRKS float64
	SQLCode := `SELECT AVG(rks) AS rks_average
	FROM(
	SELECT rks FROM mdplay WHERE user_id=? ORDER BY rks DESC LIMIT ?
	)`
	result := db.QueryRow(SQLCode, userID, bestSongs)
	result.Scan(&averageRKS)
	return averageRKS
}

func GetTotalUserNum(db *sql.DB, userInput string) int {
	var songNum int
	var SQLCode string
	if len(userInput) < 11 {
		SQLCode = `SELECT COUNT(*) FROM mduser WHERE nickname LIKE "%` + userInput + `%" AND LENGTH(nickname)<11`
	} else {
		SQLCode = `SELECT COUNT(*) FROM mduser WHERE nickname LIKE "玩家` + userInput + `" OR user_id is "` + userInput + `"`
	}
	result := db.QueryRow(SQLCode)
	result.Scan(&songNum)
	return songNum
}

func GetTotalSongNum(db *sql.DB, userID string) int {
	var userNum int
	SQLCode := `SELECT COUNT(*) FROM mdplay WHERE user_id=?`
	result := db.QueryRow(SQLCode, userID)
	result.Scan(&userNum)
	return userNum
}

func GetUserName(db *sql.DB, userID string) string {
	var userName string
	SQLCode := `SELECT nickname FROM mduser WHERE user_id is ?`
	result := db.QueryRow(SQLCode, userID)
	result.Scan(&userName)
	return userName
}

func GetSongInfoFromCode(db *sql.DB, albumCode int, songCode int) *FullSongAliasInfo {
	var FullSongAliasInfo FullSongAliasInfo
	FullSongAliasInfo.AlbumCode = albumCode
	FullSongAliasInfo.SongCode = songCode
	FullSongAliasInfo.db = db
	SQLCode := `SELECT music_pic_name,music_name FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&FullSongAliasInfo.SongPic, &FullSongAliasInfo.SongName)
	return &FullSongAliasInfo
}

func GetSongPicFromCode(db *sql.DB, albumCode int, songCode int) string {
	var songPicName string
	SQLCode := `SELECT music_pic_name FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&songPicName)
	if songPicName == "" {
		songPicName = "random_song_cover"
	}
	return songPicName
}

func GetSongNameFromCode(db *sql.DB, albumCode int, songCode int) string {
	var songName string
	SQLCode := `SELECT music_name FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&songName)
	if songName == "" {
		songName = "未知曲目"
	}
	return songName
}

func GetSongDiffFromCode(db *sql.DB, albumCode int, songCode int) []string {
	var songDiffOrigin string
	SQLCode := `SELECT music_diff FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&songDiffOrigin)
	if songDiffOrigin == "" {
		songDiffOrigin = "x,x,x,x"
	}
	return strings.Split(songDiffOrigin, ",")
}

func GeneralShortenName(name string, limit int) string {
	var nameShort string
	if len(name) > limit {
		byteCount := 0
		for _, r := range name {
			runeLen := utf8.RuneLen(r)
			if byteCount+runeLen > limit {
				nameShort = name[:byteCount]
				break
			}
			byteCount += runeLen
			if byteCount == limit {
				nameShort = name[:byteCount]
			}
		}
	} else {
		nameShort = name
	}
	return nameShort
}
