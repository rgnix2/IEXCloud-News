package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-toast/toast"
	"github.com/joho/godotenv"
)

const (
	sandUrl            string = "https://sandbox-sse.iexapis.com/stable/"
	liveUrl            string = "https://cloud-sse.iexapis.com/stable/"
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

var (
	appPath   string
	iexSand   string
	iexLive   string
	liveToken string
	sandToken string
	apiUrl    string
)

func init() {
	// get app path
	getPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	appPath = getAsset(getPath)

	// get env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sandToken = os.Getenv("TOKEN_SAND")
	liveToken = os.Getenv("TOKEN_LIVE")

	// set flags
	var isDebug bool
	var getStocks string

	flag.StringVar(&getStocks, "stock", "", "fb,appl,tsla")
	flag.BoolVar(&isDebug, "debug", false, "true")
	flag.Parse()

	// api url structure
	iexSand = sandUrl + "news-stream?token=" + sandToken + "&symbols=" + getStocks
	iexLive = liveUrl + "news-stream?token=" + liveToken + "&symbols=" + getStocks

	if isDebug {
		log.Println("Debug: Sandbox data")
		apiUrl = iexSand
	}
	if !isDebug {
		log.Println("Live: Listening for news..")
		apiUrl = iexLive
	}

}

func main() {

	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatalln(err)
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatalln(err)
		}

		stringObject := string(line)

		if len(stringObject) > 10 {
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
			err = json.Unmarshal([]byte(newString), &data)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println("New News:", data.Source)

			// Alerts
			notification := toast.Notification{
				AppID:               "NixStockNews",
				Title:               data.Headline,
				Message:             data.Summary,
				ActivationArguments: data.URL + "?token=" + liveToken,
				Icon:                appPath,
			}
			err := notification.Push()
			if err != nil {
				log.Fatalln(err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func getAsset(path string) string {

	strReplace := strings.NewReplacer("\\", "/")
	newString := strReplace.Replace(path)
	assetString := newString + "/assets/alert.png"

	return assetString
}
