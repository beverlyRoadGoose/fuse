package conversation // import "heytobi.dev/fuse/conversation"

import (
	"heytobi.dev/fuse/telegram"
)

// Orchestrator ...
type Orchestrator interface {
	// RegisterSequence ...
	RegisterSequence(chatID int64, sequence Sequence)

	// DeregisterActiveSequence ...
	DeregisterActiveSequence(chatID int64)
}

// Sequence ...
type Sequence interface {
	// Process ..
	Process(update *telegram.Update) error

	// GetName ...
	GetName() string
}
