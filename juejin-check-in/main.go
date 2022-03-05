package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var cookie string

var userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1 Edg/98.0.4758.102"

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

// getIndex = function() {
// var e = +Date.now() + Number(("" + Math.random()).slice(2, 8));
// return function() {
// 		return e += 1
// }
func getEventIndex() int64 {
	e := time.Now().UnixMilli() + rand.Int63n(1000000)
	return e + 1
}

func appLogTrace() {
	event := getEventIndex()
	data := fmt.Sprintf(`[{
		"events": [{
			"event": "applog_trace",
			"params": "{\"count\":3,\"state\":\"net\",\"key\":\"log\",\"params_for_special\":\"applog_trace\",\"aid\":2608,\"platform\":\"web\",\"_staging_flag\":1,\"sdk_version\":\"4.2.9\",\"event_index\":%d}",
			"local_time_ms": %d
		}],
		"user": {
			"user_unique_id": "6897956815070594574",
			"user_id": "1275089220809896",
			"user_is_login": true,
			"web_id": "6897956815070594574"
		},
		"header": {
			"app_id": 2608,
			"os_name": "windows",
			"os_version": "10",
			"device_model": "Windows NT 10.0",
			"ab_sdk_version": "90000611",
			"language": "en-US",
			"platform": "Web",
			"sdk_version": "4.2.9",
			"sdk_lib": "js",
			"timezone": 8,
			"tz_offset": -28800,
			"resolution": "1500x1000",
			"browser": "Microsoft Edge",
			"browser_version": "98.0.1108.62",
			"referrer": "https://juejin.cn/user/center/gains",
			"referrer_host": "juejin.cn",
			"width": 1500,
			"height": 1000,
			"screen_width": 1500,
			"screen_height": 1000,
			"utm_medium": "user_center",
			"utm_campaign": "hdjjgame",
			"custom": "{\"student_verify_status\":\"not_student\",\"user_level\":2,\"profile_id\":\"1275089220809896\"}"
		},
		"local_time": %d
	}]`, event, time.Now().UnixMilli(), time.Now().Unix())

	req, err := http.NewRequest(http.MethodPost, "https://mcs.snssdk.com/list", strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func mainSiteStay() {
	event := getEventIndex()
	data := fmt.Sprintf(`[{
		"events": [{
			"event": "main_site_stay",
			"params": "{\"_staging_flag\":0,\"user_id\":\"1275089220809896\",\"page_link\":\"https://juejin.cn/user/settings/account\",\"stay_time\":%d,\"event_index\":%d}",
			"local_time_ms": %d,
			"is_bav": 0,
			"ab_sdk_version": "90000611",
			"session_id": "ffaf3940-7072-40eb-aed8-9edb371742bb"
		}],
		"user": {
			"user_unique_id": "6897956815070594574",
			"user_id": "1275089220809896",
			"user_is_login": true,
			"web_id": "6897956815070594574"
		},
		"header": {
			"app_id": 2608,
			"os_name": "windows",
			"os_version": "10",
			"device_model": "Windows NT 10.0",
			"ab_sdk_version": "90000611",
			"language": "en-US",
			"platform": "Web",
			"sdk_version": "4.2.9",
			"sdk_lib": "js",
			"timezone": 8,
			"tz_offset": -28800,
			"resolution": "1500x1000",
			"browser": "Microsoft Edge",
			"browser_version": "98.0.1108.62",
			"referrer": "https://juejin.cn/user/center/gains",
			"referrer_host": "juejin.cn",
			"width": 1500,
			"height": 1000,
			"screen_width": 1500,
			"screen_height": 1000,
			"utm_medium": "user_center",
			"utm_campaign": "hdjjgame",
			"custom": "{\"student_verify_status\":\"not_student\",\"user_level\":2,\"profile_id\":\"1275089220809896\"}"
		},
		"local_time": %d
	}]`, rand.Int63n(2000)+1000, event, time.Now().UnixMilli(), time.Now().Unix())

	req, err := http.NewRequest(http.MethodPost, "https://mcs.snssdk.com/list", strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func CheckIn(cookie string) (*CheckInResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "https://api.juejin.cn/growth_api/v1/check_in", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", cookie)
	req.Header.Add("User-Agent", userAgent)

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

	// 添加日志，尝试避免被后台监控到使用脚本
	// 不知是否可行?
	appLogTrace()
	mainSiteStay()

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
