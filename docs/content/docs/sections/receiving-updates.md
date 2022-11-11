# Receiving Updates

Telegram provides 2 ways by which your bot can receive messages:

- [Long Polling](https://core.telegram.org/bots/api#getupdates)
- [Webhooks](https://core.telegram.org/bots/api#setwebhook)

## Getting Updates through long polling
**Steps**
1. Initialize a Bot
2. Register command handlers
3. Start polling for updates

```go
httpClient := &http.Client{}
config := &telegram.Config{
    Token:               "<YOUR TOKEN>",
    UpdateMethod:        telegram.UpdateMethodGetUpdates,
    PollingCronSchedule: "*/1 * * * *",
    PollingTimeout:      30,
    PollingUpdatesLimit: 100,
}

poller, err := telegram.NewPoller(telegramConfig, httpClient)
if err != nil {
    return nil, errors.New("failed to initialize telegram poller")
}

bot, err := telegram.NewBot(telegramConfig, httpClient)
if err != nil {
    return nil, errors.New("failed to initialize telegram instance")
}

bot = bot.WithPoller(poller)

bot.RegisterHandler("/start", func(update *telegram.Update) {
    result, err := bot.Send(telegram.SendMessageRequest{
        ChatID: update.Message.Chat.ID,
        Text:   " ¯\_(ツ)_/¯",
    })

    if err != nil {
        log.Error("failed to send telegram message")
    }

    if !result {
        log.Warn("send message result was false")
    }
})

bot.Start() // start listening for updates.
```

{{< hint info >}}
💡This method is very handy during development as you don't have to set up a webhook that Telegram servers can reach 
to test your bot.
{{< /hint >}}

## Getting Updates through a Webhook 
**Steps**
1. Initialize a Bot
2. Register a Webhook
3. Register command handlers
4. Call the process update method directly whenever your webhook is invoked

```go
httpClient := &http.Client{}
config := &telegram.Config{
    Token: "<YOUR TELEGRAM TOKEN>",
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

bot.RegisterHandler("/start", func(update *telegram.Update) {
    result, err := bot.Send(telegram.SendMessageRequest{
        ChatID: update.Message.Chat.ID,
        Text:   " ¯\_(ツ)_/¯",
    })

    if err != nil {
        log.Error("failed to send telegram message")
    }

    if !result {
        log.Warn("send message result was false")
    }
})

// In your webhook http handler:
bot.ProcessUpdate(Update{}) // the update parameter should be deserialized from the request body.
```