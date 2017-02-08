package main

import "time"

type VoteValue string

const VoteYes = "+1"
const VoteNo = "-1"

type Vote struct {
	ID      int       `json:"id"`
	User    User      `json:"user"`
	Content VoteValue `json:"content"`
	Time    time.Time `json:"created_at"`
}
