package SongDataUpdater

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

var songNoteDensityIndex = [13]float64{0.0, 0.8, 1.2, 1.4, 1.5, 1.9, 2.3, 2.6, 3.2, 3.8, 4.4, 6.0, 8.4}
var songScoreDensityIndex = [13]float64{0, 320, 450, 600, 800, 1100, 1300, 1700, 2000, 2500, 3000, 3500, 5500}
var songLengthAvg = 135.0

func getSongStaticValueData(db *sql.DB) []songValueCalcUnits {
	var songValueCalBases []songValueCalcUnits

	SQLCode := `SELECT ms.music_album,ms.music_album_number,music_diff,music_total_note,music_total_score,ml.music_audio_length
FROM mdsong ms
JOIN mdsonglength ml
ON ms.music_album = ml.music_album AND ms.music_album_number = ml.music_album_number`
	result, err := db.Query(SQLCode)
	for result.Next() {
		var songValueCalBase songValueCalcUnits
		dataBuf := make([]string, 3)
		err = result.Scan(&songValueCalBase.SongCode, &songValueCalBase.AlbumCode, &dataBuf[0], &dataBuf[1], &dataBuf[2], &songValueCalBase.SongLengthData)
		songValueCalBase.SongDiff = strings.Split(dataBuf[0], ",")
		songValueCalBase.SongNotes = strings.Split(dataBuf[1], ",")
		songValueCalBase.SongScore = strings.Split(dataBuf[2], ",")
		if err != nil {
			log.Fatal(err)
		}
		songValueCalBases = append(songValueCalBases, songValueCalBase)
	}
	return songValueCalBases
}

func SongValueTableInit(db *sql.DB) {
	play_table_sql := `
	CREATE TABLE  IF NOT EXISTS mdvalue(
		music_album INTEGER NOT NULL,
		music_album_number INTEGER NOT NULL,
		music_diff_tier INTEGER NOT NULL,
		music_value REAL,
		PRIMARY KEY (music_album, music_album_number,music_diff_tier)
		FOREIGN KEY (music_album, music_album_number) REFERENCES mdsong(music_album, music_album_number)
	);
    `
	_, err := db.Exec(play_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func insertValueData(db *sql.DB, valueMap songValueMap) error {
	queryMDData := `INSERT INTO 
	mdvalue (music_album,music_album_number,music_diff_tier,music_value) 
	VALUES (?, ?, ?, ?)
	ON CONFLICT(music_album,music_album_number,music_diff_tier)
	DO UPDATE SET music_value=excluded.music_value`
	statement, err := db.Prepare(queryMDData)
	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec(valueMap.AlbumCode, valueMap.SongCode, valueMap.SongDiffTier, valueMap.SongValue)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}
	return nil
}

func CalStaticValue(outputData bool) {
	rowDatas := getSongStaticValueData(DBConnector())
	valueList := make([][]float64, 13)
	for _, rowData := range rowDatas {
		fmt.Println(rowData)
		for i := 0; i <= 3; i++ {
			songDiffBase, err := strconv.Atoi(rowData.SongDiff[i])
			if err != nil {
				insertValueData(DBConnector(), songValueMap{rowData.AlbumCode, rowData.SongCode, i, -1})
				continue
			}
			if songDiffBase == 0 {
				insertValueData(DBConnector(), songValueMap{rowData.AlbumCode, rowData.SongCode, i, 0})
				continue
			}

			songNoteBase, err := strconv.ParseFloat(rowData.SongNotes[i], 64)
			songScoreBase, err := strconv.ParseFloat(rowData.SongScore[i], 64)
			if err != nil {
				insertValueData(DBConnector(), songValueMap{rowData.AlbumCode, rowData.SongCode, i, 0})
				continue
			}

			songNoteIndex := max((songNoteBase/rowData.SongLengthData/songNoteDensityIndex[songDiffBase]-1.0)*2.5, 0)
			songScoreIndex := max((songScoreBase/rowData.SongLengthData/songScoreDensityIndex[songDiffBase]-1.0)*1.5, 0)
			songLengthIndex := rowData.SongLengthData/songLengthAvg - 1.0

			songValue := math.Round((float64(songDiffBase)+min(max((songNoteIndex+songScoreIndex)/2, songLengthIndex), 0.9))*10) / 10

			insertValueData(DBConnector(), songValueMap{rowData.AlbumCode, rowData.SongCode, i, songValue})
			valueList[songDiffBase] = append(valueList[songDiffBase], songValue)
		}
	}

	if outputData {
		outputCalData(valueList)
	}
}

func ImportSongValueFromJsonFile() {
	valueJsonFiles := getJsonRawFile("song_value.json").(map[string]interface{})
	for key, valueJsonFile := range valueJsonFiles {
		fmt.Println(key, valueJsonFile)
		var valueMap songValueMap
		var err error
		valueMap.AlbumCode, err = strconv.Atoi(strings.Split(key, "-")[0])
		valueMap.SongCode, err = strconv.Atoi(strings.Split(key, "-")[1])
		if err != nil {
			panic(err)
		}
		for songDiff := 0; songDiff < len(valueJsonFile.([]interface{})); songDiff++ {
			valueMap.SongDiffTier = songDiff
			if !strings.ContainsAny(valueJsonFile.([]interface{})[songDiff].(string), ".") {
				continue
			}
			valueMap.SongValue, err = strconv.ParseFloat(valueJsonFile.([]interface{})[songDiff].(string), 64)
			if err != nil {
				panic(err)
			}
			insertValueData(DBConnector(), valueMap)
		}
	}
}

func outputCalData(valueList [][]float64) {
	for i := 1; i <= 12; i++ {
		sort.Slice(valueList[i], func(j, k int) bool {
			return valueList[i][j] < valueList[i][k]
		})
	}
	fmt.Println(valueList[1:])
}

func GetSongValueBySongDiff(db *sql.DB, albumCode, songCode, diffTier int) float64 {
	var songValue float64
	getSQLValueNum := `SELECT music_value FROM mdvalue WHERE music_album=? AND music_album_number=? AND music_diff_tier=?`
	result := db.QueryRow(getSQLValueNum, albumCode, songCode, diffTier)
	result.Scan(&songValue)
	return songValue
}

func calRKS(db *sql.DB, albumCode, songCode, diffTier, rank int, acc float64, enable bool) float64 {
	if !enable {
		return 0
	}
	baseValue := GetSongValueBySongDiff(db, albumCode, songCode, diffTier-1)
	rankIndex := max(0, (800-float64(rank))*0.0001) + 1
	rks := baseValue * acc * 0.01 * rankIndex
	return rks
}
