package conversation // import "heytobi.dev/fuse/conversation"

import (
	"heytobi.dev/fuse/telegram"
)

// Orchestrator ...
type Orchestrator interface {
	// RegisterSequence ...
	RegisterSequence(chatID int64, sequence Sequence) error

	// DeregisterActiveSequence ...
	DeregisterActiveSequence(chatID int64) error
}

// Sequence ...
type Sequence interface {
	// Start ...
	Start(chatID int64) error

	// Finish ...
	Finish(chatID int64) error

	// Process ..
	Process(update *telegram.Update) error

	// GetName ...
	GetName() string
}