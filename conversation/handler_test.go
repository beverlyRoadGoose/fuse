package conversation // import "heytobi.dev/fuse/conversation"

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"heytobi.dev/fuse/telegram"
	"testing"
)

func TestCanRegisterSequence(t *testing.T) {
	chatID := int64(1)

	bot := &mockBot{}
	handler := NewHandler(bot)

	handler.RegisterSequence(chatID, &mockSequence{})

	_, registered := handler.activeSequences[chatID]
	assert.True(t, registered)
}

func TestCanDeregisterSequence(t *testing.T) {
	chatID := int64(1)

	bot := &mockBot{}
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

	bot := &mockBot{}
	handler := NewHandler(bot)

	handler.RegisterSequence(chatID, firstSequence)
	assert.Equal(t, firstSequenceName, handler.activeSequences[chatID].GetName())

	handler.RegisterSequence(chatID, secondSequence)
	assert.Equal(t, secondSequenceName, handler.activeSequences[chatID].GetName())
}

func TestHandleReturnsErrorIfProcessFails(t *testing.T) {
	chatID := int64(1)
	update := &telegram.Update{
		Message: &telegram.Message{
			Chat: &telegram.Chat{ID: chatID},
		},
	}

	expectedError := errors.New("delegation failed")

	sequence := &mockSequence{}
	sequence.On("Process", update).Return(expectedError)

	bot := &mockBot{}
	handler := NewHandler(bot)
	handler.RegisterSequence(chatID, sequence)

	err := handler.Handle(update)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestHandleReturnsErrorIfSendingDefaultMessageFails(t *testing.T) {
	chatID := int64(1)
	update := &telegram.Update{
		Message: &telegram.Message{
			Chat: &telegram.Chat{ID: chatID},
		},
	}

	expectedError := errors.New("send message failed")

	bot := &mockBot{}
	bot.On("SendMessage", mock.Anything).Return(nil, expectedError)

	handler := NewHandler(bot).WithDefaultResponse("default response")

	err := handler.Handle(update)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
