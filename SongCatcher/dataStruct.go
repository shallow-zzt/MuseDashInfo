package SongCatcher

type PlayData struct {
	acc              float64
	bms_id           int
	bms_version      string
	character_uid    int
	combo            int
	created_at       string
	elfin_uid        int
	game_version     string
	hp               int
	is_check         bool
	judge            string
	lb_version       string
	miss             int
	music_difficulty int
	music_uid        string
	object_id        string
	score            int
	updated_at       string
	user_id          string
	visible          bool
}

type UserData struct {
	created_at string
	nickname   string
	object_id  string
	updated_at string
	user_id    string
}

type MDdbData struct {
	platform         int
	musicAlbum       int
	musicAlbumNumber int
	musicDiff        int
	rank             int
	userId           string
	acc              float64
	miss             int
	score            int
	judge            string
	characterId      int
	elfinId          int
	playTime         string
	nickname         string
	rks              float64
}
