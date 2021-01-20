package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/libs/fstraverser"
	"mathbattle/models/mathbattle"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expect command")
		return
	}

	switch os.Args[1] {
	case "add-problems":
		container := infrastructure.NewServerContainer(config.LoadConfig("config.yaml"))
		addProblemsToRepository(container.ProblemRepository(), os.Args[2])
	case "run-bot":
		configPath := "config.yaml"
		if len(os.Args) > 3 {
			configPath = os.Args[2]
		}
		runBot(configPath)
	case "send-message":
		sendMessage()
	case "get-info":
		getInfo()
	case "send-kb":
		sendKb()
	default:
		fmt.Println("Unknow command")
	}
}

func runBot(configPath string) {
	err := exec.Command("./mb-bot.exe", configPath).Start()
	if err != nil {
		log.Fatalf("Failed to run mbbot, error: %v", err)
	}

	err = exec.Command("./mb-server.exe", configPath).Start()
	if err != nil {
		log.Fatalf("Failed to run mbserver, error: %v", err)
	}
}

func getInfo() {
	container := infrastructure.NewServerContainer(config.LoadConfig("config.yaml"))

	b, _ := tb.NewBot(tb.Settings{
		Token:       container.Config().TelegramToken,
		Poller:      &tb.LongPoller{Timeout: 10 * time.Second},
		Synchronous: true,
		//Verbose:     true,
	})

	ids := []string{"1020341126", "1340754306", "743651179", "760356269", "146590255", "483612890", "875279773",
		"599432958", "719375391", "683841172", "810211817", "796335579", "738821177", "813617082", "930111176", "624013576"}

	for _, chatID := range ids {
		chat, err := b.ChatByID(chatID)
		if err != nil {
			log.Printf("Failed to get chat by id %s", chatID)
		}

		fmt.Printf("INSERT INTO users (tg_chat_id, tg_name, is_admin, registration_time) VALUES (%v, '%v', %v, '0001-01-01 00:00:00');\n", chatID, chat.Username, false)
	}

	names := []string{"Рената", "ка", "Ксения", "Андрей", "Саша", "Сергей", "Сабина", "Александр", "Андрей", "Егор", "Константин",
		"Слава", "Ярослав", "Екатерина", "Сух", "Anastasia"}
	grades := []int{10, 10, 10, 11, 11, 11, 10, 11, 11, 11, 11, 11, 11, 10, 11, 9}
	fmt.Println(len(names), len(grades))
	for i := 0; i < len(names); i++ {
		name := names[i]
		grade := grades[i]
		fmt.Printf("INSERT INTO participants (user_id, name, school, grade, is_active) VALUES (%d, '%s', '', %d, true);\n", i+1, name, grade)
	}
}

func sendMessage() {
	ids := []string{"1020341126", "1340754306", "743651179", "760356269", "146590255", "483612890", "875279773",
		"599432958", "719375391", "683841172", "810211817", "796335579", "738821177", "813617082", "930111176", "624013576"}

	//ids := []string{"442504899"}

	container := infrastructure.NewServerContainer(config.LoadConfig("config.yaml"))

	b, _ := tb.NewBot(tb.Settings{
		Token:       container.Config().TelegramToken,
		Poller:      &tb.LongPoller{Timeout: 10 * time.Second},
		Synchronous: true,
		//Verbose:     true,
	})

	for _, chatID := range ids {
		cID, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			log.Printf("Failed to parse")
		}

		msg := "Привет! Мы сделали систему оценивания решений и комментариев, поправили найденные баги и запускаем новый раунд.\n"
		msg += "Мы хотим подготовить участников к предстоящим матбоям, познакомить новичков с форматом задач и требованиями к решениям. Принять участие в раунде могут все, но задачи будут не очень сложные:)"

		_, err = b.Send(tb.ChatID(cID), msg)
		if err != nil {
			log.Println(err)
		}
	}

}

