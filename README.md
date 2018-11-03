# Golang package to avoid Telegram limits

[![GoDoc](https://godoc.org/github.com/soluchok/tsender?status.svg)](http://godoc.org/github.com/soluchok/tsender)
[![Build Status](https://github.com/soluchok/tsender/workflows/build/badge.svg)](https://github.com/soluchok/tsender/actions)
[![codecov](https://codecov.io/gh/soluchok/tsender/branch/master/graph/badge.svg)](https://codecov.io/gh/soluchok/tsender)

Features:
- messages per user are delivered consistently
- group and users messages are supported
- messages are sent deferred

## Install

```
go get -u github.com/soluchok/tsender
```

## Example

```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/soluchok/tsender"
)

const (
	token         = "1****:******************A"
	chatID  int64 = 64432422
	groupID int64 = -243783656
	workers int   = 8
)

func getClient() *http.Client {
	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		panic("http.DefaultClient is not an http.Client")
	}
	transport.MaxIdleConnsPerHost = workers
	return &http.Client{Transport: transport}
}

type Provider struct {
	bot *tgbotapi.BotAPI
}

func (p *Provider) Send(msg interface{}) {
	_, err := p.bot.Send(msg.(tgbotapi.Chattable))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPIWithClient(token, getClient())
	if err != nil {
		panic(err)
	}

	sender := tsender.NewSender(&Provider{bot})

	go sender.Run(workers)
	defer sender.Stop()

	for i := 1; i <= 5; i++ {
		sender.Send(chatID, tgbotapi.NewMessage(chatID, fmt.Sprintf("Hi, chat %d", i)))
		sender.Send(groupID, tgbotapi.NewMessage(groupID, fmt.Sprintf("Hi, group %d", i)))
	}

	time.Sleep(time.Second * 5)
}
```