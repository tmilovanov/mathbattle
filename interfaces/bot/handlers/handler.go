package handlers

type Handler struct {
	Name        string
	Description string
}

func noResponse() []TelegramResponse {
	return []TelegramResponse{}
}
