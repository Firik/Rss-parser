package rss

import (
	"encoding/xml"
	"net/http"
	"io/ioutil"
	"log"
	"strings"
	"golang.org/x/net/html/charset"
	"sort"
	"time"
)

import "rss_parser/Http"

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
	Version xml.Attr `xml:"version,attr"`
}

type Channel struct {
	XMLName  xml.Name `xml:"channel"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	Language string   `xml:"language"`
	Items    []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name     `xml:"item"`
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description xml.CharData `xml:"description"`
	PubDate     string       `xml:"pubDate"`
}

func (rss *Rss) DecodeXmlHttpResponse(err error, response *http.Response) {
	xmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	reader := strings.NewReader(string(xmlData))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&rss)
	if err != nil {
		log.Println(err)
	}
}

func (rss *Rss) ProcessUrl(url string, rssItemsChannel chan<- []Item) {
	request := Http.CreateRequest(url)
	response, err := Http.SendRequest(request)
	if err != nil {
		errorMessage := "Error while get data from rss " + url + "\n"
		log.Println(errorMessage, err)
		return
	}
	rss.DecodeXmlHttpResponse(err, response)
	response.Body.Close()
	rssItemsChannel <- rss.Channel.Items
}

func (rss *Rss) CombineRssItems(rssItemsChannel <-chan []Item, readyToGoChannel chan<- bool) {
	var i = 0
	for {
		items := <-rssItemsChannel
		for _, item := range items {
			rss.Channel.Items = append(rss.Channel.Items, item)
		}
		i++

		if i == 5 {
			readyToGoChannel <- true
			break
		}
	}
}

func (rss *Rss) SortRssItems() {
	sort.Slice(rss.Channel.Items, func(i, j int) bool {
		time1, err := time.Parse(time.RFC1123Z, rss.Channel.Items[i].PubDate)
		if err != nil {
			log.Println(err)
		}
		time2, err := time.Parse(time.RFC1123Z, rss.Channel.Items[j].PubDate)
		if err != nil {
			log.Println(err)
		}

		return time1.Before(time2)
	})
}

func (rss *Rss) XmlBytes() *[]byte {
	rss.setRssOptions()
	rssBytes, err := xml.Marshal(rss)
	if err != nil {
		log.Fatal(err)
	}
	return &rssBytes
}

func (rss *Rss) setRssOptions() {
	rss.Version = xml.Attr{
		Name:  xml.Name{Local: "version"},
		Value: "2.0",
	}

	rss.Channel.Title = "Local Rss"
	rss.Channel.Link = "http://localhost"
	rss.Channel.Language = "ru-RU"
}
