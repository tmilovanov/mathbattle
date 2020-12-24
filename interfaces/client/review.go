package client

import (
	"fmt"
	"mathbattle/models/mathbattle"
)

type APIReview struct {
	BaseUrl string
}

func (a *APIReview) Store(review mathbattle.Review) (mathbattle.Review, error) {
	result := mathbattle.Review{}
	err := PostJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/reviews"), review, &result)
	return result, err
}

func (a *APIReview) FindMany(descriptor mathbattle.ReviewFindDescriptor) ([]mathbattle.Review, error) {
	result := []mathbattle.Review{}
	err := SendGetJsonRecieveJson(fmt.Sprintf("%s%s", a.BaseUrl, "/reviews/find/descriptor"), descriptor, &result)
	return result, err
}

func (a *APIReview) Delete(ID string) error {
	return DeleteRecieveNone(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/reviews", ID))
}

func (a *APIReview) RevewStageDescriptors(participantID string) ([]mathbattle.SolutionDescriptor, error) {
	var result []mathbattle.SolutionDescriptor
	err := SendGetNoneRecieveJson(fmt.Sprintf("%s%s/%s", a.BaseUrl, "/reviews/descriptors", participantID), &result)
	return result, err
}
