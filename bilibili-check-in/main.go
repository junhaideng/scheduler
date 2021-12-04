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

// 仅截取一部分字段
type Data struct {
	IsLogin bool    `json:"isLogin"`
	Money   float64 `json:"money"`
}

var cookie string

func init() {
	flag.StringVar(&cookie, "c", "", "cookie")
}

func CheckIn() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.bilibili.com/x/web-interface/nav", nil)
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

	if !res.Data.IsLogin {
		return errors.New("Not login, please provide lastest cookie again")
	}
	fmt.Println("Has money: ", res.Data.Money)
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
