package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/cbroglie/mustache"
)

type WeatherData struct {
	Temperature string
	SkyText     string
	WindDisplay string
}
type Data struct {
	Greeting    string
	RefreshDate string

	Weather WeatherData
}

var ltime time.Time

func main() {
	data := Data{}
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	ltime = time.Now().In(loc)

	data.RefreshDate = ltime.Format("01/02/2006 - 03:04PM") + " GMT-3"
	setGreeting(&data)
	setWeather(&data)

	readFileAndGenerateReadme(&data)
	fmt.Println("README.md generated!")
}

func readFileAndGenerateReadme(data *Data) {
	tmpl, _ := mustache.ParseFile("README.template.md")
	var buf bytes.Buffer
	tmpl.FRender(&buf, *data)
	err := os.WriteFile("README.md", buf.Bytes(), 0644)
	check(err)
}

func setWeather(data *Data) {
	params := url.Values{
		"weasearchstr":  {"Rio de Janeiro, RJ, Paciência"},
		"weadegreetype": {"C"},
		"culture":       {"en-US"},
		"src":           {"outlook"},
	}
	r, err := http.Get("http://weather.service.msn.com/find.aspx?" + params.Encode())
	check(err)
	defer r.Body.Close()

	rdata, _ := ioutil.ReadAll(r.Body)
	var result struct {
		Info []struct {
			Data []struct {
				SkyText     string `xml:"skytext,attr"`
				Temperature string `xml:"temperature,attr"`
				WindDisplay string `xml:"winddisplay,attr"`
			} `xml:"current"`
			Location            string `xml:"weatherlocationname,attr"`
			DegreeType          string `xml:"degreetype,attr"`
		} `xml:"weather"`
	}
	xml.Unmarshal(rdata, &result)
	json.Marshal(result)

	data.Weather.Temperature = result.Info[0].Data[0].Temperature
	data.Weather.SkyText = result.Info[0].Data[0].SkyText
	data.Weather.WindDisplay = result.Info[0].Data[0].WindDisplay
}

func setGreeting(data *Data) {
	hour := ltime.Hour()
	greetings := []string{"☀️ Good morning!", "🌇 Good afternoon!", "🌃 Good evening!"}
	greeting := greetings[0]
	if hour > 12 {
		greeting = greetings[1]
	} else if hour > 18 {
		greeting = greetings[2]
	}
	data.Greeting = greeting
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
