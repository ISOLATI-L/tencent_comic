package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type chapter struct {
	name string
	url  string
}

var client *http.Client

func init() {
	var err error
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
	client = &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不进入重定向
		},
		Jar: jar,
	}
}

var CFG config

func main() {
	CFG = loadConfig()
	// chapterPatternStr: `.*?FILE.((104[3-9])|(10[5-9][0-9]))[^(href)]*?`,

	err := Login()
	if err != nil {
		log.Fatalln(err.Error())
	}
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

	// err = getImgUrl("https://ac.qq.com/ComicView/index/id/623654/cid/1060")
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	chaptersUrl, err := getChaptersUrl()
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, chapter := range chaptersUrl {
		fmt.Println(chapter.name)
	}
	err = downloadCommic(chaptersUrl)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
