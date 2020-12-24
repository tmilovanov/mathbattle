package client

import (
	"fmt"
	"mathbattle/models/mathbattle"
)

type APIProblem struct {
	BaseUrl string
}

func (a *APIProblem) GetByID(ID string) (mathbattle.Problem, error) {
	result := mathbattle.Problem{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/problems", ID), &result)
	return result, err
}
