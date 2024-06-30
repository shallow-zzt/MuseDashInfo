package MDWebPageUtils

import (
	"db/SongDataUpdater"
	"fmt"
	"testing"
)

func TestAliasCreateSQL(t *testing.T) {
	SongDataUpdater.GetBasicAilas()

}

func TestGetSongInfo(t *testing.T) {
	var db = SongDataUpdater.DBConnector()
	fmt.Println(GetSongInfoFromCode(db, 0, 0).GetAlias())

}
