package file

import (
	"os"
	"log"
	"bufio"
	"path/filepath"
)

func SaveToXmlFile(rssBytes []byte) {
	currentDir := getCurrentDirAbsolute()
	xmlFile, err := os.Create(currentDir + "output.xml")
	if err != nil {
		log.Fatal(err)
	}
	f := bufio.NewWriter(xmlFile)
	f.Write(rssBytes)
	xmlFile.Close()
}

func GetUrlsFromFile(filename string) []string {
	urls := make([]string, 0)
	currentDir := getCurrentDirAbsolute()
	file, err := os.Open(currentDir + filename)
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

	return urls
}

func getCurrentDirAbsolute() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir + "/"
}
