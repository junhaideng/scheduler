package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
)

var dataPattern = regexp.MustCompile("window.tplData = (.*?);</script>")

var emailUsername string
var emailPassword string
var to string
var city string

func init() {
	flag.StringVar(&emailUsername, "eu", "", "email username")
	flag.StringVar(&emailPassword, "ep", "", "email password")
	// 如果有多个，以 `,` 分隔
	flag.StringVar(&to, "to", "", "send report to user")
	flag.StringVar(&city, "city", "成都", "city")
}

type WeatherData struct {
	// 天气情况
	Weather map[string]string `json:"weather"`
	// 地址
	Position Position `json:"position"`
	//  pm25 值
	PSPm25  PSPm25  `json:"ps_pm25"`
	Feature Feature `json:"feature"`
	// 日期
	Base Base `json:"base"`
	// 24 小时
	The24_HourForecast The24HourForecast `json:"24_hour_forecast"`
	// 15 天
	The15_DayForecast The15DayForecast `json:"15_day_forecast"`
	// 天气情况，40天
	LongDayForecast LongDayForecast `json:"long_day_forecast"`
	// 昨天的
	Yesterday15D Yesterday15D `json:"yesterday_15d"`
	// 部分其他字段省略
}

type Base struct {
	DateShort string `json:"dateShort"`
	Date      string `json:"date"`
	Weekday   string `json:"weekday"`
	Lunar     string `json:"lunar"`
}

type Feature struct {
	Humidity    string `json:"humidity"`
	Wind        string `json:"wind"`
	SunriseTime string `json:"sunriseTime"`
	SunsetTime  string `json:"sunsetTime"`
	Ultraviolet string `json:"ultraviolet"`
}

type LongDayForecast struct {
	Info               []LongDayForecastInfo `json:"info"`
	UpdateTime         string                `json:"update_time"`
	PublishTime        string                `json:"publish_time"`
	Pm25UpdateTime     string                `json:"pm25_update_time"`
	Pm25PublishTime    string                `json:"pm25_publish_time"`
	WeatherPublishTime string                `json:"weather_publish_time"`
	WeatherUpdateTime  string                `json:"weather_update_time"`
	InfoNumBaidu       int64                 `json:"info#num#baidu"`
	Days               string                `json:"days"`
}

type LongDayForecastInfo struct {
	Moonrise               string         `json:"moonrise"`
	WindPowerNight         string         `json:"wind_power_night"`
	TemperatureNight       string         `json:"temperature_night"`
	WindDirectionNight     string         `json:"wind_direction_night"`
	NextFullMoon           string         `json:"next_full_moon"`
	MoonPhaseAngle         string         `json:"moon_phase_angle"`
	NextNewMoon            string         `json:"next_new_moon"`
	Pm25                   *Pm25          `json:"pm25"`
	WindDirectionDay       string         `json:"wind_direction_day"`
	WeatherNight           string         `json:"weather_night"`
	MoonPicNum             string         `json:"moon_pic_num"`
	WindPowerDay           string         `json:"wind_power_day"`
	SunsetTime             string         `json:"sunsetTime"`
	Date                   string         `json:"date"`
	WeatherDay             string         `json:"weather_day"`
	MoonPhase              string         `json:"moon_phase"`
	TemperatureDay         string         `json:"temperature_day"`
	Moonset                string         `json:"moonset"`
	WeatherNightForBeijing string         `json:"weather_night_for_beijing"`
	MoonExposureProportion string         `json:"moon_exposure_proportion"`
	WeatherDayForBeijing   string         `json:"weather_day_for_beijing"`
	SunriseTime            string         `json:"sunriseTime"`
	LimitLine              *HumidityClass `json:"limitLine"`
	Humidity               *HumidityClass `json:"humidity"`
}

type HumidityClass struct {
	Tip  string `json:"tip"`
	Text string `json:"text"`
}

type Pm25 struct {
	ListQuality ListQuality `json:"listQuality"`
	ListTitle   string      `json:"listTitle"`
}

type ListQuality struct {
	ListKey    string `json:"listKey"`
	ListValue  string `json:"listValue"`
	ListAqiVal string `json:"listAqiVal"`
	Site       string `json:"site"`
	Hour       string `json:"hour"`
}

type PSPm25 struct {
	Level  string `json:"level"`
	PSPm25 string `json:"ps_pm25"`
}

type Position struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type The15DayForecast struct {
	Info               []The15DayForecastInfo `json:"info"`
	UpdateTime         string                 `json:"update_time"`
	PublishTime        string                 `json:"publish_time"`
	Pm25UpdateTime     string                 `json:"pm25_update_time"`
	Pm25PublishTime    string                 `json:"pm25_publish_time"`
	WeatherPublishTime string                 `json:"weather_publish_time"`
	WeatherUpdateTime  string                 `json:"weather_update_time"`
	InfoNumBaidu       int64                  `json:"info#num#baidu"`
	Days               string                 `json:"days"`
}

