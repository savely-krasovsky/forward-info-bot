package handler

import (
	"fmt"
	"forward-info-bot/config"
	"forward-info-bot/tool"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Handler struct {
	Telegram *tgbotapi.BotAPI
	Logger   logrus.FieldLogger
	Config   *config.Config
}

func NewHandler(bot *tgbotapi.BotAPI, logger logrus.FieldLogger, conf *config.Config) *Handler {
	return &Handler{
		Telegram: bot,
		Logger:   logger,
		Config:   conf,
	}
}

func (h *Handler) Start(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		`Just forward me some message and I will send you all available information.`,
	)

	if _, err := h.Telegram.Send(msg); err != nil {
		return errors.Wrap(err, "cannot send message")
	}

	return nil
}

func (h *Handler) Default(update tgbotapi.Update) error {
	body := ""

	if update.Message.Text != "" {
		body += fmt.Sprintf("<b>Message:</b> %s\n", update.Message.Text)
	} else {
		body += "<b>Media Type:</b> "
		switch {
		case update.Message.Photo != nil:
			body += "photo"
		case update.Message.Video != nil:
			body += "video"
		case update.Message.VideoNote != nil:
			body += "video note"
		case update.Message.Audio != nil:
			body += "audio"
		case update.Message.Voice != nil:
			body += "voice"
		case update.Message.Sticker != nil:
			body += "sticker"
		case update.Message.Animation != nil:
			body += "animation"
		case update.Message.Document != nil:
			body += "document"
		case update.Message.Game != nil:
			body += "game"
		case update.Message.Contact != nil:
			body += "contact"
		case update.Message.Location != nil:
			body += "location"
		case update.Message.Venue != nil:
			body += "venue"
		default:
			body += "unknown"
		}
		body += "\n"

		if update.Message.Caption != "" {
			body += fmt.Sprintf("<b>Caption:</b> %s\n", update.Message.Caption)
		}
	}

	if update.Message.ForwardDate != 0 {
		switch {
		case update.Message.ForwardFrom != nil:
			var format string
			if update.Message.ForwardFrom.IsBot {
				format = "<b>Bot:</b> %s%s "
			} else if update.Message.ForwardFrom.LastName == "" {
				format = "<b>User:</b> %s%s "
			} else {
				format = "<b>User:</b> %s %s "
			}

			body += fmt.Sprintf(
				format,
				update.Message.ForwardFrom.FirstName,
				update.Message.ForwardFrom.LastName,
			)

			if update.Message.ForwardFrom.UserName != "" {
				body += fmt.Sprintf(
					"(<code>@%s</code> / <code>%d</code>)\n",
					update.Message.ForwardFrom.UserName,
					update.Message.ForwardFrom.ID,
				)
			} else {
				body += fmt.Sprintf(
					"(<code>%d</code>)\n",
					update.Message.ForwardFrom.ID,
				)
			}
		case update.Message.ForwardFromChat != nil:
			if update.Message.ForwardFromChat.Title != "" {
				body += fmt.Sprintf(
					"<b>%s:</b> %s ",
					strings.Title(update.Message.ForwardFromChat.Type),
					update.Message.ForwardFromChat.Title,
				)
			}

			if update.Message.ForwardFromChat.UserName != "" {
				body += fmt.Sprintf(
					"(<code>@%s</code> / <code>%d</code>)\n",
					update.Message.ForwardFromChat.UserName,
					update.Message.ForwardFromChat.ID,
				)
			} else {
				body += fmt.Sprintf(
					"(<code>%d</code>)\n",
					update.Message.ForwardFromChat.ID,
				)
			}
		}

		if update.Message.ForwardFromMessageID != 0 {
			body += fmt.Sprintf(
				"<b>ID:</b> %d\n",
				update.Message.ForwardFromMessageID,
			)
		}

		body += fmt.Sprintf(
			"<b>Date:</b> <code>%s</code>",
			time.Unix(int64(update.Message.ForwardDate), 0).UTC().String(),
		)
	} else {
		if update.Message.From != nil {
			body += fmt.Sprintf(
				"<b>User:</b> %s %s ",
				update.Message.From.FirstName,
				update.Message.From.LastName,
			)
		}

		if update.Message.From.UserName != "" {
			body += fmt.Sprintf(
				"(<code>@%s</code> / <code>%d</code>)\n",
				update.Message.From.UserName,
				update.Message.From.ID,
			)
		} else {
			body += fmt.Sprintf(
				"(<code>%d</code>)\n",
				update.Message.From.ID,
			)
		}

		body += fmt.Sprintf(
			"<b>Expected Lang:</b> <code>%s</code>\n",
			update.Message.From.LanguageCode,
		)

		body += fmt.Sprintf(
			"<b>Date:</b> <code>%s</code>",
			time.Unix(int64(update.Message.Date), 0).UTC().String(),
		)
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		body,
	)
	msg.ParseMode = "HTML"

	if _, err := h.Telegram.Send(msg); err != nil {
		return tool.NewHRError(
			"Cannot send message, maybe it's too long?",
			errors.Wrap(err, "cannot send message"),
		)
	}

	return nil
}

func (h *Handler) Error(update tgbotapi.Update, err error) {
	if err == nil {
		// Why did you call this function?
		return
	}

	// Log error
	h.Logger.WithError(err).Error("error occurred in handler")

	// Send human readable representation of error to user to let him know
	if hrerr, ok := err.(*tool.HRError); ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, hrerr.Human())
		_, err := h.Telegram.Send(msg)
		if err != nil {
			h.Logger.Error(errors.Wrap(err, "cannot send message with human readable error"))
		}
	} else {
		// ... do nothing? Unreadable error useless for people
	}
}
