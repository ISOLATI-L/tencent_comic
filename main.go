package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"tencent_comic/regexpPattern"
)

const URL string = "https://ac.qq.com/Comic/comicInfo/id/623654"

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
	getChapterUrl(cookies, URL)

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
	match := regexpPattern.TitlePattern.FindStringSubmatch(string(body))
	if len(match) > 0 {
		htmlFile, err := os.OpenFile(
			match[1]+".html",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0666,
		)
		if err != nil {
			log.Fatalln(err.Error())
		}
		htmlFile.Write(body)
		htmlFile.Close()
	}
}
