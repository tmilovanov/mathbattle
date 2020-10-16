package models

import tb "gopkg.in/tucnak/telebot.v2"

type TelegramPostman interface {
	Post(participantID string, m *tb.Message) error
}
