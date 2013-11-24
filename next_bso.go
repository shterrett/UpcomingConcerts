package main

import (
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xml"
  "net/http"
  "fmt"
  "io/ioutil"
  "encoding/json"
)

func Works(url string) string {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
  }
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
  }
  defer resp.Body.Close()
  page, err := gokogiri.ParseHtml(body)
  if err != nil {
    fmt.Println(err)
  }
  return Pieces(page)
}

func Pieces(details xml.Node) string {
  pieces, _ := details.Search(".//div[@class='program-media-collapse']/h3")
  var piecesString string
  piecesString = "<ul class=\"works\">"
  for _, piece := range(pieces) {
    piecesString += "<li>"
    piecesString += piece.Content()
    piecesString += "</li>"
  }
  piecesString += "</ul>"
  return piecesString
}

func Title(performance xml.Node) string {
  title, _ := performance.Search(".//div[@class='performance-title']")
  return title[0].Content()
}

func Time(performance xml.Node) string {
  time, _ := performance.Search(".//span[@class='performance-time']")
  return time[0].Content()
}

func Day(performance xml.Node) string {
  day, _ := performance.Search(".//span[contains(@class, 'performance-day')]")
  return day[0].Content()
}

func PerformanceList() xml.Node {
  resp, err := http.Get("http://www.bso.org/Performance/Listing")
  if err != nil {
    fmt.Println(err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
  }
  page, err := gokogiri.ParseHtml(body)
  if err != nil {
    fmt.Println(err)
  }
  list, _ := page.Search(".//ul[@id='listOfPerformances']")
  return list[0]
}

type Concert struct {
  Title string `json:"title"`
  Works string `json:"works"`
  Day string   `json:"day"`
  Time string  `json:"time"`
  Link string  `json:"link"`
  Id int       `json:"id"`
}

type Concerts struct {
  Concerts []Concert `json:"concerts"`
}

func SetWorks(concert *Concert, link string, status chan int) {
  works := Works(link)
  (*concert).Works = works
  status <- 1
}

func Link(performance xml.Node) string {
    anchor, err := performance.Search(".//a")
    if err != nil {
      fmt.Println(err)
    }
    return "http://www.bso.org" + anchor[0].Attr("href")
}

func BuildConcertList() []Concert {
  list := PerformanceList()
  performances, _ := list.Search(".//li")
  numberPerformances := len(performances) - 1
  upcomingPerformances := make([]Concert, numberPerformances)
  done := make(chan int)
  for i := 0; i < numberPerformances; i++ {
    performance := performances[i]
    link := Link(performance)
    concert := Concert{Title: Title(performance),
                       Time: Time(performance),
                       Id: i,
                       Link: link,
                       Day: Day(performance)}
    upcomingPerformances[i] = concert
    go SetWorks(&upcomingPerformances[i], link, done)
  }
  j := 0
  for j < numberPerformances {
    j += <-done
  }
  return upcomingPerformances
}

type ConcertsHandler struct{}
func (c ConcertsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
  concerts := Concerts{}
  concerts.Concerts = BuildConcertList()
  concertsJson, _ := json.Marshal(concerts)
  writer.Header().Set("Content-Type", "application/json")
  writer.Write(concertsJson)
}


func main() {
  http.Handle("/", http.FileServer(http.Dir("./public")))
  http.Handle("/concerts", ConcertsHandler{})
  http.ListenAndServe("localhost:4000", nil)
}
