package models

import tb "gopkg.in/tucnak/telebot.v2"

type TelegramResponse struct {
	Text     string
	Img      Image
	Keyboard *tb.ReplyMarkup
}

func NewResp(messageText string) TelegramResponse {
	return TelegramResponse{
		Text: messageText,
		Keyboard: &tb.ReplyMarkup{
			ReplyKeyboardRemove: true,
		},
	}
}

func NewRespImage(image Image) TelegramResponse {
	return TelegramResponse{
		Img: image,
		Keyboard: &tb.ReplyMarkup{
			ReplyKeyboardRemove: true,
		},
	}
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

	return TelegramResponse{
		Text:     messageText,
		Keyboard: keyboard,
	}
}

func NewResps(messageTexts ...string) []TelegramResponse {
	result := []TelegramResponse{}

	for _, item := range messageTexts {
		result = append(result, NewResp(item))
	}

	return result
}

func OneTextResp(messageText string) []TelegramResponse {
	return []TelegramResponse{
		NewResp(messageText),
	}
}

func OneWithKb(messageText string, buttonTexts ...string) []TelegramResponse {
	return []TelegramResponse{
		NewRespWithKeyboard(messageText, buttonTexts...),
	}
}
