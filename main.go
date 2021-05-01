package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gen2brain/dlgs"
)

const (
	symbols            string = "aal,aapl,f,tsla"
	token_sand         string = "Tpk_"
	token_live         string = "pk_"
	api_sand           string = "https://sandbox-sse.iexapis.com/stable/news-stream?token=" + token_sand + "&symbols=" + symbols
	api_live           string = "https://cloud-sse.iexapis.com/stable/news-stream?token=" + token_live + "&symbols=" + symbols
	MaxIdleConnections int    = 20
	RequestTimeout     int    = 5
)

type Data struct {
	Datetime   int64  `json:"datetime"`
	Headline   string `json:"headline"`
	Source     string `json:"source"`
	URL        string `json:"url"`
	Summary    string `json:"summary"`
	Related    string `json:"related"`
	Image      string `json:"image"`
	Lang       string `json:"lang"`
	Haspaywall bool   `json:"hasPaywall"`
}

func main() {
	log.Println(api_sand)
	resp, _ := http.Get(api_live)

	reader := bufio.NewReader(resp.Body)
	for {
		line, _ := reader.ReadBytes('\n')
		stringObject := string(line)
		//log.Println("Size: ", len(stringObject))
		if len(stringObject) > 10 {

			log.Println(string(line))
			strReplace := strings.NewReplacer("data:", "", "[", "", "]", "", "\r", "", "\n", "")
			newString := strReplace.Replace(stringObject)

			// Trim quotes
			if newString[0] == '"' {
				newString = newString[1:]
			}
			if i := len(newString) - 1; newString[i] == '"' {
				newString = newString[:i]
			}

			var data Data
			_ = json.Unmarshal([]byte(newString), &data)
			//log.Println(data.Headline)
			log.Printf("Headline %s, related: %s, Link: %s Sum: %s", data.Headline, data.Related, data.URL, data.Summary)

			// Testing Alerts
			/// could be used for windows notification (not clickable ) to send to news alert site
			err := beeep.Alert(data.Headline, data.Summary+" "+data.Source, "assets/warning.png")
			if err != nil {
				panic(err)
			}

			/// pops dialog window - could be work around for allowing click and open browser to send you to news website of alert message
			yes, err := dlgs.Question(data.Headline, "Go to news articale?", true)
			if err != nil {
				panic(err)
			}
			fmt.Println(yes)

		} else {
			log.Println("empty")

		}

		time.Sleep(5 * time.Second)
	}
}
