package MDWebPageUtils

import (
	"database/sql"
	"log"
	"math"
)

type FullSongValueInfo struct {
	AlbumCode       int
	SongCode        int
	SongPic         string
	SongName        string
	SongValueEasy   float64
	SongValueHard   float64
	SongValueMaster float64
	SongValueHidden float64
}

func GetAllSongValueInfo(db *sql.DB, diffMode ...int) []FullSongValueInfo {
	var valueInfos []FullSongValueInfo
	SQLCode := `SELECT music_album,music_album_number,music_pic_name,music_name FROM mdsong`
	result, err := db.Query(SQLCode)
	for result.Next() {
		var valueInfo FullSongValueInfo
		valueAddress := [4]*float64{&valueInfo.SongValueEasy, &valueInfo.SongValueHard, &valueInfo.SongValueMaster, &valueInfo.SongValueHidden}
		err = result.Scan(&valueInfo.AlbumCode, &valueInfo.SongCode, &valueInfo.SongPic, &valueInfo.SongName)
		if err != nil {
			log.Fatal(err)
		}
		getSQLValueNum := `SELECT music_value FROM mdvalue WHERE music_album=? AND music_album_number=?`
		result2, err := db.Query(getSQLValueNum, valueInfo.AlbumCode, valueInfo.SongCode)
		if err != nil {
			log.Fatal(err)
		}
		counter := 0
		for result2.Next() {
			result2.Scan(valueAddress[counter])
			if len(diffMode) != 0 && diffMode[0] == 1 {
				*valueAddress[counter] = math.Floor(*valueAddress[counter])
			}
			counter += 1
		}
		valueInfos = append(valueInfos, valueInfo)
	}
	return valueInfos
}
