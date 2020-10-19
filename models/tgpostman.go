package models

import tb "gopkg.in/tucnak/telebot.v2"

type TelegramPostman interface {
	Post(chatID int64, m *tb.Message) error
}
