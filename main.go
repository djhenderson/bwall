// make_http_request.go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	verbose     = false
	logFileName = `D:/Users/Doug/Repositories/GitHub/bwall/bwall.log`
)

// getImageLink parses Bing and gets img url
// https://www.devdungeon.com/content/web-scraping-go
func getImageURL(domain string) string {
	// Make HTTP GET request
	if verbose {
		fmt.Println("Domain          :", domain)
	}
	response, err := http.Get(domain)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer response.Body.Close()

	if verbose {
		fmt.Println("Status          :", response.Status)
		fmt.Println("StatusCode      :", response.StatusCode)
		fmt.Println("Proto           :", response.Proto)
		fmt.Println("ProtoMajor      :", response.ProtoMajor)
		fmt.Println("ProtoMinor      :", response.ProtoMinor)
		fmt.Println("Header          :")
		for k, v := range response.Header {
			fmt.Println("Header          :", k, ":", v)
		}
		fmt.Println()
		fmt.Println("ContentLength   :", response.ContentLength)
		fmt.Println("TransferEncoding:", response.TransferEncoding)
		fmt.Println("Uncompressed    :", response.Uncompressed)
		//fmt.Println("TLS             :", response.TLS) // too long
	}

	body_bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if len(body_bytes) <= 0 {
		fmt.Println("body_bytes - len, cap:", len(body_bytes), cap(body_bytes))
		log.Fatalln("Empty body")
	}

	expr := `<link(( id="(?P<id>bgLink)")|( rel="(?P<rel>[^"]+)")|( href="(?P<href>[^"]+)")|( as="(?P<as>[^"]+)"))+ />`
	re, err := regexp.Compile(expr)

	subexpNames := re.SubexpNames()

	var page string
	if verbose {
		fmt.Printf("SubexpNames: %q\n", subexpNames)

		fmt.Println("Match:", re.Match(body_bytes))
		fmt.Println()
		fmt.Println("FindSubmatch:")
	}
	for i, v := range re.FindSubmatch(body_bytes) {
		if len(subexpNames[i]) > 0 {
			if verbose {
				fmt.Printf("%2d %-4s: %s\n", i, subexpNames[i], string(v))
			}
			if subexpNames[i] == "href" {
				page = string(v)
			}
		}
	}

	url := domain + page
	if verbose {
		fmt.Println()
		fmt.Println("url:", url)
	}

	return url
}

// exists returns whether the given file or directory exists
// https://stackoverflow.com/a/10510783
func exists(path string) (result bool) {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		result = true
	} else if os.IsNotExist(err) {
		result = false
	} else {
		log.Fatalln(err.Error())
	}
	return result
}

// downloadImg saves image to ./.data directory
// and returns file path
// https://stackoverflow.com/questions/22417283/save-an-image-from-url-to-file
func downloadImg(url string) string {
	dir := "D:/Users/Doug/Pictures/Wallpaper/Bing"

	if verbose {
		fmt.Println("url:", url)
		fmt.Println("dir:", dir)
	}

	found := exists(dir)
	if !found {
		err := os.MkdirAll(dir, 0766)
		if err != nil {
			log.Println("dir:", dir)
			log.Fatalln("Couldn't Create Directory:", err)
		}
	}

	// get image
	response, err := http.Get(url)
	if err != nil {
		log.Println("url:", url)
		log.Fatalln("Couldn't Download Image:", err)
	}
	defer response.Body.Close()

	lastSlash := strings.LastIndex(url, "/")
	fileName := url[lastSlash+1:]
	goodFilename := false
	if strings.HasSuffix(fileName, ".jpg") {
		// log.Println("HasSuffix of '.jpg'")
		goodFilename = true
	} else if strings.HasPrefix(fileName, "th?id=") && strings.Contains(fileName, ".jpg&amp;") {
		fileName = fileName[6:]
		// log.Println("fileName:", fileName)
		endX := strings.Index(fileName, ".jpg&amp;")
		fileName = fileName[:endX+4]
		// log.Println("fileName:", fileName)
		goodFilename = true
	}
	if !goodFilename {
		log.Fatalln("Cannot guess local filename.")
	}
	filePath := dir + "/" + fileName

	// create and/or (re)open file
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("filePath:", filePath)
		log.Fatalln("Couldn't Create and (re)open File:", err)
	}
	defer file.Close()

	// Copy to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatalln("Couldn't Save Image:", err)
	}

	return filePath
}

//func setImageAsWallpaper() {
//url := getImageURL()
//file := downloadImg(url)
//wallpaper.SetFromFile(file)
//fmt.Println(file)
//}

//// https://gist.github.com/ryanfitz/4191392
//func routine() {
//setImageAsWallpaper()
//for range time.Tick(24 * time.Hour) {
//setImageAsWallpaper()
//}
//}

func routine() {
	domain := "https://www.bing.com"
	url := getImageURL(domain)
	log.Println("Image found at:", url)
	filePath := downloadImg(url)
	log.Println("Image saved to:", filePath)
}

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("")

	routine()
}
