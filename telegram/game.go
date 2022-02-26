package telegram // import "heytobi.dev/fuse/telegram"

// Game represents a game.
// See https://core.telegram.org/bots/api#dice
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Animation    *Animation      `json:"animation"`
}
