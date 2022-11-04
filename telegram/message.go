package telegram // import "heytobi.dev/fuse/telegram"

// Message represents a message.
// See https://core.telegram.org/bots/api#message
type Message struct {
	ID                            int64                          `json:"message_id"`
	Sender                        *User                          `json:"from"`
	SenderChat                    *Chat                          `json:"sender_chat"`
	Date                          int                            `json:"unix_date"`
	Chat                          *Chat                          `json:"chat"`
	ForwardedBy                   *User                          `json:"forward_from"`
	OriginalChat                  *Chat                          `json:"forwarded_from_chat"`
	OriginalMessageId             int64                          `json:"forward_from_message_id"`
	ForwardSignature              string                         `json:"forward_signature"`
	ForwardSenderName             string                         `json:"forward_sender_name"`
	ForwardDate                   int                            `json:"forward_date"`
	AutomaticallyForwarded        bool                           `json:"is_automatic_forward"`
	RepliedMessage                *Message                       `json:"reply_to_message"`
	Bot                           *User                          `json:"via_bot"`
	EditDate                      int                            `json:"edit_date"`
	HasProtectedContent           bool                           `json:"has_protected_content"`
	MediaGroupID                  string                         `json:"media_group_id"`
	AuthorSignature               string                         `json:"author_signature"`
	Text                          string                         `json:"text"`
	Entities                      []MessageEntity                `json:"entities"`
	Animation                     *Animation                     `json:"animation"`
	Audio                         *Audio                         `json:"audio"`
	Document                      *Document                      `json:"document"`
	Photo                         []PhotoSize                    `json:"photo"`
	Sticker                       *Sticker                       `json:"sticker"`
	Video                         *Video                         `json:"video"`
	VideoNote                     *VideoNote                     `json:"video_note"`
	Voice                         *Voice                         `json:"voice"`
	Caption                       string                         `json:"caption"`
	CaptionEntities               []MessageEntity                `json:"caption_entities"`
	Contact                       *Contact                       `json:"contact"`
	Dice                          *Dice                          `json:"dice"`
	Game                          *Game                          `json:"game"`
	Poll                          *Poll                          `json:"poll"`
	Venue                         *Venue                         `json:"venue"`
	Location                      *Location                      `json:"location"`
	NewChatMembers                []User                         `json:"new_chat_members"`
	LeftChatMember                *User                          `json:"left_chat_member"`
	NewChatTitle                  string                         `json:"new_chat_title"`
	NewChatPhoto                  []PhotoSize                    `json:"new_chat_photo"`
	ChatPhotoDeleted              bool                           `json:"chat_photo_deleted"`
	GroupChatCreated              bool                           `json:"group_chat_created"`
	SuperGroupChatCreated         bool                           `json:"super_group_chat_created"`
	ChannelChatCreated            bool                           `json:"channel_chat_created"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed"`
	MigrateToChatID               int64                          `json:"migrate_to_chat_id"`
	MigrateFromChatID             int64                          `json:"migrate_from_chat_id"`
	PinnedMessage                 *Message                       `json:"pinned_message"`
	Invoice                       *Invoice                       `json:"invoice"`
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment"`
	ConnectedWebsite              string                         `json:"connected_website"`
	PassportData                  *PassportData                  `json:"passport_data"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered"`
	VoiceChatScheduled            *VoiceChatScheduled            `json:"voice_chat_scheduled"`
	VoiceChatEnded                *VoiceChatEnded                `json:"voice_chat_ended"`
	VoiceChatParticipantsInvited  *VoiceChatParticipantsInvited  `json:"voice_chat_participants_invited"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup"`
}

// MessageEntity one special entity in a text message. For example, hashtags, usernames, URLs, etc.
// See https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url"`
	User     *User  `json:"user"`
	Language string `json:"language"`
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
	// ReplyMarkup supports different value types, the supported types are listed in the telegram api docs
	ReplyMarkup any `json:"reply_markup"`
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
