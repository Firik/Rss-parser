package main

import "os"

import (
	"rss_parser/rss"
	"rss_parser/file"
)

func main() {
	var rssItemsChannel = make(chan []rss.Item)
	var readyToGoChannel = make(chan bool)
	var Rss rss.Rss

	urls := file.GetUrlsFromFile("rss_sources.txt")
	urlsCount := len(*urls)
	for _, url := range *urls {
		go Rss.ProcessUrl(url, rssItemsChannel)
	}

	go Rss.CombineItems(rssItemsChannel, readyToGoChannel, urlsCount)

	for {
		if <-readyToGoChannel {
			Rss.SortItems()
			rssBytes := Rss.XmlBytes()
			file.SaveToXmlFile(*rssBytes)

			os.Exit(0)
		}
	}
}
