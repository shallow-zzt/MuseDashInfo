package SongCatcher

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SetUserTotalRks(db *sql.DB) {
	db.Exec("PRAGMA journal_mode=WAL;")
	getUserSQLCode := `SELECT user_id,nickname FROM mduser`
	setRKSSQLCode := `INSERT OR REPLACE INTO 
	mduser (user_id,nickname,user_rks) 
	VALUES (?, ?, ?)
	ON CONFLICT(user_id)
	DO UPDATE SET user_rks=excluded.user_rks`
	result, err := db.Query(getUserSQLCode)
	if err != nil {
		fmt.Println("查询失败:", err)
	}

	var userTotalRksData UserRksData
	for result.Next() {
		result.Scan(&userTotalRksData.userId, &userTotalRksData.nickname)
		userTotalRksData.totalRks = getUserTotalRKS(db, userTotalRksData.userId, 30)
		statement, err := db.Prepare(setRKSSQLCode)
		if err != nil {
			fmt.Errorf("创建表失败: %v", err)
		}
		defer statement.Close()
		_, err = statement.Exec(userTotalRksData.userId, userTotalRksData.nickname, userTotalRksData.totalRks)
		if err != nil {
			fmt.Errorf("创建表失败: %v", err)
		}
		setUserTotalRksHistory(db, userTotalRksData.userId, userTotalRksData.totalRks)
	}
	defer result.Close()
}

func getUserTotalRKS(db *sql.DB, userID string, bestSongs int) float64 {
	var averageRKS float64
	SQLCode := `SELECT AVG(rks) AS rks_average
	FROM(
	SELECT rks FROM mdplay WHERE user_id=? ORDER BY rks DESC LIMIT ?
	)`
	result := db.QueryRow(SQLCode, userID, bestSongs)
	result.Scan(&averageRKS)
	return averageRKS
}

func (c *UserRksHistory) SetRksFloatFormat() {
	for i := 0; i < len(c.rksStrList); i++ {
		convertRKS, err := strconv.ParseFloat(c.rksStrList[i], 64)
		if err != nil {
			c.rksList = append(c.rksList, -1)
			continue
		}
		c.rksList = append(c.rksList, convertRKS)
	}
}

func setUserTotalRksHistory(db *sql.DB, userId string, rks float64) {
	db.Exec("PRAGMA journal_mode=WAL;")
	var userRksHistory UserRksHistory
	var rksListStr, timeListStr string
	getRksDataSQLCode := `SELECT rks_list,updated_at_list FROM mdrks WHERE user_id=?`
	setRksDataSQLCode := `INSERT OR REPLACE INTO 
	mdrks (user_id,rks_list,updated_at_list) 
	VALUES (?, ?, ?)
	ON CONFLICT(user_id)
	DO UPDATE SET rks_list=excluded.rks_list,updated_at_list=excluded.updated_at_list`
	result := db.QueryRow(getRksDataSQLCode, userId)
	result.Scan(&rksListStr, &timeListStr)
	userRksHistory.userId = userId
	userRksHistory.rksStrList = strings.Split(rksListStr, ",")
	userRksHistory.updatedAtList = strings.Split(timeListStr, ",")
	if len(userRksHistory.rksStrList) >= 30 {
		userRksHistory.rksStrList = append(userRksHistory.rksStrList[1:], strconv.FormatFloat(rks, 'f', -1, 64))
		userRksHistory.updatedAtList = append(userRksHistory.updatedAtList[1:], time.Now().Format(time.RFC3339))
	} else {
		userRksHistory.rksStrList = append(userRksHistory.rksStrList, strconv.FormatFloat(rks, 'f', -1, 64))
		userRksHistory.updatedAtList = append(userRksHistory.updatedAtList, time.Now().Format(time.RFC3339))
	}
	statement, err := db.Prepare(setRksDataSQLCode)
	if err != nil {
		fmt.Errorf("创建表失败: %v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(userRksHistory.userId, strings.Join(userRksHistory.rksStrList, ","), strings.Join(userRksHistory.updatedAtList, ","))
	if err != nil {
		fmt.Errorf("创建表失败: %v", err)
	}
	fmt.Println(userId, userRksHistory.rksStrList)
}
