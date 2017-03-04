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
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const apiParam = "dc6zaTOxFJmzC"
const apiKey = "api_key"
const apiURL = "http://api.giphy.com/v1/gifs/"

var (
	app = kingpin.New("Terminal-GIPHY", "A terminal client for GIPHY")

	trending       = app.Command("trending", "Look at trending GIFs")
	trendingNumber = trending.Arg("number", "The index of GIF to view (indexed at 0 and < 10)").Default("-1").Int()

	search       = app.Command("search", "Search for GIFs")
	query        = search.Arg("query", "Search string").Required().String()
	searchNumber = search.Arg("number", "The index of GIF to view (indexed at 0 and < 10)").Default("-1").Int()
)

func constructURL(endpoint string, extraKey string, extraParam string) (*url.URL, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, endpoint)
	q := u.Query()
	q.Add(apiKey, apiParam)
	q.Add("limit", "10")
	q.Add(extraKey, extraParam)
	u.RawQuery = q.Encode()
	return u, nil
}

func getGifData(u *url.URL, response *terminalGiphy.Response) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	return nil
}

func displayGif(gifURL string) error {
	gifResp, err := http.Get(gifURL)
	if err != nil {
		return err
	}
	defer gifResp.Body.Close()

	imgcat.Cat(gifResp.Body, os.Stdout)
	return nil
}

func displayGifs(response terminalGiphy.Response, number int) error {
	if number > -1 {
		gifImageData := response.Data[number].Images
		gifURL := gifImageData.FixedHeight["url"].(string)

		err := displayGif(gifURL)
		if err != nil {
			return err
		}
	} else {
		for i := 0; i < len(response.Data); i++ {
			gifImageData := response.Data[i].Images
			gifURL := gifImageData.Preview["url"].(string)
			fmt.Printf("[%d] ", i)
			err := displayGif(gifURL)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	var (
		endpoint   string
		extraKey   string
		extraParam string
		number     int
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case trending.FullCommand():
		endpoint = "trending"
		number = *trendingNumber
		if number > 10 {
			fmt.Println("Index very large, try again with less than 10")
			os.Exit(1)
		}

	case search.FullCommand():
		endpoint = "search"
		extraKey = "q"
		extraParam = *query

		number = *searchNumber
		if number > 10 {
			fmt.Println("Index very large, try again with less than 10")
			os.Exit(1)
		}
	}

	u, err := constructURL(endpoint, extraKey, extraParam)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	response := terminalGiphy.Response{}
	err = getGifData(u, &response)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(response.Data) == 0 {
		fmt.Println("No GIFs found :(")
		os.Exit(1)
	}
	err = displayGifs(response, number)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
