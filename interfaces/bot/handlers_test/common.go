package handlerstest

import (
	"bytes"

	"mathbattle/infrastructure"
	"mathbattle/interfaces/bot/handlers"

	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

type reqRespTextSequence struct {
	request  string
	response string
	step     int
}

type reqRespSequence struct {
	request  tb.Message
	response handlers.TelegramResponse
	step     int
}

func text(input string) tb.Message {
	return tb.Message{Text: input}
}

func photo(input string, fakeFilePath string, fakeFileContent []byte) tb.Message {
	return tb.Message{
		Text: input,
		Photo: &tb.Photo{
			File: tb.File{
				FilePath:   fakeFilePath,
				FileReader: bytes.NewReader(fakeFileContent),
			},
		},
	}
}

func sendAndTest(req *require.Assertions, handler handlers.TelegramCommandHandler, ctx infrastructure.TelegramUserContext,
	msg *tb.Message, expect handlers.TelegramResponse, expectedStep int) infrastructure.TelegramUserContext {

	result := ctx
	if ctx.CurrentStep == -1 {
		ctx.CurrentStep = 0
	}
	step, resp, err := handler.Handle(ctx, msg)
	req.Nil(err)
	req.Equal(expect, resp[0])
	req.Equal(expectedStep, step)
	result.CurrentStep = step

	return result
}

func sendTextExpectTextSequence(req *require.Assertions, handler handlers.TelegramCommandHandler, ctx infrastructure.TelegramUserContext,
	seq []reqRespTextSequence) {

	for _, elem := range seq {
		msg := tb.Message{Text: elem.request}
		expectResp := handlers.NewResp(elem.response)
		expectStep := elem.step
		ctx = sendAndTest(req, handler, ctx, &msg, expectResp, expectStep)
	}
}

func sendReqExpectRespSequence(req *require.Assertions, handler handlers.TelegramCommandHandler, ctx infrastructure.TelegramUserContext,
	seq []reqRespSequence) infrastructure.TelegramUserContext {

	result := ctx
	for _, elem := range seq {
		result = sendAndTest(req, handler, result, &elem.request, elem.response, elem.step)
	}
	return result
}
