package client

import (
	"fmt"

	"mathbattle/models/mathbattle"
)

type APIParticipant struct {
	BaseUrl string
}

func (a *APIParticipant) Store(participant mathbattle.Participant) (mathbattle.Participant, error) {
	result := participant
	err := PostJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/participants"), &result, &result)
	return result, err
}

func (a *APIParticipant) GetByID(ID string) (mathbattle.Participant, error) {
	result := mathbattle.Participant{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/participants", ID), &result)
	return result, err
}

func (a *APIParticipant) GetByTelegramID(TelegramID int64) (mathbattle.Participant, error) {
	result := mathbattle.Participant{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%d", a.BaseUrl, "/participants/telegram", TelegramID), &result)
	return result, err
}

func (a *APIParticipant) GetAll() ([]mathbattle.Participant, error) {
	result := []mathbattle.Participant{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/participants"), &result)
	return result, err
}

func (a *APIParticipant) Update(participant mathbattle.Participant) error {
	return PutJsonRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/participants", participant.ID), participant)
}

func (a *APIParticipant) Delete(ID string) error {
	return DeleteRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/participants", ID))
}

func (a *APIParticipant) Unsubscribe(ID string) error {
	return PostNoneRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/participants/unsubscribe", ID))
}
