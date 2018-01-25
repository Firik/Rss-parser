package file

import (
	"os"
	"log"
	"bufio"
)

func SaveToXmlFile(rssBytes []byte) {
	xmlFile, err := os.Create("output.xml")
	if err != nil {
		log.Fatal(err)
	}
	f := bufio.NewWriter(xmlFile)
	f.Write(rssBytes)
	xmlFile.Close()
}

func GetUrlsFromFile(filename string) *[]string {
	urls := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	f := bufio.NewReader(file)
	for {
		readLine, err := f.ReadString('\n')
		if err != nil && readLine == "" {
			break
		}
		urls = append(urls, readLine)
	}
	file.Close()

	return &urls
}
