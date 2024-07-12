package SongDataUpdater

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var songNoteDensityIndex = [12]float64{0.9, 1.1, 1.2, 1.6, 2.0, 2.6, 2.8, 3.6, 4.5, 5.2, 6.4, 8.0}
var songScoreDensityIndex = [12]float64{350, 500, 500, 800, 1100, 1400, 1800, 2300, 2700, 3300, 4300, 5500}
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
	for _, rowData := range rowDatas {
		for i := 0; i <= 3; i++ {
			songDiffBase, err := strconv.Atoi(rowData.SongDiff[i])
			if songDiffBase == 0 {
				continue
			}

			songNoteBase, err := strconv.ParseFloat(rowData.SongNotes[i], 64)
			songScoreBase, err := strconv.ParseFloat(rowData.SongScore[i], 64)
			if err != nil {
				continue
			}
			fmt.Println(rowData.AlbumCode, "-", rowData.SongCode, "-", i, "难度", songDiffBase, "物量密度:", songNoteBase/rowData.SongLengthData, "分数密度:", songScoreBase/rowData.SongLengthData, "曲长:", rowData.SongLengthData)
		}
	}
}
