package conversation // import "heytobi.dev/fuse/conversation"

import (
	"github.com/sirupsen/logrus"
	"heytobi.dev/fuse/telegram"
)

// Handler is a suggested default handler. It acts as an orchestrator of the non-command messages received
// by the bot, by keeping track of the conversation context per chat, and delegating actions to the appropriate Sequence.
//
// If it doesn't work well with your use case, you can implement & register a custom one as your default handler.
type Handler struct {
	bot             *telegram.Bot
	activeSequences map[int64]Sequence
	defaultResponse string
}

// NewHandler ...
func NewHandler(bot *telegram.Bot) *Handler {
	return &Handler{
		bot:             bot,
		activeSequences: make(map[int64]Sequence),
	}
}

// Handle ...
func (h *Handler) Handle(update *telegram.Update) {
	if update != nil && update.Message != nil {
		// check if there is an active sequence for this user, delegate to that sequence if there is one.
		if sequence, hasActiveSequence := h.activeSequences[update.Message.Chat.ID]; hasActiveSequence {
			sequence.Process(update)
			return
		}

		if h.defaultResponse != "" {
			result, err := h.bot.SendMessage(&telegram.SendMessageRequest{
				ChatID: update.Message.Chat.ID,
				Text:   h.defaultResponse,
			})

			if err != nil {
				logrus.WithError(err).Error("failed to send telegram message")
			}

			if !result.Successful {
				logrus.WithFields(logrus.Fields{
					"response": result,
				}).Warn("send message result was false")
			}
		}
	}
}

// RegisterSequence registers the active sequence for the given user. New registrations always override any already
// registered sequence, and there can only be 1 or no active sequences for each user, tracked by the telegram chat ID.
func (h *Handler) RegisterSequence(chatID int64, sequence Sequence) {
	h.activeSequences[chatID] = sequence
}

// DeregisterActiveSequence ...
func (h *Handler) DeregisterActiveSequence(chatID int64) {
	delete(h.activeSequences, chatID)
}

func (h *Handler) WithDefaultResponse(response string) {
	h.defaultResponse = response
}
