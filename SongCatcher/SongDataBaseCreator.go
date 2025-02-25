package SongCatcher

import (
	"database/sql"
	"fmt"
)

func createTable(db *sql.DB) {
	createPlayTable(db)
	createUserTable(db)
	createRksHistoryTable(db)
}

func createPlayTable(db *sql.DB) {
	play_table_sql := `
		CREATE TABLE IF NOT EXISTS ` + PLAY_TABLE + ` (
			platform INTEGER NOT NULL,
			music_album_id INTEGER NOT NULL,
			music_album_song_id INTEGER NOT NULL,
			music_difficulty INTEGER NOT NULL,
			rank INTEGER NOT NULL,
			user_id TEXT,
			acc REAL,
			miss INTEGER,
			score INTEGER,
			judge TEXT,
			character_uid INTEGER,
			elfin_uid INTEGER,
			updated_at TEXT,
			rks REAL,
			PRIMARY KEY (platform, music_album_id,music_album_song_id, music_difficulty, rank)
		);
	`
	_, err := db.Exec(play_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func createUserTable(db *sql.DB) {
	user_table_sql := `
		CREATE TABLE IF NOT EXISTS ` + USER_TABLE + ` (
			user_id TEXT PRIMARY KEY NOT NULL,
			nickname TEXT,
			user_rks REAL
		);
	`
	_, err := db.Exec(user_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}

func createRksHistoryTable(db *sql.DB) {
	rks_history_table_sql := `
		CREATE TABLE IF NOT EXISTS ` + RKS_HISTORY_TABLE + ` (
			user_id TEXT PRIMARY KEY NOT NULL,
			rks_list TEXT,
			updated_at_list TEXT
		);
	`
	_, err := db.Exec(rks_history_table_sql)
	if err != nil {
		fmt.Println("创建表失败:", err)
	}
}