type The15DayForecastInfo struct {
	WeatherDay         string         `json:"weather_day"`
	WindDirectionNight string         `json:"wind_direction_night"`
	TemperatureDay     string         `json:"temperature_day"`
	WindDirectionDay   string         `json:"wind_direction_day"`
	WeatherNight       string         `json:"weather_night"`
	WindPowerDay       string         `json:"wind_power_day"`
	WindPowerNight     string         `json:"wind_power_night"`
	Date               string         `json:"date"`
	SunsetTime         string         `json:"sunsetTime"`
	TemperatureNight   string         `json:"temperature_night"`
	SunriseTime        string         `json:"sunriseTime"`
	Pm25               Pm25           `json:"pm25"`
	LimitLine          *HumidityClass `json:"limitLine"`
	Humidity           *HumidityClass `json:"humidity"`
}

type The24HourForecast struct {
	Info                  []The24HourForecastInfo `json:"info"`
	LOC                   string                  `json:"loc"`
	UpdateTime            string                  `json:"update_time"`
	ResourceName          string                  `json:"resource_name"`
	Delflag               string                  `json:"@delflag"`
	ID                    string                  `json:"@id"`
	Site                  string                  `json:"@site"`
	The24Hourforecastsite string                  `json:"site"`
	WeatherID             string                  `json:"weather_id"`
	The24HourForecastId   string                  `json:"id"`
	ChangeFreq            string                  `json:"changefreq"`
	Updatetime            string                  `json:"@updatetime"`
	PublishTime           string                  `json:"publish_time"`
	Templateid            string                  `json:"@templateid"`
	Type                  string                  `json:"type"`
	LastMod               string                  `json:"lastmod"`
	InfoNumBaidu          int64                   `json:"info#num#baidu"`
}

type The24HourForecastInfo struct {
	Temperature   string      `json:"temperature"`
	Hour          string      `json:"hour"`
	WindDirection string      `json:"wind_direction"`
	Uv            string      `json:"uv"`
	Site          string      `json:"site"`
	UvNum         string      `json:"uv_num"`
	WindPower     string      `json:"wind_power"`
	Weather       string      `json:"weather"`
	WindPowerNum  string      `json:"wind_power_num"`
	Precipitation string      `json:"precipitation"`
	Pm25          ListQuality `json:"pm25"`
}

type Yesterday15D struct {
	WeatherDay         string `json:"weather_day"`
	WindDirectionNight string `json:"wind_direction_night"`
	TemperatureDay     string `json:"temperature_day"`
	WindDirectionDay   string `json:"wind_direction_day"`
	WeatherNight       string `json:"weather_night"`
	WindPowerDay       string `json:"wind_power_day"`
	WindPowerNight     string `json:"wind_power_night"`
	Date               string `json:"date"`
	SunsetTime         string `json:"sunsetTime"`
	TemperatureNight   string `json:"temperature_night"`
	SunriseTime        string `json:"sunriseTime"`
	Pm25               Pm25   `json:"pm25"`
}

func getWeather(city string) (*WeatherData, error) {
	// 发送网络请求，获取源代码
	resp, err := http.Get(fmt.Sprintf("http://weathernew.pae.baidu.com/weathernew/pc?query=%s&srcid=4982", city))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取网页源码
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 进行匹配，获取数据
	match := dataPattern.FindSubmatch(data)
	if len(match) == 0 {
		return nil, errors.New("Do not find weather data")
	}

	// 反序列化数据
	var res WeatherData
	err = json.Unmarshal(match[1], &res)
	if err != nil {
		return nil, err
	}
	return &res, err
}

// 简单的示例，可以进行加工
const tpl = `{{.city}}天气情况：

今天是{{.date}}，农历{{.lunar}}，天气{{.weather}}，伴有{{.wind_direction}}，预测今日{{.precipitation_type}}
紫外线{{.uv}}, {{.uv_info}}，pm2.5指标为{{.pm25}}, 属于{{.pm25_level}}
`

func sendMail(to, content string) error {
	// 这里采用的是 126 邮箱，不同的邮箱 host 设置不同
	auth := smtp.PlainAuth("", emailUsername, emailPassword, "smtp.126.com")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n %s",
		emailUsername,
		to,
		"每日天气",
		"text/plain",
		content,
	)
	err := smtp.SendMail("smtp.126.com:25", auth, emailUsername, []string{to}, []byte(msg))
	if err != nil {
		fmt.Println("send email failed: ", err)
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	data, err := getWeather(city)
	if err != nil {
		fmt.Println(err)
		return
	}

	t, err := template.New("weather").Parse(tpl)
	if err != nil {
		fmt.Println("Parse template failed: ", err)
		return
	}

	var buf = &strings.Builder{}

	err = t.Execute(buf, map[string]string{
		"city":               city,
		"date":               data.Base.Date,
		"lunar":              data.Base.Lunar,
		"weather":            data.Weather["weather"],
		"wind_direction":     data.Weather["wind_direction"],
		"uv":                 data.Weather["uv"],
		"uv_info":            data.Weather["uv_info"],
		"precipitation_type": data.Weather["precipitation_type"],
		"pm25":               data.PSPm25.PSPm25,
		"pm25_level":         data.PSPm25.Level,
	})

	if err != nil {
		fmt.Println("Execute template failed: ", err)
		return
	}

	err = sendMail(to, buf.String())
	if err != nil {
		fmt.Println("send mail failed: ", err)
		panic(err) // 崩溃了 GitHub 会进行提示
	}
	fmt.Println("发送成功")
}
