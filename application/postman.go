package application

import (
	"log"
	"mathbattle/models/mathbattle"
)

type PostmanService struct {
	Users   mathbattle.UserRepository
	Postman mathbattle.PostmanService
}

func (s *PostmanService) SendSimpleToUsers(msg mathbattle.SimpleMessage) error {
	if len(msg.UsersIDS) == 0 {
		users, err := s.Users.GetAll()
		if err != nil {
			log.Printf("[PostmanService][SendSimpleToUsers] Failed to get all users, error: %v", err)
			return err
		}

		for _, user := range users {
			err = s.SendSimpleMessage(user.TelegramID, msg.Text)
			if err != nil {
				log.Printf("[PostmanService][SendSimpleToUsers] Failed to send to user, error: %v", err)
				return err
			}
		}
	} else {
		for _, userID := range msg.UsersIDS {
			user, err := s.Users.GetByID(userID)
			if err != nil {
				log.Printf("[PostmanService][SendSimpleToUsers] Failed to get user, error: %v", err)
				return err
			}

			err = s.SendSimpleMessage(user.TelegramID, msg.Text)
			if err != nil {
				log.Printf("[PostmanService][SendSimpleToUsers] Failed to send to user, error: %v", err)
				return err
			}
		}
	}

	return nil
}

func (s *PostmanService) SendSimpleMessage(chatID int64, message string) error {
	return s.Postman.SendSimpleMessage(chatID, message)
}

func (s *PostmanService) SendImage(chatID int64, imageCaption string, image []byte) error {
	return s.Postman.SendImage(chatID, imageCaption, image)
}

func (s *PostmanService) SendAlbum(chatID int64, albumCaption string, images [][]byte) error {
	return s.Postman.SendAlbum(chatID, albumCaption, images)
}
