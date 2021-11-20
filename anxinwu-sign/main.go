package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type response struct {
	AddCoins int `json:"addCoins,omitempty"`
	Days     int `json:"days,omitempty"`
}

var auth string

func init() {
	flag.StringVar(&auth, "auth", "", "鉴权token")
}

func main() {
	flag.Parse()
	const url = "https://aixinwu.sjtu.edu.cn/wechat/continousLoginDay"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println(auth)

	req.Header.Add("Authorization", auth)
	req.Header.Add("Host", "aixinwu.sjtu.edu.cn")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.52")

	var client = http.Client{}
	maxTimes := 10
	count := 0
	for count < maxTimes {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		resData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
		var resp response
		json.Unmarshal(resData, &resp)
		if resp.AddCoins != 0 {
			fmt.Printf("获得爱心币: %d, 已经连续登录: %d天\n", resp.AddCoins, resp.Days)
			break
		}
		fmt.Println(string(resData))
		count++
	}

}
