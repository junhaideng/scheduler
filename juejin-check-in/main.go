package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var cookie string

func init() {
	flag.StringVar(&cookie, "c", "", "juejin cookie")
}

type CheckInResponse struct {
	ErrNo  int64  `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   Data   `json:"data"`
}

type Data struct {
	IncrPoint int64 `json:"incr_point"`
	SumPoint  int64 `json:"sum_point"`
}

func CheckIn(cookie string) (*CheckInResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "https://api.juejin.cn/growth_api/v1/check_in", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var res CheckInResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err

	}
	return &res, nil
}

func main() {
	flag.Parse()

	res, err := CheckIn(cookie)
	if err != nil {
		fmt.Println("Checkin in failed: ", err)
		return
	}
	if res.ErrNo == 0 {
		fmt.Println("Check in succeed")
		return
	}
	if res.ErrNo == 15001 {
		fmt.Println("Already check in")
		return
	}
	fmt.Printf("%#v\n", res)
	panic("Check in encounter a problem")
}
