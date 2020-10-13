package handlerstest

import (
	"bytes"
	mathbattle "mathbattle/models"

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
	response mathbattle.TelegramResponse
	step     int
}

func textReq(input string) tb.Message {
	return tb.Message{Text: input}
}

func photoReq(input string, fakeFilePath string, fakeFileContent []byte) tb.Message {
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

func sendAndTest(req *require.Assertions, handler mathbattle.TelegramCommandHandler, ctx mathbattle.TelegramUserContext,
	msg *tb.Message, expect mathbattle.TelegramResponse, expectedStep int) mathbattle.TelegramUserContext {

	result := ctx
	step, resp, err := handler.Handle(ctx, msg)
	req.Nil(err)
	req.Equal(resp, expect)
	req.Equal(step, expectedStep)
	result.CurrentStep = step

	return result
}

func sendTextExpectTextSequence(req *require.Assertions, handler mathbattle.TelegramCommandHandler, ctx mathbattle.TelegramUserContext,
	seq []reqRespTextSequence) {

	for _, elem := range seq {
		msg := tb.Message{Text: elem.request}
		expectResp := mathbattle.NewResp(elem.response)
		expectStep := elem.step
		ctx = sendAndTest(req, handler, ctx, &msg, expectResp, expectStep)
	}
}

func sendReqExpectRespSequence(req *require.Assertions, handler mathbattle.TelegramCommandHandler, ctx mathbattle.TelegramUserContext,
	seq []reqRespSequence) {

	for _, elem := range seq {
		ctx = sendAndTest(req, handler, ctx, &elem.request, elem.response, elem.step)
	}
}

func getTestDbName() string       { return "test_mathbattle.sqlite" }
func getTestSolutionName() string { return "test_solutions" }
func getTestProblemsName() string { return "test_problems" }
