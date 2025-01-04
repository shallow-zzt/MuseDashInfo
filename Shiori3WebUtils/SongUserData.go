package MDWebPageUtils

import (
	"database/sql"
	"db/SongDataMapping"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

type SongData struct {
	AlbumCode      int
	SongCode       int
	DiffTierCode   int
	DiffTier       string
	RKSValue       float64
	RKSValueSimple float64
	RKSValueString string
	Playtime       string
	PlaytimeStamp  int
	PlayTimeBefore string
	PlatformCode   int
	Platform       string
	DiffValue      string
	SongPic        string
	SongName       string
	SongNameShort  string
	SongAcc        float64
	SongAccBig     string
	SongAccSmall   string
	SongScore      int
	SongRank       int
	SongChara      int
	SongCharaName  string
	SongElfin      int
	SongElfinName  string
	IsUnknown      int
}

type UserSongInfo struct {
	UserName            string
	UserId              string
	TotalRKSValue       float64
	TotalRKSValueSimple float64
	UserSongInfoList    []SongData
}

func (c *SongData) ConvertPlayTime() {
	timeBeforeRuneCount := 0
	timeFormat, err := time.Parse(time.RFC3339, c.Playtime)
	if err != nil {
		return
	}
	c.PlaytimeStamp = int(timeFormat.Unix())
	playTimeStampTemp := int(time.Now().Unix()) - c.PlaytimeStamp
	if int(playTimeStampTemp/(86400*365)) >= 1 && timeBeforeRuneCount <= 2 {
		c.PlayTimeBefore += (strconv.Itoa(int(playTimeStampTemp/(86400*365))) + "y")
		playTimeStampTemp = playTimeStampTemp % (86400 / 365)
		timeBeforeRuneCount++
	}
	if int(playTimeStampTemp/(86400)) >= 1 && timeBeforeRuneCount <= 2 {
		c.PlayTimeBefore += (strconv.Itoa(int(playTimeStampTemp/(86400))) + "d")
		playTimeStampTemp = playTimeStampTemp % (86400)
		timeBeforeRuneCount++
	}
	if int(playTimeStampTemp/(3600)) >= 1 && timeBeforeRuneCount <= 2 {
		c.PlayTimeBefore += (strconv.Itoa(int(playTimeStampTemp/(3600))) + "h")
		playTimeStampTemp = playTimeStampTemp % (3600)
		timeBeforeRuneCount++
	}
	if int(playTimeStampTemp/(60)) <= 0 && timeBeforeRuneCount == 0 {
		c.PlayTimeBefore += "<1m"
	}
}

func (c *SongData) ConvertCode2Name() {
	c.DiffTier = SongDataMapping.DiffTierMap[c.DiffTierCode]
	c.Platform = SongDataMapping.PlatFormMap[c.PlatformCode]
	c.SongCharaName = SongDataMapping.CharaNameMap[c.SongChara]
	c.SongElfinName = SongDataMapping.ElfinNameMap[c.SongElfin]
}

func (c *SongData) ShortenSongName(limit int) {
	if len(c.SongName) > limit {
		byteCount := 0
		for _, r := range c.SongName {
			runeLen := utf8.RuneLen(r)
			if byteCount+runeLen > limit {
				c.SongNameShort = c.SongName[:byteCount] + "…"
				break
			}
			byteCount += runeLen
			if byteCount == limit {
				c.SongNameShort = c.SongName[:byteCount] + "…"
			}
		}
	} else {
		c.SongNameShort = c.SongName
	}
}

func (c *SongData) ConvertSongAcc() {
	totalSongAcc := strconv.Itoa(int(c.SongAcc * 100))
	c.SongAccBig = totalSongAcc[:len(totalSongAcc)-2]
	c.SongAccSmall = totalSongAcc[len(totalSongAcc)-2:]
}

func (c *SongData) GetSimpleRKS() {
	if c.RKSValueSimple <= 0 {
		c.RKSValueString = "N/A"
		return
	}
	c.RKSValueString = strconv.FormatFloat(c.RKSValueSimple, 'f', -1, 64)
}

func GetUserSongList(db *sql.DB, songdb *sql.DB, userID string, bestSongs int, bestSongsOffset int) UserSongInfo {
	var userSongInfo UserSongInfo
	var sqlTemp interface{}
	userSongInfo.UserId = userID
	userSongInfo.UserName = GetUserName(db, userID)
	userSongInfo.TotalRKSValue = GetTotalRKS(db, userID, 30)
	userSongInfo.TotalRKSValueSimple = math.Floor(userSongInfo.TotalRKSValue*100) / 100
	var songInfoList []SongData
	SQLCode := `SELECT * FROM mdplay WHERE user_id=? ORDER BY rks DESC LIMIT ? OFFSET ?`
	result, err := db.Query(SQLCode, userID, bestSongs, bestSongsOffset)
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var songData SongData
		err = result.Scan(&songData.PlatformCode, &songData.AlbumCode, &songData.SongCode, &songData.DiffTierCode, &songData.SongRank, &sqlTemp, &songData.SongAcc, &sqlTemp, &songData.SongScore, &sqlTemp, &songData.SongChara, &songData.SongElfin, &songData.Playtime, &songData.RKSValue)
		if err != nil {
			log.Fatal(err)
		}
		songData.IsUnknown = 0
		songData.RKSValueSimple = math.Floor(songData.RKSValue*100) / 100
		songData.SongPic = GetSongPicFromCode(songdb, songData.AlbumCode, songData.SongCode)
		songData.DiffValue = GetSongDiffFromCode(songdb, songData.AlbumCode, songData.SongCode)[songData.DiffTierCode-1]
		songData.SongName = GetSongNameFromCode(songdb, songData.AlbumCode, songData.SongCode)
		songData.GetSimpleRKS()
		songData.ConvertCode2Name()
		songData.ConvertPlayTime()
		songData.ConvertSongAcc()
		songData.ShortenSongName(9)
		if songData.SongPic == "random_song_cover" {
			songData.IsUnknown = 1
		}
		songInfoList = append(songInfoList, songData)

	}
	userSongInfo.UserSongInfoList = songInfoList
	return userSongInfo
}
