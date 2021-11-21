package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"
)

// 邮件用户
var emailUsername string
var emailPassword string

// 发送给哪些用户
var to string

var eaiSess string
var uukey string

var max int

// 签到情况
type checkInType int

const (
	CheckInSuccess checkInType = iota
	CheckInFailed
	AlreadyCheckIn
)

var daring = []string{
	"小姐姐", "亲爱的", "宝宝",
}

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return nil
	},
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
	},
}

const UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"

func init() {
	flag.StringVar(&emailUsername, "eu", "", "email username")
	flag.StringVar(&emailPassword, "ep", "", "emial password")
	flag.StringVar(&to, "to", "201648748@qq.com", "users send email to")
	flag.StringVar(&eaiSess, "eai", "", "eai-sess")
	flag.StringVar(&uukey, "uukey", "", "uukey")
	flag.IntVar(&max, "max", 10, "max try times")
	rand.Seed(time.Now().UnixNano())
}

type SubmitResponse struct {
	E int64           `json:"e"`
	M string          `json:"m"`
	D json.RawMessage `json:"d"`
}

// 提交签到内容
func submit(eaiSess, uukey string) (*SubmitResponse, error) {
	form := url.Values{}
	form.Add("zgfxdq", "0") // 不在中高风险地区
	form.Add("mjry", "0")   // 今日是否接触密接人员
	form.Add("csmjry", "0") // 近14日内本人/共同居住者是否去过疫情发生场所
	form.Add("szxqmc", "望江校区")
	form.Add("sfjzxgym", "1")
	form.Add("jzxgymrq", "2021-05-12")
	form.Add("sfjzdezxgym", "1")
	form.Add("jzdezxgymrq", "2021-06-11")
	form.Add("tw", "3")     // 体温
	form.Add("sfcxtz", "0") // 没有出现发热、乏力、干咳、呼吸困难等症状
	form.Add("sfjcbh", "0") // 今日是否接触无症状感染/疑似/确诊人群
	form.Add("sfcxzysx", "0")
	form.Add("qksm", "") // 其他情况
	form.Add("sfyyjc", "0")
	form.Add("jcjgqr", "0")
	form.Add("remark", "")
	form.Add("address", "四川省成都市武侯区望江路街道四川大学四川大学望江校区")
	form.Add("geo_api_info",
		`{"type":"complete","position":{"Q":30.634781629775,"R":104.08064317491403,"lng":104.080643,"lat":30.634782},"location_type":"html5","message":"Get sdkLocation failed.Get geolocation success.Convert Success.Get address success.","accuracy":50,"isConverted":true,"status":1,"addressComponent":{"citycode":"028","adcode":"510107","businessAreas":[{"name":"小天竺","id":"510107","location":{"Q":30.639354,"R":104.06894199999999,"lng":104.068942,"lat":30.639354}},{"name":"跳伞塔","id":"510107","location":{"Q":30.636149,"R":104.07122400000003,"lng":104.071224,"lat":30.636149}}],"neighborhoodType":"科教文化服务;学校;高等院校","neighborhood":"四川大学","building":"","buildingType":"","street":"新南路","streetNumber":"51号","country":"中国","province":"四川省","city":"成都市","district":"武侯区","township":"望江路街道"},"formattedAddress":"四川省成都市武侯区望江路街道四川大学四川大学望江校区","roads":[],"crosses":[],"pois":[],"info":"SUCCESS"}`)
	form.Add("area", "四川省 成都市 武侯区")
	form.Add("province", "四川省")
	form.Add("city", "成都市")
	form.Add("sfzx", "1") // 是否在校
	form.Add("sfjcwhry", "0")
	form.Add("sfjchbry", "0")
	form.Add("sfcyglq", "0")
	form.Add("gllx", "")
	form.Add("glksrq", "")
	form.Add("jcbhlx", "")
	form.Add("jcbhrq", "")
	form.Add("bztcyy", "1")
	form.Add("sftjhb", "0")
	form.Add("sftjwh", "0")
	form.Add("szcs", "")
	form.Add("szgj", "")
	form.Add("bzxyy", "")
	form.Add("jcjg", "")
	form.Add("hsjcrq", "")
	form.Add("hsjcdd", "")
	form.Add("hsjcjg", "0")
	form.Add("date", time.Now().Format("20060102"))
	// form.Add("uid", "3678")
	form.Add("created", fmt.Sprintf("%d", time.Now().Unix()))
	form.Add("jcqzrq", "")
	form.Add("sfjcqz", "")
	form.Add("szsqsfybl", "0")
	form.Add("sfsqhzjkk", "0")
	form.Add("sqhzjkkys", "")
	form.Add("sfygtjzzfj", "0")
	form.Add("gtjzzfjsj", "")
	form.Add("created_uid", "0")
	// form.Add("id", "46708101")
	form.Add("gwszdd", "")
	form.Add("sfyqjzgc", "")
	form.Add("jrsfqzys", "")
	form.Add("jrsfqzfy", "")
	form.Add("szgjcs", "")
	form.Add("ismoved", "0")

	req, err := http.NewRequest(http.MethodPost, "https://wfw.scu.edu.cn/ncov/wap/default/save", strings.NewReader(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", fmt.Sprintf("eai-sess=%s; UUkey=%s", eaiSess, uukey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status: " + resp.Status)
	}

	var res SubmitResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		// decode 错误，返回的数据不是 json
		// 有时候 POST 之后返回 200 text/plain
		return nil, err
	}
	return &res, nil
}

func sendMail(to, content string) {
	auth := smtp.PlainAuth("", emailUsername, emailPassword, "smtp.126.com")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n %s",
		"D先生 <"+emailUsername+">",
		to,
		"每日打卡",
		"text/plain",
		content,
	)
	err := smtp.SendMail("smtp.126.com:25", auth, emailUsername, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("send email failed: ", err)
		panic("send email failed")
	}
}

func send(typ checkInType) {
	dear := daring[rand.Intn(len(daring))]
	// 打卡成功
	switch typ {
	case CheckInSuccess:
		sendMail(to, fmt.Sprintf("%s, 今天的打卡已经在 %s 完成了哦, 笔芯", dear, time.Now().Format("2006-01-02 15:04:05")))
	case CheckInFailed:
		sendMail(to, fmt.Sprintf("呜呜呜, %s, 今天打卡失败了, 快让D先生给你手动打!!", dear))
	case AlreadyCheckIn:
		fmt.Println("已经打卡成功，不需要发送邮件啦")
	}
}

// GitHub Actions 为 0 区，我们这取东八区 => 16pm 打卡，实际上应该再晚一点
func main() {
	flag.Parse()
	count := 0
	var err error
	var resp *SubmitResponse
	var buf strings.Builder
	for count < max {
		count++
		resp, err = submit(eaiSess, uukey)
		if err != nil {
			fmt.Println("submit err: ", err)
			buf.WriteString(time.Now().Format("2006-01-02 15:04:05") + " 出现错误： " + err.Error() + "\n")
			continue
		}
		break
	}
	fmt.Printf("%#v\n", resp)

	if buf.Len() > 0 {
		sendMail("201648748@qq.com", buf.String())
	}
	// 成功
	if strings.Contains(resp.M, "已经") {
		send(AlreadyCheckIn)
	} else if strings.Contains(resp.M, "成功") {
		send(CheckInSuccess)
	} else {
		send(CheckInFailed)
	}
}