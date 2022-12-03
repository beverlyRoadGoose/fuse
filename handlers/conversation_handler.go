package handlers // import "heytobi.dev/fuse/handlers

import (
	"heytobi.dev/fuse/telegram"
)

// ConversationHandler is a suggested default handler. It acts as an orchestrator of the non-command messages received
// by the bot, by keeping track of the conversation context and delegating actions to the appropriate handler.
//
// If it doesn't work well with your use case, you can implement a custom one and register it as your default handler.
type ConversationHandler struct {
}

// Handle ...
func (h *ConversationHandler) Handle(update *telegram.Update) {

}
