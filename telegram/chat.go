package telegram // import "heytobi.dev/fuse/telegram"

// ChatPhoto represents a chat photo.
// See https://core.telegram.org/bots/api#chatphoto
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatLocation represents a location to which a chat is connected.
// See https://core.telegram.org/bots/api#chatlocation
type ChatLocation struct {
	Location *Location `json:"location"`
	Address  string    `json:"address"`
}

// ChatPermissions describes actions that a non-administrator user is allowed to take in a chat.
// See https://core.telegram.org/bots/api#chatpermissions
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

// Chat represents a chat.
// See https://core.telegram.org/bots/api#chat
type Chat struct {
	ID                    int64            `json:"id"`
	Type                  string           `json:"type"`
	Title                 string           `json:"title"`
	Username              string           `json:"username"`
	FirstName             string           `json:"first_name"`
	LastName              string           `json:"last_name"`
	Photo                 *ChatPhoto       `json:"photo"`
	Bio                   string           `json:"bio"`
	HasPrivateForwards    bool             `json:"has_private_forwards"`
	Description           string           `json:"description"`
	InviteLink            string           `json:"invite_link"`
	PinnedMessage         *Message         `json:"pinned_message"`
	Permissions           *ChatPermissions `json:"permissions"`
	SlowModeDelay         int              `json:"slow_mode_delay"`
	MessageAutoDeleteTime int              `json:"message_auto_delete_time"`
	HasProtectedContent   bool             `json:"has_protected_content"`
	StickerSetName        string           `json:"sticker_set_name"`
	CanSetStickerSet      bool             `json:"can_set_sticker_set"`
	LinkedChatID          int64            `json:"linked_chat_id"`
	Location              *ChatLocation    `json:"location"`
}
