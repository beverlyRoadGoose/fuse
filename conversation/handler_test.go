package conversation // import "heytobi.dev/fuse/conversation"

import (
	"github.com/stretchr/testify/assert"
	"heytobi.dev/fuse/telegram"
	"testing"
)

func TestCanRegisterSequence(t *testing.T) {
	chatID := int64(1)

	bot := &telegram.Bot{}
	handler := NewHandler(bot)

	handler.RegisterSequence(chatID, &mockSequence{})

	_, registered := handler.activeSequences[chatID]
	assert.True(t, registered)
}

func TestCanDeregisterSequence(t *testing.T) {
	chatID := int64(1)

	bot := &telegram.Bot{}
	handler := NewHandler(bot)

	handler.RegisterSequence(chatID, &mockSequence{})
	handler.DeregisterActiveSequence(chatID)

	_, registered := handler.activeSequences[chatID]
	assert.False(t, registered)
}

func TestRegisterOverridesExistingSequence(t *testing.T) {
	chatID := int64(1)
	firstSequenceName := "first"
	secondSequenceName := "second"

	firstSequence := &mockSequence{}
	secondSequence := &mockSequence{}

	firstSequence.On("GetName").Return(firstSequenceName)
	secondSequence.On("GetName").Return(secondSequenceName)

	bot := &telegram.Bot{}
	handler := NewHandler(bot)

	handler.RegisterSequence(chatID, firstSequence)
	assert.Equal(t, firstSequenceName, handler.activeSequences[chatID].GetName())

	handler.RegisterSequence(chatID, secondSequence)
	assert.Equal(t, secondSequenceName, handler.activeSequences[chatID].GetName())
}
