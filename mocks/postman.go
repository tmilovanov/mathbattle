package mocks

import mathbattle "mathbattle/models"

type Postman struct {
	impl map[int64][]string
}

func NewPostman() *Postman {
	return &Postman{
		impl: make(map[int64][]string),
	}
}

func (pm *Postman) PostText(chatID int64, message string) error {
	pm.impl[chatID] = append(pm.impl[chatID], message)
	return nil
}

func (pm *Postman) PostPhoto(chatID int64, caption string, image mathbattle.Image) error {
	return nil
}

func (pm *Postman) PostAlbum(chatID int64, caption string, images []mathbattle.Image) error {
	return nil
}
