package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"unicode"

	"mathbattle/config"
	"mathbattle/infrastructure"
	"mathbattle/libs/fstraverser"
	"mathbattle/models/mathbattle"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expect command")
		return
	}

	switch os.Args[1] {
	case "add-problems":
		container := infrastructure.NewContainer(config.LoadConfig("config.yaml"))
		addProblemsToRepository(container.ProblemRepository(), os.Args[2])
	default:
		fmt.Println("Unknow command")
	}
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
