package main

import (
	"errors"
	mathbattle "mathbattle/models"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type TelegramPostman2 struct {
	bot *tgbotapi.BotAPI
}

func (pm *TelegramPostman2) PostText(chatID int64, message string) error {
	_, err := pm.bot.Send(tgbotapi.NewMessage(chatID, message))
	return err
}

func (pm *TelegramPostman2) PostPhoto(chatID int64, caption string, image mathbattle.Image) error {
	return errors.New("Not implemented")
}

func (pm *TelegramPostman2) PostAlbum(chatID int64, caption string, images []mathbattle.Image) error {
	return errors.New("Not implemented")
}
