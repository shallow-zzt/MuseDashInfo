package SongDataUpdater

import (
	"db/SongDataMapping"
	"math"
	"strings"
)

func charaElfinNameConvert(code int, codeType int) string {
	if codeType == 0 {
		if code < 0 || code >= len(SongDataMapping.CharaNameMap) {
			return "未知人物"
		}
		return SongDataMapping.CharaNameMap[code]
	} else if codeType == 1 {
		if code < 0 || code >= len(SongDataMapping.ElfinNameMap) {
			return "未知精灵"
		}
		return SongDataMapping.ElfinNameMap[code]
	}
	return "未知"
}

func GetSongRankData(albumCode, songCode, diffTier, platform int) songRankInfo {
	var songRankInfo songRankInfo
	var songDiffBuf string
	var songChartAuthorBuf string
	rankDB := DBConnector("MDRankData.db")
	songDB := DBConnector()

	SongInfoCode := `SELECT music_pic_name,music_name,music_author,music_diff,music_sheet_author FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	songInfoResult := songDB.QueryRow(SongInfoCode, albumCode, songCode)
	songInfoResult.Scan(&songRankInfo.SongPic, &songRankInfo.SongName, &songRankInfo.SongAuthor, &songDiffBuf, &songChartAuthorBuf)
	songRankInfo.SongDiff = strings.Split(songDiffBuf, ",")
	songRankInfo.SongChartAuthor = strings.Split(songChartAuthorBuf, "|")

	RankInfoCode := `SELECT rank,user_id,score,acc,character_uid,elfin_uid FROM mdplay WHERE music_album_id = ? AND music_album_song_id = ? AND music_difficulty = ? AND platform = ?`
	RankInfoResult, err := rankDB.Query(RankInfoCode, albumCode, songCode, diffTier, platform)
	if err != nil {
		return songRankInfo
	}
	for RankInfoResult.Next() {
		var songRankData songRankData
		var userNameBuf string
		var charaBuf, elfinBuf int
		RankInfoResult.Scan(&songRankData.Rank, &userNameBuf, &songRankData.Score, &songRankData.Acc, &charaBuf, &elfinBuf)
		songRankData.UserId = userNameBuf
		songRankData.Acc = math.Round(songRankData.Acc*100) / 100
		getUserNameCode := `SELECT nickname FROM mduser WHERE user_id=?`
		getUserNameResult := rankDB.QueryRow(getUserNameCode, userNameBuf)
		getUserNameResult.Scan(&songRankData.UserName)
		songRankData.CharaElfin = charaElfinNameConvert(charaBuf, 0) + " / " + charaElfinNameConvert(elfinBuf, 1)
		songRankInfo.SongRankData = append(songRankInfo.SongRankData, songRankData)
	}
	return songRankInfo
}
