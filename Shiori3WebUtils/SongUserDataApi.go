package MDWebPageUtils

import (
	"database/sql"
	"strings"
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
