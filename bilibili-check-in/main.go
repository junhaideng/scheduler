package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
)

type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	TTL     int64  `json:"ttl"`
	Data    Data   `json:"data"`
}

type Data struct {
	Following    int64 `json:"following"`
	Follower     int64 `json:"follower"`
	DynamicCount int64 `json:"dynamic_count"`
}

var cookie string

func init() {
	flag.StringVar(&cookie, "c", "", "cookie")
}

func CheckIn() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/web-interface/nav/stat", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res Response

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		return errors.New("Not login, please provide lastest cookie again")
	}
	
	fmt.Printf("Res: %#v\n", res)
	return nil
}

func main() {
	flag.Parse()
	err := CheckIn()
	if err != nil {
		fmt.Println("Check in failed: ", err)
		return
	}
	fmt.Println("Check In Success")
}
