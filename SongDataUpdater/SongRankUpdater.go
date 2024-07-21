package SongDataUpdater

import "strings"

func charaElfinNameConvert(code int, codeType int) string {
	charaNameMap := []string{"贝斯手", "问题少女", "梦游少女", "兔女郎", "飞行员", "偶像", "僵尸少女", "华服小丑", "提琴少女", "女仆", "魔法少女", "小恶魔", "黑衣少女", "圣诞礼物", "制服少女", "领航员", "游戏主播", "打工战士", "博丽灵梦", "重生的少女", "修女", "雾雨魔理沙", "阿米娅", "拳击手", "道士", "初音未来", "镜音铃·连", "摩托车手"}
	elfinNameMap := []string{"喵斯", "安吉拉", "塔纳托斯", "Rabot-233", "小护士", "小女巫", "小龙女", "莉莉丝", "佩奇医生", "静音灵", "霓虹彩蛋", "Beta狗"}
	if codeType == 0 {
		if code >= len(charaNameMap) {
			return "未知人物"
		}
		return charaNameMap[code]
	} else if codeType == 1 {
		if code >= len(elfinNameMap) {
			return "未知精灵"
		}
		return elfinNameMap[code]
	}
	return "未知"
}

func GetSongRankData(albumCode, songCode, diffTier, platform int) songRankInfo {
	var songRankInfo songRankInfo
	var songDiffBuf string
	rankDB := DBConnector("MDRankData.db")
	songDB := DBConnector()

	SongInfoCode := `SELECT music_pic_name,music_name,music_author,music_diff FROM mdsong WHERE music_album = ? AND music_album_number = ?`
	songInfoResult := songDB.QueryRow(SongInfoCode, albumCode, songCode)
	songInfoResult.Scan(&songRankInfo.SongPic, &songRankInfo.SongName, &songRankInfo.SongAuthor, &songDiffBuf)
	songRankInfo.SongDiff = strings.Split(songDiffBuf, ",")

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
		getUserNameCode := `SELECT nickname FROM mduser WHERE user_id=?`
		getUserNameResult := rankDB.QueryRow(getUserNameCode, userNameBuf)
		getUserNameResult.Scan(&songRankData.UserName)
		songRankData.CharaElfin = charaElfinNameConvert(charaBuf, 0) + " / " + charaElfinNameConvert(elfinBuf, 1)
		songRankInfo.SongRankData = append(songRankInfo.SongRankData, songRankData)
	}
	return songRankInfo
}
