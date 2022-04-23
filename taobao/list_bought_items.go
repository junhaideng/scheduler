package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
)

// 邮件用户
var emailUsername string
var emailPassword string

// 发送给哪些用户
var to string
var itemId int64
var orderId string

var cookie string
var author string

func init() {
	flag.StringVar(&emailUsername, "eu", "", "email username")
	flag.StringVar(&emailPassword, "ep", "", "emial password")
	flag.StringVar(&to, "to", "201648748@qq.com", "users send email to")
	flag.Int64Var(&itemId, "item_id", 2577079479835126857, "item id")
	flag.StringVar(&orderId, "order_id", "2577079479834126857", "order id")
	flag.StringVar(&cookie, "cookie", "", "cookie")
	flag.StringVar(&author, "author", "201648748@qq.com", "author email")
}

type Data struct {
	MainOrders []MainOrders `json:"mainOrders,omitempty"`
	Error      string       `json:"error,omitempty"`
}

type Operation struct {
	Style   string  `json:"style"`
	Text    string  `json:"text"`
	URL     *string `json:"url,omitempty"`
	Action  *string `json:"action,omitempty"`
	DataURL *string `json:"dataUrl,omitempty"`
}

type MainOrders struct {
	ID        string     `json:"id"`
	SubOrders []SubOrder `json:"subOrders"`
}

type SubOrder struct {
	ID        int64       `json:"id"`
	ItemInfo  ItemInfo    `json:"itemInfo"`
	PriceInfo PriceInfo   `json:"priceInfo"`
	Quantity  string      `json:"quantity"`
	Operation []Operation `json:"operations"`
}
type Extra struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Visible string `json:"visible"`
}

type ItemInfo struct {
	Extra     []Extra `json:"extra"`
	ID        int64   `json:"id"`
	Pic       string  `json:"pic"`
	SkuID     int64   `json:"skuId"`
	Title     string  `json:"title"`
	XtCurrent bool    `json:"xtCurrent"`
}

type PriceInfo struct {
	Original  string `json:"original"`
	RealTotal string `json:"realTotal"`
}

var dataPattern = regexp.MustCompile(`(?m)var data = JSON.parse\('(.*?)'\)`)

func getData(cookie string) {
	url := "https://buyertrade.taobao.com/trade/itemlist/list_bought_items.htm"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("cookie", cookie)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// os.WriteFile("item.html", data, 0666)
	res := dataPattern.FindSubmatch(data)
	if len(res) < 2 {
		fmt.Println("未找到")
		_ = errors.New("没有找到")
		sendMail(author, "更新淘宝 cookie")
		return
	}
	res[1] = bytes.ReplaceAll(res[1], []byte(`\`), []byte{})

	// fmt.Println(string(res[1]))
	var item Data
	err = json.Unmarshal(res[1], &item)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("%#v\n", item)
	// f, err := os.Create("item.json")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// encoder := json.NewEncoder(f)
	// encoder.SetEscapeHTML(false)
	// encoder.SetIndent(" ", "  ")
	// encoder.Encode(item)

	for _, order := range item.MainOrders {
		if order.ID == orderId {
			for _, item := range order.SubOrders {
				if item.ID == itemId {
					// fmt.Println(item)
					// item.Operation[0].Text
					flag := false
					for _, op := range item.Operation {
						if strings.EqualFold(op.Text, "u672Au53D1u8D27") {
							flag = true
						}
					}
					if !flag {
						sendMail(to, "已经发货啦，可以关闭定时任务了")
					} else {
						fmt.Println("尚未发货")
						d, _ := json.MarshalIndent(item, " ", "  ")
						fmt.Printf("%v\n", string(d))
					}
				}
			}
		}
	}
}

func sendMail(to, content string) {
	tos := make([]string, 0)
	for _, t := range strings.Split(to, ",") {
		tos = append(tos, strings.TrimSpace(t))
	}
	auth := smtp.PlainAuth("", emailUsername, emailPassword, "smtp.126.com")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s; charset=\"utf-8\"\r\n\r\n %s",
		"通知 <"+emailUsername+">",
		strings.Join(tos, ","),
		"淘宝发货状态",
		"text/html",
		content,
	)
	err := smtp.SendMail("smtp.126.com:25", auth, emailUsername, tos, []byte(msg))
	if err != nil {
		fmt.Println("send email failed: ", err)
		panic("send email failed")
	}
}
func main() {
	flag.Parse()
	getData(cookie)
}
