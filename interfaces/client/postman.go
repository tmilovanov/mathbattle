package client

import (
	"errors"
	"fmt"

	"mathbattle/models/mathbattle"
)

type APIPostman struct {
	BaseUrl string
}

func (a *APIPostman) SendSimpleToUsers(msg mathbattle.SimpleMessage) error {
	return PostJsonRecieveNone(fmt.Sprintf("%s%s", a.BaseUrl, "/postman/send_to_users"), msg)
}

func (a *APIPostman) SendSimpleMessage(chatID int64, message string) error {
	return errors.New("Not implemented")
}

func (a *APIPostman) SendImage(chatID int64, imageCaption string, image []byte) error {
	return errors.New("Not implemented")
}

func (a *APIPostman) SendAlbum(chatID int64, albumCaption string, images [][]byte) error {
	return errors.New("Not implemented")
}
