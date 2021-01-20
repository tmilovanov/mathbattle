package mathbattle

type SimpleMessage struct {
	Text     string   `json:"text"`
	UsersIDS []string `json:"users_ids"`
}

type PostmanService interface {
	SendSimpleToUsers(msg SimpleMessage) error
	SendSimpleMessage(chatID int64, message string) error
	SendImage(chatID int64, imageCaption string, image []byte) error
	SendAlbum(chatID int64, albumCaption string, images [][]byte) error
}
