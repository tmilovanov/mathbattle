package mathbattle

type Postman interface {
	SendSimpleMessage(chatID int64, message string) error
	SendImage(chatID int64, imageCaption string, image []byte) error
	SendAlbum(chatID int64, albumCaption string, images [][]byte) error
}
