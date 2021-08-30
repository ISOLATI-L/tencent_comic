package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type config struct {
	id             string
	chapterPattern string
}

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
	// chapterPatternStr: `.*?FILE.\d+[^(href)]*?`,

	err := Login()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = downloadComic("https://ac.qq.com/ComicView/index/id/623654/cid/1060")
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

	// chaptersUrl, err := getChaptersUrl(cookies, cfg)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// for _, chapterUrl := range chaptersUrl {
	// 	log.Println(chapterUrl)
	// }
	// log.Println(len(chaptersUrl))

	req, err := http.NewRequest(
		"POST",
		"https://ac.qq.com/ComicView/index/id/623654/cid/1060",
		nil,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// for _, cookie := range cookies {
	// 	req.AddCookie(cookie)
	// }
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err.Error())
	}
	savehtml(body)
}
