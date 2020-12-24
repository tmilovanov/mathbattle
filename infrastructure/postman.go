package infrastructure

import (
	"bytes"
	"errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

type TelegramPostman struct {
	bot *tb.Bot
}

func NewTelegramPostman(APIToken string) (*TelegramPostman, error) {
	bot, err := tb.NewBot(tb.Settings{Token: APIToken})
	if err != nil {
		return nil, err
	}

	return &TelegramPostman{bot: bot}, nil
}

func (pm *TelegramPostman) SendSimpleMessage(chatID int64, message string) error {
	_, err := pm.bot.Send(tb.ChatID(chatID), message)
	return err
}

func (pm *TelegramPostman) SendImage(chatID int64, caption string, image []byte) error {
	_, err := pm.bot.Send(tb.ChatID(chatID), &tb.Photo{
		Caption: caption,
		File:    tb.FromReader(bytes.NewReader(image)),
	})
	return err
}

func (pm *TelegramPostman) SendAlbum(chatID int64, caption string, images [][]byte) error {
	if len(images) < 1 {
		return errors.New("Not enough items to send")
	}

	inputMedia := []tb.InputMedia{}
	inputMedia = append(inputMedia, &tb.Photo{
		Caption: caption,
		File:    tb.FromReader(bytes.NewReader(images[0])),
	})
	for i := 1; i < len(images); i++ {
		inputMedia = append(inputMedia, &tb.Photo{
			File: tb.FromReader(bytes.NewReader(images[i])),
		})
	}

	_, err := pm.bot.SendAlbum(tb.ChatID(chatID), inputMedia)
	return err
}
