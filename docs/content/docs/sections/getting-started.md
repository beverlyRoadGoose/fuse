---
weight: 1
---

# Getting Started

## Installation

Add the library as a dependency in your mod file:
```console
go get -u heytobi.dev/fuse
```

## Initialise your bot

```go
httpClient := &http.Client{}
config := &telegram.Config{
    Token:        "<YOUR TELEGRAM TOKEN>",
    UpdateMethod: telegram.UpdateMethodGetUpdates,
}

bot, err := telegram.NewBot(telegramConfig, httpClient)
if err != nil {
    return nil, errors.New("failed to initialize telegram instance")
}

bot.Start()
```

## Using a Local Bot API Server
If you are [running a Local Bot API Server](https://core.telegram.org/bots/api#using-a-local-bot-api-server), you can
specify the host and the port (if applicable) using the fields exposed in the config struct:

```go
config := &telegram.Config{
    BotApiServer: "https://localserver.net",
    BotApiServerPort: 1234,
}
```