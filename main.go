package main

import (
	"flag"
	"fmt"
	"forward-info-bot/config"
	"forward-info-bot/handler"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	var (
		path   string
	)

	flag.StringVar(
		&path,
		"config",
		"",
		"enter path to config file",
	)

	// Parse -config argument
	flag.Parse()

	// Init logger
	logger := logrus.New()

	// Get config
	conf, err := config.NewConfig(path)
	if err != nil {
		logger.WithError(err).Fatal("incorrect path or config itself")
	}

	// Set log level from config
	lvl, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		logger.WithError(err).Fatal("cannot parse log level")
	}

	logger.SetLevel(lvl)

	// Init Bot API
	bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		fmt.Println("Telegram bot cannot be initialized! See, error:")
		panic(err)
	}

	fmt.Printf("Authorized on account @%s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Graceful shutdown
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)

	go func() {
		<-s
		updates.Clear()
		os.Exit(1)
	}()

	h := handler.NewHandler(bot, logger, conf)
	for update := range updates {
		go func(u tgbotapi.Update) {
			if u.Message == nil { // ignore any non-Message Updates
				return
			}

			switch u.Message.Command() {
			case "start":
				if err := h.Start(u); err != nil {
					h.Error(u, err)
				}
			default:
				if err := h.Default(u); err != nil {
					h.Error(u, err)
				}
			}
		}(update)
	}
}
