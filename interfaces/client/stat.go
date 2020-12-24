package client

import (
	"mathbattle/models/mathbattle"
)

type APIStat struct {
	BaseUrl string
}

func (a *APIStat) Stat() (mathbattle.Stat, error) {
	return mathbattle.Stat{}, nil
}
