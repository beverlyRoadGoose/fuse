package telegram // import "heytobi.dev/fuse/telegram"

// BotCommand represents a bot command.
// See https://core.telegram.org/bots/api#botcommand
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// MaskPosition describes the position on faces where a mask should be placed by default.
// See https://core.telegram.org/bots/api#maskposition
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// Contact represents a phone contact.
// See https://core.telegram.org/bots/api#contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int64  `json:"user_id"`
	VCard       string `json:"vcard"`
}

// Dice represents an animated emoji that displays a random value.
// See https://core.telegram.org/bots/api#dice
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// MessageAutoDeleteTimerChanged represents a service message about a change in auto-delete timer settings.
// See https://core.telegram.org/bots/api#messageautodeletetimerchanged
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

// ProximityAlertTriggered represents the content of a service message, sent whenever a user in the chat triggers a
// proximity alert set by another user.
// See https://core.telegram.org/bots/api#proximityalerttriggered
type ProximityAlertTriggered struct {
	Traveler *User `json:"traveler"`
	Watcher  *User `json:"watcher"`
	Distance int   `json:"distance"`
}
