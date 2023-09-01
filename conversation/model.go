package conversation // import "heytobi.dev/fuse/conversation"

import (
	"context"
	"heytobi.dev/fuse/telegram"
)

// Orchestrator defines a type responsible for orchestrating sequences.
type Orchestrator interface {
	// RegisterActiveSequence registers the given sequence as active for the given user.
	RegisterActiveSequence(chatID int64, sequence Sequence) error

	// DeregisterActiveSequence clears the active sequence for the given user.
	DeregisterActiveSequence(chatID int64) error
}

// Sequence can be thought of as the context of a conversation. It is responsible for its own state management
// and making sense of how an individual message fits into the broader conversation.
type Sequence interface {
	// Start initiates the sequence
	Start(ctx context.Context, update *telegram.Update) error

	// Finish wraps up a sequence
	Finish() error

	// Process processes the given update as part of the sequence
	Process(ctx context.Context, update *telegram.Update) error

	// GetName returns the name of the sequence
	GetName() string

	// TODO explore time based self-closing sequences.
}
