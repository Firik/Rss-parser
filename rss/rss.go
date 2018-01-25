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

import rssHttp "rss_parser/http"

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

func (rss *Rss) DecodeXmlHttpResponse(response *http.Response) {
	url := response.Request.URL.String()
	xmlData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(url, err)
	}
	reader := strings.NewReader(string(xmlData))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&rss)
	if err != nil {
		log.Println(url, err)
	}
}

func (rss *Rss) ProcessUrl(url string, rssItemsChannel chan<- []Item) {
	request := rssHttp.CreateRequest(url)
	response := rssHttp.SendRequest(request)
	rss.DecodeXmlHttpResponse(response)
	response.Body.Close()
	rssItemsChannel <- rss.Channel.Items
}

func (rss *Rss) CombineItems(rssItemsChannel <-chan []Item, readyToGoChannel chan<- bool) {
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

func (rss *Rss) SortItems() {
	sort.Slice(rss.Channel.Items, func(i, j int) bool {
		time1, err := time.Parse(time.RFC1123Z, rss.Channel.Items[i].PubDate)
		if err != nil {
			log.Println("Can't parse time1: ", err)
		}
		time2, err := time.Parse(time.RFC1123Z, rss.Channel.Items[j].PubDate)
		if err != nil {
			log.Println("Can't parse time2: ", err)
		}

		return time1.Before(time2)
	})
}

func (rss *Rss) XmlBytes() *[]byte {
	rss.setDefaultAttributes()
	rssBytes, err := xml.Marshal(rss)
	if err != nil {
		log.Fatal(err)
	}
	return &rssBytes
}

func (rss *Rss) setDefaultAttributes() {
	rss.Version = xml.Attr{
		Name:  xml.Name{Local: "version"},
		Value: "2.0",
	}

	rss.Channel.Title = "Local Rss"
	rss.Channel.Link = "http://localhost"
	rss.Channel.Language = "ru-RU"
}
