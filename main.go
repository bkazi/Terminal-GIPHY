package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	imgcat "github.com/martinlindhe/imgcat/lib"
)

const apiParam = "dc6zaTOxFJmzC"
const apiKey = "api_key"

const trendingURL = "http://api.giphy.com/v1/gifs/trending"

func main() {

	u, err := url.Parse(trendingURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Add(apiKey, apiParam)
	u.RawQuery = q.Encode()

	var responseJSON map[string]interface{}
	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	error := json.Unmarshal(body, &responseJSON)
	if error != nil {
		log.Fatal(error)
	}

	data := responseJSON["data"].([]interface{})
	gif := data[0].(map[string]interface{})
	gifImageData := gif["images"].(map[string]interface{})
	gifURL := gifImageData["fixed_height"].(map[string]interface{})["url"].(string)

	gifResp, _ := http.Get(gifURL)
	defer gifResp.Body.Close()

	imgcat.Cat(gifResp.Body, os.Stdout)
}
