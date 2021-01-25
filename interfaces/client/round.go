package client

import (
	"fmt"

	"mathbattle/models/mathbattle"
)

type APIRound struct {
	BaseUrl string
}

func (a *APIRound) StartNew(startOrder mathbattle.StartOrder) (mathbattle.StartResult, error) {
	result := mathbattle.StartResult{}
	err := PostJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/start"), startOrder, &result)
	return result, err
}

func (a *APIRound) StartReviewStage(startOrder mathbattle.StartOrder) (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := PostJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/start_review"), startOrder, &result)
	return result, err
}

func (a *APIRound) ReviewStageDistributionDesc() (mathbattle.ReviewDistributionDesc, error) {
	var result mathbattle.ReviewDistributionDesc
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/review_stage_distribution"), &result)
	return result, err
}

func (a *APIRound) GetAll() ([]mathbattle.Round, error) {
	result := []mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds"), &result)
	return result, err
}

func (a *APIRound) GetByID(ID string) (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/rounds", ID), &result)
	return result, err
}

func (a *APIRound) GetRunning() (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/running"), &result)
	return result, err
}

func (a *APIRound) GetReviewPending() (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/review_pending"), &result)
	return result, err
}

func (a *APIRound) GetReviewRunning() (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/review_running"), &result)
	return result, err
}

func (a *APIRound) GetLast() (mathbattle.Round, error) {
	result := mathbattle.Round{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/rounds/last"), &result)
	return result, err
}

func (a *APIRound) GetProblemDescriptors(participantID string) ([]mathbattle.ProblemDescriptor, error) {
	result := []mathbattle.ProblemDescriptor{}
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/rounds/problem_descriptors", participantID), &result)
	return result, err
}
