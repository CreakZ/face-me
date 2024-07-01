package constants

const (
	BaseURL      = "https://www.faceit.com/api"
	IDInfo       = "users/v1/nicknames"
	MatchInfoURL = "stats/v1/stats/users/%v/games/%v"                                    // %v represents 2 binvars: player's id and game (csgo or cs2) respectively
	MatchInfo    = "stats/v1/stats/time/users/%v/games/%v?page=%v&size=%v&game_mode=5v5" // %v represents 4 binvars: player's id, game (csgo or cs2), page number, page size respectively
	BaseSize     = 100
)