func sendKb() {
	container := infrastructure.NewServerContainer(config.LoadConfig("config.yaml"))

	b, _ := tb.NewBot(tb.Settings{
		Token:       container.Config().TelegramToken,
		Poller:      &tb.LongPoller{Timeout: 10 * time.Second},
		Synchronous: true,
		//Verbose:     true,
	})

	srep := container.SolutionRepository()
	solutions, err := srep.FindMany("1", "", "")
	if err != nil {
		log.Panic(err)
	}

	curSolutionI := 0
	controlBlock := []*tb.Message{}
	solutionBlock := []tb.Message{}
	b.Handle("/try", func(m *tb.Message) {
		controlBlock = append(controlBlock, m)

		// Solution content
		album := tb.Album{}
		for _, part := range solutions[curSolutionI].Parts {
			album = append(album, &tb.Photo{
				File: tb.FromReader(bytes.NewReader(part.Content)),
			})
		}
		messages, err := b.SendAlbum(m.Sender, album)
		if err != nil {
			log.Panic(err)
		}
		for _, m := range messages {
			solutionBlock = append(solutionBlock, m)
		}

		// Solution description
		menu := &tb.ReplyMarkup{}
		rows := []tb.Row{}
		rows = append(rows, menu.Row(
			tb.Btn{Text: "Комментировать", Data: "comment"},
			tb.Btn{Text: "Оценить", Data: "mark"},
		))
		rows = append(rows, menu.Row(
			tb.Btn{Text: "Комментарии других участников (0)", Data: "data2"},
		))
		if curSolutionI != len(solutions)-1 {
			rows = append(rows, menu.Row(
				tb.Btn{Text: "Следующее решение", Data: "data3"},
				tb.Btn{Text: "Отмена", Data: "data4"},
			))
		} else {
			rows = append(rows, menu.Row(
				tb.Btn{Text: "Отмена", Data: "data4"},
			))
		}
		menu.Inline(rows...)
		solutionDescription := fmt.Sprintf("Решение %d/%d", curSolutionI+1, len(solutions))
		if solutions[curSolutionI].Mark == -1 {
			solutionDescription += " НЕ ОЦЕНЕНО"
		}
		if solutions[curSolutionI].JuriComment == "" {
			solutionDescription += " НЕ ПРОКОММЕНТИРОВАННО"
		}
		msg, err := b.Send(m.Sender, solutionDescription, menu)
		if err != nil {
			log.Panic(err)
		}

		controlBlock = append(controlBlock, msg)
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		log.Printf("Got query: %v", q)
	})

	b.Handle(tb.OnCallback, func(cb *tb.Callback) {
		log.Printf("Got callback: '%s'", cb.Data)

		switch cb.Data {
		case "comment":
			b.Send(cb.Sender, "Введите комментарий:")

		case "data1":
			// Show solution
			album := tb.Album{}
			for _, part := range solutions[curSolutionI].Parts {
				album = append(album, &tb.Photo{
					File: tb.FromReader(bytes.NewReader(part.Content)),
				})
			}
			msg, err := b.SendAlbum(cb.Sender, album)
			if err != nil {
				log.Panic(err)
			}
			for _, m := range msg {
				solutionBlock = append(solutionBlock, m)
			}

			txt := "Нажмите /comment для комментирования решения\n"
			txt += "Нажмите /mark для выставления оценки\n"
			b.Send(cb.Sender, txt)

			log.Println(len(msg))
		case "data2":
			// Show comments
		case "data3":
			// Next solution
		case "data4":
			// Cancel
			for _, msg := range controlBlock {
				b.Delete(msg)
			}
			controlBlock = []*tb.Message{}
			for _, msg := range solutionBlock {
				b.Delete(&msg)
			}
			solutionBlock = []tb.Message{}
		}

		b.Respond(cb, &tb.CallbackResponse{})
	})

	b.Start()

	//album := tb.Album{}
	//for _, part := range solutions[0].Parts {
	//album = append(album, &tb.Photo{
	//File: tb.FromReader(bytes.NewReader(part.Content)),
	//})
	//}
	//b.Send(tb.ChatID(chatID), "Hello!", menu)
}

func addProblemsToRepository(repository mathbattle.ProblemRepository, problemsPath string) {
	fstraverser.TraverseStartingFrom(problemsPath, func(fileInfo fstraverser.FileInformation) {
		f, err := os.Open(fileInfo.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadFile(fileInfo.Path)
		if err != nil {
			log.Fatal(err)
		}

		fileName := filepath.Base(fileInfo.Path)
		if err != nil {
			log.Fatal(err)
		}
		gradeMin, gradeMax, err := extartGradesFromName([]rune(fileName))
		if err != nil {
			log.Fatal(err)
		}
		sha256sum := hex.EncodeToString(h.Sum(nil))

		fmt.Printf("%s [%d;%d] %s\n", fileInfo.Path, gradeMin, gradeMax, sha256sum)
		problem := mathbattle.Problem{
			ID:        sha256sum,
			MinGrade:  gradeMin,
			MaxGrade:  gradeMax,
			Extension: filepath.Ext(fileInfo.Path),
			Content:   content,
		}

		_, err = repository.Store(problem)
		if err != nil {
			log.Fatal(err)
		}
	})

}

func extartGradesFromName(fileName []rune) (int, int, error) {
	if len(fileName) < 3 {
		return 0, 0, errors.New("Filename is too short")
	}

	var err error

	i := 0
	minGrade := 0
	for ; i < len(fileName); i++ {
		if !unicode.IsDigit(fileName[i]) {
			if fileName[i] == '_' {
				minGrade, err = strconv.Atoi(string(fileName[:i]))
				if err != nil {
					return 0, 0, fmt.Errorf("Faild to convert grade part: %v", string(fileName[:i]))
				}
				if !mathbattle.IsValidGrade(minGrade) {
					return 0, 0, fmt.Errorf("Invalid minimal grade: %v", minGrade)
				}
				break
			} else {
				return 0, 0, fmt.Errorf("Unexpected filename format")
			}
		}
	}

	i++
	maxGrade := 0
	maxGradeBegin := i
	for ; i < len(fileName); i++ {
		if !unicode.IsDigit(fileName[i]) {
			maxGrade, err = strconv.Atoi(string(fileName[maxGradeBegin:i]))
			if err != nil {
				return 0, 0, fmt.Errorf("Faild to convert grade part: %v", string(fileName[maxGradeBegin:i]))
			}
			if !mathbattle.IsValidGrade(maxGrade) {
				return 0, 0, fmt.Errorf("Invalid maximal grade: %v", minGrade)
			}
			break
		}
	}

	if minGrade > maxGrade {
		return 0, 0, fmt.Errorf("Invalid grades: minimum grade > maximum grade")
	}

	return minGrade, maxGrade, nil
}
