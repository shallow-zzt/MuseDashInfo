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
		err = result.Scan(&songValueCalBase.AlbumCode, &songValueCalBase.SongCode, &dataBuf[0], &dataBuf[1], &dataBuf[2], &songValueCalBase.SongLengthData)
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
		music_value_easy REAL,
		music_value_hard REAL,
		music_value_master REAL,
		music_value_hidden REAL,
		PRIMARY KEY (music_album, music_album_number)
		FOREIGN KEY (music_album, music_album_number) REFERENCES mdsong(music_album, music_album_number)
	);
    `
	_, err := db.Exec(play_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func CalStaticValue() {
	rowDatas := getSongStaticValueData(DBConnector())
	valueList := make([][]float64, 13)
	for _, rowData := range rowDatas {
		for i := 0; i <= 3; i++ {
			songDiffBase, err := strconv.Atoi(rowData.SongDiff[i])
			if songDiffBase == 0 || err != nil {
				continue
			}

			songNoteBase, err := strconv.ParseFloat(rowData.SongNotes[i], 64)
			songScoreBase, err := strconv.ParseFloat(rowData.SongScore[i], 64)
			if err != nil {
				continue
			}

			songNoteIndex := max((songNoteBase/rowData.SongLengthData/songNoteDensityIndex[songDiffBase]-1.0)*2.5, 0)
			songScoreIndex := max((songScoreBase/rowData.SongLengthData/songScoreDensityIndex[songDiffBase]-1.0)*1.5, 0)
			songLengthIndex := rowData.SongLengthData/songLengthAvg - 1.0

			songValue := math.Round((float64(songDiffBase)+min(max((songNoteIndex+songScoreIndex)/2, songLengthIndex), 0.9))*10) / 10
			valueList[songDiffBase] = append(valueList[songDiffBase], songValue)
		}
	}
}

func outputCalData(valueList [][]float64) {
	for i := 1; i <= 12; i++ {
		sort.Slice(valueList[i], func(j, k int) bool {
			return valueList[i][j] < valueList[i][k]
		})
	}
	fmt.Println(valueList)
}
