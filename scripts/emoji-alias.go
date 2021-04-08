package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Emoji struct {
	Aliases  []string `json:"aliases"`
	Filename string   `json:"filename"`
}

func main() {
	resp, err := http.Get(`https://raw.githubusercontent.com/mattermost/mattermost-webapp/master/utils/emoji.json`)
	if err != nil {
		log.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("invalid response code: %s", resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error)
	}
	var e []Emoji
	if err := json.Unmarshal(b, &e); err != nil {
		log.Fatal(err.Error())
	}

	table := []string{
		"| emoji | name | aliases |",
		"|:------|:-----|:--------|",
	}
	for _, v := range e {
		if len(v.Aliases) >= 2 {
			orig := v.Aliases[0]
			aliases := []string{}
			for _, a := range v.Aliases[1:] {
				aliases = append(aliases, fmt.Sprintf("`:%s:`", a))
			}
			table = append(table, fmt.Sprintf("| [png](%s) | `:%s:` | %s |", fmt.Sprintf(`https://github.com/mattermost/mattermost-webapp/blob/master/images/emoji/%s.png`, v.Filename), orig, strings.Join(aliases, ", ")))
		}
	}

	fmt.Println("## Emoji aliases")
	fmt.Println()
	fmt.Println(strings.Join(table, "\n"))
}
