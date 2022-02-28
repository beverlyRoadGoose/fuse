package telegram // import "heytobi.dev/fuse/telegram"

type setWebhookResponse struct {
	Ok          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}

// SendMessageRequest defines messages that can be sent by the bot.
// https://core.telegram.org/bots/api#sendmessage
type SendMessageRequest struct {
	ChatID                   int64           `json:"chat_id"`
	Text                     string          `json:"text"`
	ParseMode                string          `json:"parse_mode"`
	Entities                 []MessageEntity `json:"entities"`
	DisableWebPagePreview    bool            `json:"disable_web_page_preview"`
	DisableNotification      bool            `json:"disable_notification"`
	ProtectContent           bool            `json:"protect_content"`
	ReplyToMessageID         int             `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool            `json:"allow_sending_without_reply"`
	// TODO add missing optional reply_markup
}

type sendMessageResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int    `json:"message_id"`
		From      *User  `json:"from"`
		Chat      *Chat  `json:"chat"`
		Date      int    `json:"date"`
		Text      string `json:"text"`
	} `json:"result"`
}

// Webhook defines an endpoint for receiving telegram updates.
// See https://core.telegram.org/bots/api#setwebhook
type Webhook struct {
	Url            string   `json:"url"`
	IPAddress      string   `json:"ip_address"`
	MaxConnections int      `json:"max_connections"`
	AllowedUpdates []string `json:"allowed_updates"`
}
