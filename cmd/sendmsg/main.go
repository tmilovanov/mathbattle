package main

import (
	"fmt"
	"log"

	"mathbattle/repository/sqlite"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	telegramToken := "1284920547:AAGRFXY0l-jp_OWhx_W3NlHp2_5gtd7O-P8"
	//telegramToken := "1415332759:AAE2r0BtZiYvrg4Ar2ySmajCMgV3PAkyeF8"

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	participantRep, err := sqlite.NewParticipantRepository("mathbattle.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	participants, err := participantRep.GetAll()
	if err != nil {
		log.Fatal(err)
	}

	msg := "Мы проверили присланные решения, вы можете посмотреть комментарии к ним и авторские решения по ссылке: "
	msg += "\n"
	msg += "https://www.notion.so/1f3d2ea04473460f92d59b2c2db0f919"
	msg += "\n"
	msg += "Спасибо за участие! Мы планируем доработать систему и провести следующий раунд в новогодние каникулы."
	for _, p := range participants {
		_, err := bot.Send(tgbotapi.NewMessage(p.TelegramID, msg))
		if err != nil {
			fmt.Printf("Failed to send message to participant %s (chatID: %d), error: %v\n", p.ID, p.TelegramID, err)
		} else {
			fmt.Printf("Successfuly sent message to participant %s\n", p.ID)
		}
	}
}
