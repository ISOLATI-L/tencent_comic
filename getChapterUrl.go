package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"tencent_comic/regexpPattern"
)

func getChapterUrl(cookies []*http.Cookie, url string) (err error) {
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	html := string(body)
	match := regexpPattern.TitlePattern.FindStringSubmatch(html)
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
	return nil
}
