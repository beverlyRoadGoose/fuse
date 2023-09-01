# Sending Messages

## Replying a Message

In this snippet, we reply a user with the text `¯\_(ツ)_/¯` when they send the `/start` command to our bot.

```go
httpClient := &http.Client{}
config := &telegram.Config{
    Token:        "<YOUR TELEGRAM TOKEN>",
    UpdateMethod: telegram.UpdateMethodWebhook,
}

bot, err := telegram.Init(config, httpClient)
if err != nil {
    log.Fatal("failed to initialize telegram bot")
}

bot.RegisterWebhook(telegram.Webhook{url: "mywebhook.com/notify"})
if err != nil {
    log.Fatal("failed to register webhook")
}

bot.RegisterHandler("/start", func(ctx context.Context, update *telegram.Update) {
    result, err := bot.Send(telegram.SendMessageRequest{
        ChatID: update.Message.Chat.ID,
        Text:   " ¯\_(ツ)_/¯",
    })

    if err != nil {
        log.Error("failed to send telegram message")
    }

    if !result.Successful {
        log.Warn(fmt.Sprintf("failed to send telegram message: %s", result.Description))
    }
})
```

The `messageSent` flag is set based on the `ok` field in the response for the telegram API call. You can read more about
the `ok` field [in the Telegram docs](https://core.telegram.org/bots/api#making-requests).