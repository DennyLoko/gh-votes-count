package main

import (
	"encoding/json"
	"time"
)

type PollState string

const PollStateOpen = "open"
const PollStateClosed = "closed"

type Poll struct {
	ID           int                    `json:"id"`
	Number       int                    `json:"number"`
	Title        string                 `json:"title"`
	Description  string                 `json:"body"`
	Author       User                   `json:"user"`
	State        PollState              `json:"state"`
	StartTime    time.Time              `json:"created_at"`
	VotesSummary map[string]interface{} `json:"reactions"`

	Votes       []Vote `json:"-"`
	votesLoaded bool
}

func NewPollFromJSON(data []byte) *Poll {
	p := &Poll{}

	err := json.Unmarshal(data, p)
	if err != nil {
		panic(err)
	}

	return p
}

func (p *Poll) LoadVotes(data []byte) error {
	votes := make([]Vote, 0)
	json.Unmarshal(data, &votes)

	p.Votes = votes
	p.votesLoaded = true

	return nil
}
