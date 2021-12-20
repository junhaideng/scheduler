package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	ErrNo  int64  `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   Data   `json:"data"`
}

type Data struct {
	ID              int64  `json:"id"`
	LotteryID       string `json:"lottery_id"`
	LotteryName     string `json:"lottery_name"`
	LotteryType     int64  `json:"lottery_type"`
	LotteryImage    string `json:"lottery_image"`
	LotteryDesc     string `json:"lottery_desc"`
	HistoryID       string `json:"history_id"`
	TotalLuckyValue int64  `json:"total_lucky_value"`
	DrawLuckyValue  int64  `json:"draw_lucky_value"`
}

var client = http.Client{}

func getLottery() error {
	url := "https://api.juejin.cn/growth_api/v1/lottery/draw"
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Host", "api.juejin.cn")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62")
	req.Header.Add("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r Response

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&r)
	if err != nil {
		fmt.Println(string(data))
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")
	return encoder.Encode(r)
}

var cookie string

func init() {
	flag.StringVar(&cookie, "c", "", "cookie")
}

func main() {
	flag.Parse()
	err := getLottery()
	if err != nil {
		fmt.Println(err)
	}
}
