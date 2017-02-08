package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
)

var contributors = []string{
	"DennyLoko", "hernandev", "joubertredrat", "lucasmezencio", "marcossffilho",
	"marcusagm", "naroga", "pedrommone",
}
var contributorsVotes = map[string]bool{}
var colorMap = map[VoteValue]map[bool]color.Attribute{
	"+1": map[bool]color.Attribute{
		true:  color.FgHiGreen,
		false: color.FgGreen,
	},
	"-1": map[bool]color.Attribute{
		true:  color.FgHiRed,
		false: color.FgRed,
	},
}

func main() {
	if len(os.Args) != 2 {
		color.Red("Invalid arguments")
		os.Exit(1)
	}

	fmt.Println("Getting poll details...")
	fmt.Println()

	url := fmt.Sprintf("https://api.github.com/repos/phpol/phpol/issues/%s", os.Args[1])
	voteJSON, err := getJSON(url)
	if err != nil {
		panic(err)
	}

	p := NewPollFromJSON(voteJSON)

	votesJSON, err := getJSON(p.VotesSummary["url"].(string))
	if err != nil {
		panic(err)
	}
	p.LoadVotes(votesJSON)

	fmt.Println(fmt.Sprintf("Poll Title: %s", p.Title))
	fmt.Println(fmt.Sprintf(
		"Author: %s | Date: %s",
		p.Author.Login,
		p.StartTime.UTC().Format("2006-01-02 15:04:05 MST")))

	fmt.Print("State: ")
	color.Set(color.FgHiGreen)
	if p.State == PollStateClosed {
		color.Set(color.FgHiRed)
	}

	fmt.Println(p.State)
	color.Unset()

	fmt.Println("------------------------------------------------------------")
	fmt.Println(p.Description)
	fmt.Println("------------------------------------------------------------")

	for _, vote := range p.Votes {
		if fgColor, ok := colorMap[vote.Content][isContributor(vote.User.Login)]; ok {
			color.Set(fgColor)
		}

		fmt.Print(fmt.Sprintf("%s ", vote.User.Login))
		contributorsVotes[vote.User.Login] = true
		color.Unset()
	}

	color.Set(color.FgYellow)
	for _, contributor := range contributors {
		if _, ok := contributorsVotes[contributor]; !ok {
			fmt.Print(fmt.Sprintf("%s ", contributor))
		}
	}
	color.Unset()

	fmt.Println()
}

func getJSON(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = map[string][]string{
		"Accept": {"application/vnd.github.squirrel-girl-preview"},
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub replied with %d status code", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}

func isContributor(name string) bool {
	for _, contributor := range contributors {
		if name == contributor {
			return true
		}
	}

	return false
}
