package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	terminalGiphy "github.com/bkazi/Terminal-GIPHY/lib"
	imgcat "github.com/martinlindhe/imgcat/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

const apiParam = "dc6zaTOxFJmzC"
const apiKey = "api_key"

const apiURL = "http://api.giphy.com/v1/gifs/"

var (
	app = kingpin.New("Terminal-GIPHY", "A terminal client for GIPHY")

	trending = app.Command("trending", "Look at trending GIFs")

	search = app.Command("search", "Search for GIFs")
	query  = search.Arg("query", "Search string").Required().String()
)

func constructURL(endpoint string, extraKey string, extraParam string) (*url.URL, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)
	q := u.Query()
	q.Add(apiKey, apiParam)
	q.Add("limit", "1")
	q.Add(extraKey, extraParam)
	u.RawQuery = q.Encode()
	return u, nil
}

func main() {
	var (
		endpoint   string
		extraKey   string
		extraParam string
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case trending.FullCommand():
		endpoint = "trending"

	case search.FullCommand():
		endpoint = "search"
		extraKey = "q"
		extraParam = *query
	}

	u, err := constructURL(endpoint, extraKey, extraParam)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	response := terminalGiphy.Response{}
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(response.Data) == 0 {
		fmt.Println("No GIFs found :(")
		os.Exit(1)
	}
	gifImageData := response.Data[0].Images
	gifURL := gifImageData.FixedHeight["url"].(string)

	gifResp, err := http.Get(gifURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer gifResp.Body.Close()

	imgcat.Cat(gifResp.Body, os.Stdout)
}
