package MDWebPageUtils

import (
	"database/sql"
	"strings"
)

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
	return songPicName
}

func GetSongNameFromCode(db *sql.DB, albumCode int, songCode int) string {
	var songName string
	SQLCode := `SELECT music_name FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&songName)
	return songName
}

func GetSongDiffFromCode(db *sql.DB, albumCode int, songCode int) []string {
	var songDiffOrigin string
	SQLCode := `SELECT music_diff FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	result := db.QueryRow(SQLCode, albumCode, songCode)
	result.Scan(&songDiffOrigin)
	return strings.Split(songDiffOrigin, ",")
}
