<h2>
    <img src="./assets/img/fuse.png">
<div>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/heytobi.dev/fuse?status.svg)](https://godoc.org//heytobi.dev/fuse)
[![GitHub Actions](https://github.com/beverlyRoadGoose/fuse/actions/workflows/ci.yml/badge.svg)](https://github.com/beverlyRoadGoose/fuse/actions/workflows/ci.yml)
[![codecov.io](https://codecov.io/gh/beverlyRoadGoose/fuse/coverage.svg?branch=main)](https://codecov.io/gh/beverlyRoadGoose/fuse)
</div>
</h2>

Fuse is a Go library for developing [Telegram](https://telegram.org/) bots, using the [Telegram Bot API](https://core.telegram.org/bots/api).

⚠️ I'm developing this for use in a hobby project, so in the initial phase I'm only adding features as needed in the main project. Overtime I'll aim to cover much of what the Telegram API provides.

## Installation
```console
you@pc:~$ go get -u heytobi.dev/fuse
```

## Current Features
✔️ Register Webhooks  
✔️ Receive updates through Webhooks  
✔️ Send Messages  

## Usage
### Getting Updates through long polling
#### Steps
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

bot.RegisterHandler("/start", func(u interface{}) {
    update, ok := u.(telegram.Update)
    if !ok {
        log.Error("start handler received unexpected type, should be a telegram update")
    }
    
    // Do stuff with the update, for example send a reply:
    
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

### Getting Updates through a Webhook
#### Steps
1. Initialize a Bot
2. Register a Webhook
3. Register command handlers
4. Call the process update method directly whenever your webhook is invoked

```go
httpClient := &http.Client{}
config := &telegram.Config{
    Token: "<YOUR TOKEN>",
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

bot.RegisterHandler("/start", func(u interface{}) {
    update, ok := u.(telegram.Update)
    if !ok {
        log.Error("start handler received unexpected type, should be a telegram update")
    }
    
    // Do stuff with the update, for example send a reply:
    
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

## License
```
MIT License

Copyright (c) 2022 Oluwatobi Adeyinka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
