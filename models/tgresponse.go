package models

import tb "gopkg.in/tucnak/telebot.v2"

type TelegramResponse struct {
	Text     string
	Keyboard *tb.ReplyMarkup
}

func NewResp(messageText string) TelegramResponse {
	return TelegramResponse{messageText, &tb.ReplyMarkup{
		ReplyKeyboardRemove: true,
	}}
}

func NewRespWithKeyboard(messageText string, buttonTexts ...string) TelegramResponse {
	keyboard := &tb.ReplyMarkup{
		ResizeReplyKeyboard: true,
	}

	buttons := []tb.Btn{}
	for _, txt := range buttonTexts {
		buttons = append(buttons, keyboard.Text(txt))
	}

	keyboard.Reply(keyboard.Row(buttons...))

	return TelegramResponse{messageText, keyboard}
}
