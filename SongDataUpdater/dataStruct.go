package SongDataUpdater

type AliasBasicData struct {
	MusicAlbum       int    `json:"album-code"`
	MusicAlbumNumber int    `json:"song-code"`
	MusicAlias       string `json:"input-alias"`
}

type AlbumData struct {
	albumCode    string
	albumNameCN  string
	albumNameCNT string
	albumNameEN  string
	albumNameJA  string
	albumNameKR  string
}

type SongData struct {
	musicAlbum       int
	musicAlbumNumber int
	musicName        string
	musicAuthor      string
	musicBPM         string
	musicPicName     string
	musicSceneName   string
	musicSheetAuthor []string
	musicDiff        []string
	musicSpecialDiff string
	musicNotes       []string
	musicScore       []string
}

type songValueCalcUnits struct {
	SongCode       int
	AlbumCode      int
	SongLengthData float64
	SongDiff       []string
	SongNotes      []string
	SongScore      []string
}

type songValueMap struct {
	SongCode     int
	AlbumCode    int
	SongDiffTier int
	SongValue    float64
}

type songRankData struct {
	Rank       int
	UserId     string
	UserName   string
	Score      int
	Acc        float64
	CharaElfin string
}

type songRankInfo struct {
	SongCode     int
	AlbumCode    int
	SongName     string
	SongPic      string
	SongAuthor   string
	SongDiff     []string
	SongRankData []songRankData
}
