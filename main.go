package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

const ID string = "623654"

// var ChapterPatternStr = `.*?FILE.\d+[^(href)]*?`

var ChapterPatternStr = ``

type chapter struct {
	name string
	url  string
}

func main() {
	// cookies, err := Login()
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	cookies := make([]*http.Cookie, 0)
	// log.Println("cookies: ", cookies)

	// QRcodeFile, err := os.OpenFile(
	// 	"temp.png",
	// 	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	// 	0666,
	// )
	// QRcodeFile.Write(QRcode)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// QRcodeFile.Close()
	if ChapterPatternStr == "" {
		ChapterPatternStr = ".*?"
	}
	chaptersPattern, err := regexp.Compile(
		fmt.Sprintf(
			`<a\s+target\s*=\s*"_blank"\s+title\s*=\s*"(%s)"`,
			ChapterPatternStr,
		),
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	chaptersUrl, err := getChaptersUrl(cookies, ID, chaptersPattern)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, chapterUrl := range chaptersUrl {
		log.Println(chapterUrl)
	}
	log.Println(len(chaptersUrl))

	req, err := http.NewRequest(
		"GET",
		"https://ac.qq.com/ComicView/index/id/623654/cid/1060",
		nil,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	// resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// log.Println(string(body))
	savehtml(body)
}
