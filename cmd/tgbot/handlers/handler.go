package handlers

import (
	mathbattle "mathbattle/models"
)

type Handler struct {
	Name        string
	Description string
}

func noResponse() mathbattle.TelegramResponse {
	var result mathbattle.TelegramResponse
	return result
}
