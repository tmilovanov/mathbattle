package client

import (
	"fmt"

	"mathbattle/models/mathbattle"
)

type APISolution struct {
	BaseUrl string
}

func (a *APISolution) Create(solution mathbattle.Solution) (mathbattle.Solution, error) {
	result := solution
	err := PostJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/solutions"), &result, &result)
	return result, err
}

func (a *APISolution) Get(ID string) (mathbattle.Solution, error) {
	result := mathbattle.Solution{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/solutions", ID), &result)
	return result, err
}

func (a *APISolution) Find(findDescriptor mathbattle.FindDescriptor) ([]mathbattle.Solution, error) {
	result := []mathbattle.Solution{}
	err := SendGetJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/solutions/find/descriptor"), findDescriptor, &result)
	return result, err
}

func (a *APISolution) AppendPart(ID string, part mathbattle.Image) error {
	return PostJsonRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/solutions/append_part", ID), part)
}

func (a *APISolution) Update(solution mathbattle.Solution) error {
	return PutJsonRecieveNone(fmt.Sprintf("%s%s", a.BaseUrl, "/solutions"), solution)
}

func (a *APISolution) Delete(ID string) error {
	return DeleteRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/solutions", ID))
}

func (a *APISolution) GetProblemDescriptors(participantID string) ([]mathbattle.ProblemDescriptor, error) {
	result := []mathbattle.ProblemDescriptor{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/solutions/descriptors", participantID), &result)
	return result, err
}
