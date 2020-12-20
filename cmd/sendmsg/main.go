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

	msg := "Здравствуйте! Сегодня в 23:59 по московскому времени закончится приём решений по задачам. "
	msg += "Убедитесь, что вы отправили все решения, которые хотели. "
	msg += "Если у вас возникли какие-то проблемы с отправкой решений напишите: @mathbattle_support"
	for _, p := range participants {
		_, err := bot.Send(tgbotapi.NewMessage(p.TelegramID, msg))
		if err != nil {
			fmt.Printf("Failed to send message to participant %s (chatID: %d), error: %v\n", p.ID, p.TelegramID, err)
		} else {
			fmt.Printf("Successfuly sent message to participant %s\n", p.ID)
		}
	}
}
