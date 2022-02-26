package telegram // import "heytobi.dev/fuse/telegram"

// VoiceChatScheduled represents a service message about a voice chat scheduled in the chat.
// See https://core.telegram.org/bots/api#voicechatscheduled
type VoiceChatScheduled struct {
	StartDate int `json:"start_date"`
}

// VoiceChatStarted represents a service message about a voice chat started in the chat. It currently holds no information.
// See https://core.telegram.org/bots/api#voicechatstarted
type VoiceChatStarted struct {
}

// VoiceChatEnded represents a service message about a voice chat ended in the chat.
// See https://core.telegram.org/bots/api#voicechatended
type VoiceChatEnded struct {
	Duration int `json:"duration"`
}

// VoiceChatParticipantsInvited represents a service message about new members invited to a voice chat.
// See https://core.telegram.org/bots/api#voicechatparticipantsinvited
type VoiceChatParticipantsInvited struct {
	Users []User `json:"users"`
}
