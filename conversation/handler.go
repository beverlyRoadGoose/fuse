package conversation // import "heytobi.dev/fuse/conversation"

import (
	"github.com/sirupsen/logrus"
	"heytobi.dev/fuse/telegram"
)

type bot interface {
	SendMessage(message *telegram.SendMessageRequest) (*telegram.ActionResult, error)
}

// Handler is a suggested default handler. It acts as an orchestrator of the non-command messages received
// by the bot, by keeping track of the conversation context per chat, and delegating actions to the appropriate Sequence.
//
// If it doesn't work well with your use case, you can implement & register a custom one as your default handler.
type Handler struct {
	bot             bot
	activeSequences map[int64]Sequence
	defaultResponse string
}

// NewHandler ...
func NewHandler(bot bot) *Handler {
	return &Handler{
		bot:             bot,
		activeSequences: make(map[int64]Sequence),
	}
}

// Handle ...
func (h *Handler) Handle(update *telegram.Update) error {
	if update != nil && update.Message != nil {
		// check if there is an active sequence for this user, delegate to that sequence if there is one.
		if sequence, hasActiveSequence := h.activeSequences[update.Message.Chat.ID]; hasActiveSequence {
			err := sequence.Process(update)
			if err != nil {
				return err
			}
			return nil
		}

		if h.defaultResponse != "" {
			result, err := h.bot.SendMessage(&telegram.SendMessageRequest{
				ChatID: update.Message.Chat.ID,
				Text:   h.defaultResponse,
			})

			if err != nil {
				logrus.WithError(err).Error("failed to send telegram message")
				return err
			}

			if !result.Successful {
				logrus.WithFields(logrus.Fields{
					"response": result,
				}).Warn("send message result was false")
			}
		}
	}

	return nil
}

// RegisterSequence registers the active sequence for the given user. New registrations always override any already
// registered sequence, and there can only be 1 or no active sequences for each user, tracked by the telegram chat ID.
func (h *Handler) RegisterSequence(chatID int64, sequence Sequence) error {
	h.activeSequences[chatID] = sequence
	return nil
}

// DeregisterActiveSequence ...
func (h *Handler) DeregisterActiveSequence(chatID int64) error {
	delete(h.activeSequences, chatID)
	return nil
}

func (h *Handler) WithDefaultResponse(response string) *Handler {
	h.defaultResponse = response
	return h
}
