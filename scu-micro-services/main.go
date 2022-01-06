package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
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

var author = "201648748@qq.com"

var daring = []string{
	"小姐姐", "亲爱的", "宝宝",
}

var sweet = []string{
	"自从遇见了你，余生便是欢喜，余生便都是你",
	"人生只有两次幸运就好，一次遇见你，一次走到底。",
	"如果我不讨你喜欢，你直接爱上我好了。",
	"谁要你的飞吻，有本事真亲过来啊~",
	"有你，我什么都不缺。",
	"我想你应该很忙吧，那就只看前三个字就好了！",
	"你知道我喜欢谁吗，不知道就看看第一个字。",
	"我不喜欢等，我只喜欢你。",
	"想做你的充分必要条件！",
	"对你的喜欢单调递增，没有上限。",
	"希望有一天，我可以成为你的定义域。",
}

var timezone = time.FixedZone("CST", 8*3600)

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

var text = `<!DOCTYPE html>
<html>

<head>
  <meta charset='utf-8'>
  <meta http-equiv='X-UA-Compatible' content='IE=edge'>
  <title>打卡提醒</title>
  <meta name='viewport' content='width=device-width, initial-scale=1'>
  <style>
    .sweet {
      width: 100%;
      display: inline-block;
    }

    p {
      display: inline;
    }

    hr {
      margin: 10px 0;
      border-color: mediumpurple;
      border-width: 0.2px;
      border-style: dashed;
    }

    div {
      padding: 3px 0;
    }

    .from {
      margin-top: 10px;
      font-size: 14px;
      float: right;
      margin-right: 20px;
    }

    .from a {
      text-decoration: none;
      color: black;
    }
  </style>
</head>

<body>
  <div class="sweet">
    <p class="tip">
      每日情话✨:
    </p>
    <p class="content">
      {{.content}} 🥳
    </p>
  </div>
  <hr />
  <div class="check-in">
    {{.dear}}，今天的打卡已经在 {{.time}} 完成了哦，今天也要一块好好学习呀
  </div>
  <div class="from">
    From <a href="mailto:{{.author}}">D先生</a>
  </div>
</body>

</html>`

const UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"

func init() {
	flag.StringVar(&emailUsername, "eu", "", "email username")
	flag.StringVar(&emailPassword, "ep", "", "emial password")
	flag.StringVar(&to, "to", "201648748@qq.com", "users send email to")
	flag.StringVar(&eaiSess, "eai", "", "eai-sess")
	flag.StringVar(&uukey, "uukey", "", "uukey")
	flag.IntVar(&max, "max", 10, "max try times")
	rand.Seed(time.Now().In(timezone).UnixNano())
}

type SubmitResponse struct {
	E int64           `json:"e"`
	M string          `json:"m"`
	D json.RawMessage `json:"d"`
}

// 提交签到内容
func submit(eaiSess, uukey string) (*SubmitResponse, error) {
	form := url.Values{}
	form.Add("zgfxdq", "0")    // 不在中高风险地区
	form.Add("mjry", "0")      // 今日是否接触密接人员
	form.Add("csmjry", "0")    // 近14日内本人/共同居住者是否去过疫情发生场所
	form.Add("szxqmc", "") // 所在校区
	form.Add("sfjzxgym", "1") // 是否接种过新冠疫苗
	form.Add("jzxgymrq", "2021-05-12") // 接种第一剂疫苗时间
	form.Add("sfjzdezxgym", "1") // 是否接种第二剂新冠疫苗
	form.Add("jzdezxgymrq", "2021-06-11") // 接种第二剂疫苗时间
	form.Add("tw", "3")                   // 体温
	form.Add("sfcxtz", "0")               // 没有出现发热、乏力、干咳、呼吸困难等症状
	form.Add("sfjcbh", "0")               // 今日是否接触无症状感染/疑似/确诊人群
	form.Add("sfcxzysx", "0")
	form.Add("qksm", "") // 其他情况
	form.Add("sfyyjc", "0")
	form.Add("jcjgqr", "0")
	form.Add("remark", "")

	form.Add("geo_api_info", `{"type":"complete","position":{"Q":26.988307562935,"R":111.28154595269098,"lng":111.281546,"lat":26.988308},"location_type":"html5","message":"Get geolocation success.Convert Success.Get address success.","accuracy":41883,"isConverted":true,"status":1,"addressComponent":{"citycode":"0739","adcode":"430523","businessAreas":[],"neighborhoodType":"","neighborhood":"","building":"","buildingType":"","street":"凤凰街","streetNumber":"46号","country":"中国","province":"湖南省","city":"邵阳市","district":"邵阳县","towncode":"430523100000","township":"塘渡口镇"},"formattedAddress":"湖南省邵阳市邵阳县塘渡口镇凤凰街46号","roads":[],"crosses":[],"pois":[],"info":"SUCCESS"}`)
	form.Add("address", "湖南省邵阳市邵阳县塘渡口镇凤凰街46号")
	form.Add("area", "湖南省 邵阳市 邵阳县")
	form.Add("province", "湖南省")
	form.Add("city", "邵阳市")
	form.Add("sfzx", "0") // 不在校
	form.Add("bzxyy", "回家")


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
	form.Add("jcjg", "")
	form.Add("hsjcrq", "")
	form.Add("hsjcdd", "")
	form.Add("hsjcjg", "0")
	form.Add("date", time.Now().In(timezone).Format("20060102"))
	// form.Add("uid", "3678")
	form.Add("created", fmt.Sprintf("%d", time.Now().In(timezone).Unix()))
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
	tos := make([]string, 0)
	for _, t := range strings.Split(to, ",") {
		tos = append(tos, strings.TrimSpace(t))
	}
	auth := smtp.PlainAuth("", emailUsername, emailPassword, "smtp.126.com")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s; charset=\"utf-8\"\r\n\r\n %s",
		"D先生 <"+emailUsername+">",
		strings.Join(tos, ","),
		"每日疫情填报",
		"text/html",
		content,
	)
	err := smtp.SendMail("smtp.126.com:25", auth, emailUsername, tos, []byte(msg))
	if err != nil {
		fmt.Println("send email failed: ", err)
		panic("send email failed")
	}
}

func send(typ checkInType) {
	tpl, err := template.New("notice").Parse(text)
	if err != nil {
		sendMail(author, err.Error())
		fmt.Println("", err)
		return
	}
	buf := &strings.Builder{}
	now := time.Now().In(timezone).Format("2006-01-02 15:04:05")

	s := sweet[rand.Intn(len(sweet))]
	dear := daring[rand.Intn(len(daring))]

	tpl.Execute(buf, map[string]string{
		"content": s,
		"time":    now,
		"dear":    dear,
		"author":  author,
	})

	// 打卡成功
	switch typ {
	case CheckInSuccess:
		sendMail(to, buf.String())
	case CheckInFailed:
		sendMail(to, fmt.Sprintf("呜呜呜😭, %s, 今天打卡失败了, 快让D先生给你手动打!!", dear))
	case AlreadyCheckIn:
		fmt.Printf("已经打卡成功，不需要发送邮件啦, 运行时间: %s\n", now)
	}
}

// GitHub Actions 为 0 区，我们这取东八区 => 16pm 打卡
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
			buf.WriteString(time.Now().In(timezone).Format("2006-01-02 15:04:05") + " 出现错误： " + err.Error() + "\n")
			continue
		}
		break
	}
	fmt.Printf("%#v\n", resp)

	if buf.Len() > 0 {
		sendMail(author, buf.String())
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
