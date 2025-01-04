package MDWebPageUtils

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
)

type SongUserSearch struct {
	UserNum     int                    `json:"user_num"`
	UserDetails []SongUserSearchDetail `json:"user_details"`
}

type SongUserSearchDetail struct {
	UserName      string  `json:"user_name"`
	UserId        string  `json:"user_id"`
	UserRKS       float64 `json:"user_rks"`
	UserRKSSimple float64 `json:"-"`
	UserRKSString string  `json:"-"`
}

func GetRKSString(RKSValueSimple float64) string {
	if RKSValueSimple < 0 {
		RKSValueString := "N/A"
		return RKSValueString
	}
	RKSValueString := strconv.FormatFloat(RKSValueSimple, 'f', -1, 64)
	return RKSValueString
}

func GetSongUserSearchResult(db *sql.DB, userInput string) SongUserSearch {
	var SongUserSearch SongUserSearch
	var SongUserSearchList []SongUserSearchDetail
	var SQLCode string
	if len(userInput) < 11 {
		SQLCode = `SELECT user_id,nickname FROM mduser WHERE nickname LIKE "%` + userInput + `%" AND LENGTH(nickname)<11 LIMIT 10`
	} else {
		SQLCode = `SELECT user_id,nickname FROM mduser WHERE nickname LIKE "玩家` + userInput + `" OR user_id is "` + userInput + `" LIMIT 10`
	}
	result, err := db.Query(SQLCode)
	if err != nil {
		fmt.Println(err)
	}
	for result.Next() {
		var SongUserSearchDetail SongUserSearchDetail
		err = result.Scan(&SongUserSearchDetail.UserId, &SongUserSearchDetail.UserName)
		if err != nil {
			fmt.Println(err)
		}
		SongUserSearchDetail.UserRKS = GetTotalRKS(db, SongUserSearchDetail.UserId, 30)
		SongUserSearchDetail.UserRKSSimple = math.Floor(SongUserSearchDetail.UserRKS*100) / 100
		SongUserSearchDetail.UserRKSString = GetRKSString(SongUserSearchDetail.UserRKSSimple)
		SongUserSearchList = append(SongUserSearchList, SongUserSearchDetail)
	}
	SongUserSearch.UserNum = len(SongUserSearchList)
	SongUserSearch.UserDetails = SongUserSearchList
	return SongUserSearch
}
