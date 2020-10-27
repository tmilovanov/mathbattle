package models

type TelegramPostman interface {
	PostText(chatID int64, message string) error
	PostPhoto(chatID int64, caption string, image Image) error
	PostAlbum(chatID int64, caption string, images []Image) error
}
