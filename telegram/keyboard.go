package telegram // import "heytobi.dev/fuse/telegram"

// InlineKeyboardButton represents one button of an inline keyboard.
// See https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          string        `json:"url"`
	LoginUrl                     *LoginUrl     `json:"login_url"`
	CallbackData                 string        `json:"callback_data"`
	SwitchInlineQuery            string        `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string        `json:"switch_inline_query_current_chat"`
	CallbackGame                 *CallbackGame `json:"callback_game"`
	Pay                          bool          `json:"pay"`
}

// InlineKeyboardMarkup represents an inline keyboard that appears right next to the message it belongs to.
// See https://core.telegram.org/bots/api#inlinekeyboardmarkup
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// LoginUrl represents a parameter of the inline keyboard button used to automatically authorize a user.
// See https://core.telegram.org/bots/api#loginurl
type LoginUrl struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text"`
	BotUsername        string `json:"bot_username"`
	RequestWriteAccess bool   `json:"request_write_access"`
}

// CallbackGame A placeholder, currently holds no information.
// See https://core.telegram.org/bots/api#callbackgame
type CallbackGame struct {
}
