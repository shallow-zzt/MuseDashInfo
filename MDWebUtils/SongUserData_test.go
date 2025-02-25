package MDWebPageUtils

import (
	"db/SongDataUpdater"
	"fmt"
	"testing"
)

func TestLoadingDatabase(t *testing.T) {
	var db = SongDataUpdater.DBConnector("../mdsong.db")
	SQLCode := `SELECT * FROM mdsonglength`
	result, err := db.Query(SQLCode)
	fmt.Println(result)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetUserSongInfo(t *testing.T) {
	var db = SongDataUpdater.DBConnector("../MDRankData.db")
	var songdb = SongDataUpdater.DBConnector("../mdsong.db")
	fmt.Println(GetUserSongList(db, songdb, "6bb31610695511eb962d0242ac11005e", 200, 0))
}

func TestGetSearchUser(t *testing.T) {
	var db = SongDataUpdater.DBConnector("../MDRankData.db")
	fmt.Println(GetSongUserSearchResult(db, "檐下汐", 10, 0))
}
