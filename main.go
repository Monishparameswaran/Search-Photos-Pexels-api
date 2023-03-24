package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoApi = "https://api.pexels.com/v1"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

// json is a data formatter tool that can be used to make your language dependent data to language indepent one
// usually the JSON format is {"key":"value"}

type SearchResult struct {
	Page         int32   `json:"page"` // this say JSON that use key page while parsing and encoding
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos	`
}
type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}
type Photo struct {
	Id     int32       `json:"id"`
	Width  int32       `json:"width"`
	Height int32       `json:"height`
	Src    PhotoSource `json:"src"`
}
type PhotoSource struct {
	Original string `json:"original"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
}

func (c *Client) SearchPhotos(query string, perPage int, page int) (*SearchResult, error) { // it is a struct method
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&perpage=%d&page=%d", query, perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}
func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.requestDoWithAuth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result CuratedResult
	err = json.Unmarshal(data, &result)
	return &result, err
}
func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	/*
	   Client:=http.Client{};   // creates a new virtual client
	   req,err:=http.NewRequest("GET","URL",nil)  // create a NewRequest with the sepcified method and URL
	   if err!=nil{"error handle"}
	   resp,err:=Client.do(req)    response can be get by doing a request
	   defer resp.Body.Close();

	*/
	req.Header.Add("Authorization", c.Token)
	resp, err := c.hc.Do(req)
	if err != nil {
		return resp, nil
	}
	times, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-remaining")) //X-Ratelimit-Remaining is the value returned by the API that indicates number of API calll remaining
	if err != nil {
		return resp, nil
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, nil
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id) //it builds the URL to be searched
	resp, err := c.requestDoWithAuth("GET", url)  // user defined function for processing the URL
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // this is mandatory incase to end up the connection with  server to avoid the overlaod on the server
	data, err := ioutil.ReadAll(resp.Body)
	var result Photo
	err = json.Unmarshal(data, &result) // Unmarshal, data from json format to language specific,exactly opposite to the marshal function
	return &result, err
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, randNum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}
func main() {
	os.Setenv("PexelsToken", "lg5UFE3HjhgCBl2vCRPp8oDTJ4tjIVAbicEGXuaTZeWnmQmBNCY81zmX")
	TOKEN := os.Getenv("PexelsToken")
	var c = NewClient(TOKEN)
	result, err := c.SearchPhotos("lion ", 15, 1) // here it begins with c.function which means it is struct method
	if err != nil {
		fmt.Errorf("search error:%v", err)
		os.Exit(1)
	}
	if result.Page == 0 {
		fmt.Errorf("search result wrong")
	}
	fmt.Println(result)
}
