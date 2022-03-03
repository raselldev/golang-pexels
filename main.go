package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	photoApi = "https://api.pexels.com/v1"
	videoApi = "https://api.pexels.com/videos"
)

type Client struct {
	Token          string
	Hc             http.Client
	RemainingTimes int32
}

type SearchResult struct {
	Page        int32   `json:"page"`
	PerPage     int32   `json:"perPage"`
	TotalResult int32   `json:"totalResult"`
	NextPage    string  `json:"nextPage"`
	Photos      []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographerUrl"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Potrait   string `json:"potrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

func (c *Client) searchPhotos(query string, perPage, page int) (*SearchResult, error) {
	url := fmt.Sprintf(photoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	response, err := c.requestDoWithAuth("GET", url)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)
	response, err := c.Hc.Do(req)
	if err != nil {
		return response, err
	}
	times, err := strconv.Atoi(response.Header.Get("X-Ratelimit-Remaining"))
	if err != nil {
		return response, nil
	} else {
		c.RemainingTimes = int32(times)
	}
	return response, nil
}

func newClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, Hc: c}
}

func main() {
	os.Setenv("pexelsToken", "YOUR_PEXELS_TOKEN")
	TOKEN := os.Getenv("pexelsToken")
	c := newClient(TOKEN)

	result, err := c.searchPhotos("high", 15, 1)

	if err != nil {
		fmt.Errorf("Search Error:%v", err)
	}

	if result.Page == 0 {
		fmt.Println("Search Result Wrong")
	}

	fmt.Println(result)
	file, _ := json.MarshalIndent(result, "", "")
	_ = ioutil.WriteFile("data.json", file, 0644)
}
