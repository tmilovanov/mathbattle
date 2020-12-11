package main

import (
	"bytes"
	"errors"
	mathbattle "mathbattle/models"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramPostman struct {
	bot *tb.Bot
}

func (pm *TelegramPostman) PostText(chatID int64, message string) error {
	_, err := pm.bot.Send(tb.ChatID(chatID), message)
	return err
}

func (pm *TelegramPostman) PostPhoto(chatID int64, caption string, image mathbattle.Image) error {
	_, err := pm.bot.Send(tb.ChatID(chatID), &tb.Photo{
		Caption: caption,
		File:    tb.FromReader(bytes.NewReader(image.Content)),
	})
	return err
}

func (pm *TelegramPostman) PostAlbum(chatID int64, caption string, images []mathbattle.Image) error {
	if len(images) < 1 {
		return errors.New("Not enough items to sned")
	}

	inputMedia := []tb.InputMedia{}
	inputMedia = append(inputMedia, &tb.Photo{
		Caption: caption,
		File:    tb.FromReader(bytes.NewReader(images[0].Content)),
	})
	for i := 1; i < len(images); i++ {
		inputMedia = append(inputMedia, &tb.Photo{
			File: tb.FromReader(bytes.NewReader(images[i].Content)),
		})
	}

	_, err := pm.bot.SendAlbum(tb.ChatID(chatID), inputMedia)
	return err
}
