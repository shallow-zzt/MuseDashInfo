package MDWebPageUtils

import (
	"database/sql"
	"log"
	"math"
	"strconv"
	"strings"
)

type FullSongValueInfoWithNum struct {
	SongCount         int
	SongValueInfoList []FullSongValueInfo
}

type FullSongValueInfo struct {
	AlbumCode                int
	AlbumName                string
	AlbumNameShort           string
	SongCode                 int
	SongPic                  string
	SongName                 string
	SongValueEasy            float64
	SongValueHard            float64
	SongValueMaster          float64
	SongValueHidden          float64
	SongValueEasyString      string
	SongValueHardString      string
	SongValueMasterString    string
	SongValueHiddenString    string
	SongValueEasyHighlight   bool
	SongValueHardHighlight   bool
	SongValueMasterHighlight bool
	SongValueHiddenHighlight bool
}

func (c *FullSongValueInfo) SetAlbumName(db *sql.DB) {
	queryString := "ALBUM" + strconv.Itoa(c.AlbumCode+1)
	SQLCode := `SELECT music_album_name_CN FROM mdalbum WHERE music_album=?`
	result := db.QueryRow(SQLCode, queryString)
	result.Scan(&c.AlbumName)
	albumNameSplit := strings.Split(c.AlbumName, " ")
	if len(albumNameSplit[0]) > 5 {
		c.AlbumNameShort = GeneralShortenName(albumNameSplit[0], 18)
	} else {
		c.AlbumNameShort = GeneralShortenName(albumNameSplit[len(albumNameSplit)-1], 18)
	}
}

func GetAllSongValueInfo(db *sql.DB, diffMode ...int) FullSongValueInfoWithNum {
	var valueInfos []FullSongValueInfo
	SQLCode := `SELECT music_album,music_album_number,music_pic_name,music_name FROM mdsong`
	result, err := db.Query(SQLCode)
	for result.Next() {
		var valueInfo FullSongValueInfo
		valueAddress := [4]*float64{&valueInfo.SongValueEasy, &valueInfo.SongValueHard, &valueInfo.SongValueMaster, &valueInfo.SongValueHidden}
		valueHighlightAddress := [4]*bool{&valueInfo.SongValueEasyHighlight, &valueInfo.SongValueHardHighlight, &valueInfo.SongValueMasterHighlight, &valueInfo.SongValueHiddenHighlight}
		valueStringAddress := [4]*string{&valueInfo.SongValueEasyString, &valueInfo.SongValueHardString, &valueInfo.SongValueMasterString, &valueInfo.SongValueHiddenString}
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
			if *valueAddress[counter] == -1.0 {
				*valueStringAddress[counter] = GetSongDiffFromCode(db, valueInfo.AlbumCode, valueInfo.SongCode)[counter]
			}
			*valueHighlightAddress[counter] = true
			if len(diffMode) != 0 && diffMode[0] == 1 {
				*valueAddress[counter] = math.Floor(*valueAddress[counter])
			}
			counter += 1
		}
		valueInfo.SetAlbumName(db)
		valueInfos = append(valueInfos, valueInfo)
	}
	return FullSongValueInfoWithNum{SongCount: len(valueInfos), SongValueInfoList: valueInfos}
}

func GetPartialSongValueInfo(db *sql.DB, userInput string, floor float64, ceil float64, diffMode ...int) FullSongValueInfoWithNum {
	var valueInfos []FullSongValueInfo

	searchSQLCode := `SELECT DISTINCT music_album,music_album_number,music_pic_name,music_name
	FROM(SELECT ms.music_album,ms.music_album_number,ms.music_pic_name,ms.music_name,ma.music_alias
	FROM mdsong ms
	LEFT OUTER JOIN mdalias ma
	ON ms.music_album = ma.music_album AND ms.music_album_number = ma.music_album_number) msa
	WHERE msa.music_name LIKE "%` + userInput + `%" OR msa.music_alias LIKE "%` + userInput + `%"
	INTERSECT 
	SELECT DISTINCT mvs.music_album, mvs.music_album_number, ms.music_pic_name, ms.music_name
	FROM (SELECT DISTINCT music_album, music_album_number
		  FROM mdvalue
		  WHERE music_value >= ? AND music_value < ?) mvs
	LEFT OUTER JOIN mdsong ms
	ON ms.music_album = mvs.music_album AND ms.music_album_number = mvs.music_album_number;`
	result, err := db.Query(searchSQLCode, floor, ceil+1)
	if err != nil {
		log.Fatal(err)
	}
	for result.Next() {
		var valueInfo FullSongValueInfo
		valueAddress := [4]*float64{&valueInfo.SongValueEasy, &valueInfo.SongValueHard, &valueInfo.SongValueMaster, &valueInfo.SongValueHidden}
		valueHighlightAddress := [4]*bool{&valueInfo.SongValueEasyHighlight, &valueInfo.SongValueHardHighlight, &valueInfo.SongValueMasterHighlight, &valueInfo.SongValueHiddenHighlight}
		valueStringAddress := [4]*string{&valueInfo.SongValueEasyString, &valueInfo.SongValueHardString, &valueInfo.SongValueMasterString, &valueInfo.SongValueHiddenString}
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
			if *valueAddress[counter] == -1.0 {
				*valueStringAddress[counter] = GetSongDiffFromCode(db, valueInfo.AlbumCode, valueInfo.SongCode)[counter]
			}
			if *valueAddress[counter] >= floor && *valueAddress[counter] < ceil+1 {
				*valueHighlightAddress[counter] = true
			}
			if len(diffMode) != 0 && diffMode[0] == 1 {
				*valueAddress[counter] = math.Floor(*valueAddress[counter])
			}
			counter += 1
		}
		valueInfo.SetAlbumName(db)
		valueInfos = append(valueInfos, valueInfo)
	}
	return FullSongValueInfoWithNum{SongCount: len(valueInfos), SongValueInfoList: valueInfos}
}
