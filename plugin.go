package main

import (
	"github.com/FundingCircle/drone-marathon/marathon"
)

type Plugin struct {
	Marathon marathon.Marathon
}

func (p Plugin) Exec() error {
	client := marathon.NewClient(&p.Marathon)
	return client.CreateOrUpdateApplication()
}
