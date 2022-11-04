package telegram // import "heytobi.dev/fuse/telegram"

// InlineQuery represents an incoming inline query.
// See https://core.telegram.org/bots/api#inlinequery
type InlineQuery struct {
	ID       string    `json:"id"`
	Sender   *User     `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type"`
	Location *Location `json:"location"`
}

// ChosenInlineResult represents a result of an inline query that was chosen by the user and sent to their chat partner.
// See https://core.telegram.org/bots/api#choseninlineresult
type ChosenInlineResult struct {
	ID              string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// CallbackQuery represents an incoming callback query from a callback button in an inline keyboard.
// See https://core.telegram.org/bots/api#callbackquery
type CallbackQuery struct {
	ID              string   `json:"result_id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageID string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

// ShippingQuery contains information about an incoming shipping query.
// See https://core.telegram.org/bots/api#shippingquery
type ShippingQuery struct {
	ID              string           `json:"result_id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
// See https://core.telegram.org/bots/api#precheckoutquery
type PreCheckoutQuery struct {
	ID               string     `json:"result_id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id"`
	OrderInfo        *OrderInfo `json:"order_info"`
}

// Update represents an incoming update.
// See https://core.telegram.org/bots/api#update
type Update struct {
	ID                 int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
	Poll               *Poll               `json:"poll"`
	PollAnswer         *PollAnswer         `json:"poll_answer"`
	// TODO add my_chat_member, chat_member, chat_join_request
}

// See https://core.telegram.org/bots/api#getupdates
type getUpdatesRequest struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type getUpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}
