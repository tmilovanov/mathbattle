package handlers

import (
	"errors"
	mreplier "mathbattle/application"
	"mathbattle/infrastructure"
	"mathbattle/models/mathbattle"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

type SendServiceMessage struct {
	Handler
	Replier        mreplier.Replier
	PostmanService mathbattle.PostmanService
}

func (h *SendServiceMessage) Name() string {
	return h.Handler.Name
}

func (h *SendServiceMessage) Description() string {
	return h.Handler.Description
}

func (h *SendServiceMessage) IsShowInHelp(ctx infrastructure.TelegramUserContext) bool {
	res, _, _ := h.IsCommandSuitable(ctx)
	return res
}

func (h *SendServiceMessage) IsCommandSuitable(ctx infrastructure.TelegramUserContext) (bool, string, error) {
	return true, "", nil
}

func (h *SendServiceMessage) IsAdminOnly() bool {
	return true
}

func (h *SendServiceMessage) Handle(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	switch ctx.CurrentStep {
	case 0:
		return 1, OneWithKb(h.Replier.ServiceMsgGetText(), h.Replier.ServiceMsgCancelSend()), nil
	case 1:
		return h.stepAcceptText(ctx, m)
	case 2:
		return h.stepAcceptRecievers(ctx, m)
	case 3: // If recievers not all, accept them here
		return h.stepAcceptParticularRecievers(ctx, m)
	case 4:
		return h.send(ctx, m)
	default:
		return -1, noResponse(), nil
	}
}

func (h *SendServiceMessage) stepAcceptText(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text == h.Replier.ServiceMsgCancelSend() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	if m.Text == "" {
		return 1, OneWithKb(h.Replier.ServiceMsgTextIsEmpty(), h.Replier.ServiceMsgCancelSend()), nil
	}

	ctx.Variables["msg_text"] = infrastructure.NewContextVariableStr(m.Text)

	return 2, OneWithKb(h.Replier.ServiceMsgAskRecieversType(), h.Replier.ServiceMsgCancelSend(),
		h.Replier.ServiceMsgRecieversTypeAll(), h.Replier.ServiceMsgRecieversTypeSome()), nil
}

func (h *SendServiceMessage) stepAcceptRecievers(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text == h.Replier.ServiceMsgCancelSend() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	if m.Text == h.Replier.ServiceMsgRecieversTypeAll() {
		ctx.Variables["recievers"] = infrastructure.NewContextVariableStr("all")
		return 4, OneWithKb(h.Replier.ServiceMsgFinalAsk(h.Replier.ServiceMsgRecieversTypeAll()),
			h.Replier.ServiceMsgCancelSend(), h.Replier.Yes()), nil
	}

	if m.Text == h.Replier.ServiceMsgRecieversTypeSome() {
		return 3, OneWithKb(h.Replier.ServiceMsgInputRecievers(), h.Replier.ServiceMsgCancelSend()), nil
	}

	return 3, OneWithKb(h.Replier.ServiceMsgWrongRecieversType(), h.Replier.ServiceMsgCancelSend()), nil
}

func (h *SendServiceMessage) stepAcceptParticularRecievers(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text == h.Replier.ServiceMsgCancelSend() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	recievers := []string{}
	for _, r := range strings.Split(m.Text, ",") {
		recievers = append(recievers, strings.Trim(r, " \n\t"))
	}

	ctx.Variables["recievers"] = infrastructure.NewContextVariableStr(strings.Join(recievers, ","))

	return 4, OneWithKb(h.Replier.ServiceMsgFinalAsk(h.Replier.ServiceMsgRecieversTypeSome(), recievers...),
		h.Replier.ServiceMsgCancelSend(), h.Replier.Yes()), nil
}

func (h *SendServiceMessage) send(ctx infrastructure.TelegramUserContext, m *tb.Message) (int, []TelegramResponse, error) {
	if m.Text == h.Replier.ServiceMsgCancelSend() {
		return -1, OneTextResp(h.Replier.Cancel()), nil
	}

	msgText, exist := ctx.Variables["msg_text"]
	if !exist {
		return -1, noResponse(), errors.New("Context variable doesn't exist")
	}

	recievers, exist := ctx.Variables["recievers"]
	if !exist {
		return -1, noResponse(), errors.New("Context variable doesn't exist")
	}

	var err error
	if recievers.AsString() == "all" {
		err = h.PostmanService.SendSimpleToUsers(mathbattle.SimpleMessage{
			Text: msgText.AsString(),
		})
	} else {
		err = h.PostmanService.SendSimpleToUsers(mathbattle.SimpleMessage{
			Text:     msgText.AsString(),
			UsersIDS: strings.Split(recievers.AsString(), ","),
		})
	}

	if err != nil {
		return -1, noResponse(), errors.New("Failed to send")
	}

	return -1, OneTextResp(h.Replier.ServiceMsgSendSuccess()), nil
}
