package conversation // import "heytobi.dev/fuse/conversation"

import (
	"context"
	"heytobi.dev/fuse/telegram"
)

type bot interface {
	SendMessage(message *telegram.SendMessageRequest) (*telegram.ActionResult, error)
}

// Handler is a suggested default handler. It acts as an orchestrator of the non-command messages received
// by the bot by keeping track of the conversation context per chat, and delegating actions to the appropriate Sequence.
// Whenever a message is received, the handler checks if an active Sequence is registered for the user that sent the message,
// if an active sequence exists, the message is relayed to that Sequence to be processed. Sequences are responsible for
// their own state management.
//
// If it doesn't work well with your use case, you can implement & register a custom one as your default handler.
type Handler struct {
	bot             bot
	activeSequences map[int64]Sequence
	defaultSequence Sequence
}

// NewHandler ...
func NewHandler(bot bot) *Handler {
	return &Handler{
		bot:             bot,
		activeSequences: make(map[int64]Sequence),
	}
}

// Handle handles every incoming message that doesn't have a dedicated handler.
func (h *Handler) Handle(ctx context.Context, update *telegram.Update) error {
	if update != nil && update.Message != nil {
		// check if there is an active sequence for this user, delegate to that sequence if there is one.
		if sequence, hasActiveSequence := h.activeSequences[update.Message.Chat.ID]; hasActiveSequence {
			err := sequence.Process(ctx, update)
			if err != nil {
				return err
			}
			return nil
		}

		if h.defaultSequence != nil {
			err := h.defaultSequence.Process(ctx, update)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// RegisterActiveSequence registers the active sequence for the given user. New registrations always override any already
// registered sequence. There can be at most 1 active sequences for a user, tracked by the telegram chat ID.
func (h *Handler) RegisterActiveSequence(chatID int64, sequence Sequence) error {
	h.activeSequences[chatID] = sequence
	return nil
}

// DeregisterActiveSequence deletes the active sequence for a user. Sequences can call this method once their flow has
// been completed.
func (h *Handler) DeregisterActiveSequence(chatID int64) error {
	delete(h.activeSequences, chatID)
	return nil
}

// WithDefaultSequence sets a fallback sequence for messages without a dedicated sequence.
func (h *Handler) WithDefaultSequence(sequence Sequence) *Handler {
	h.defaultSequence = sequence
	return h
}
