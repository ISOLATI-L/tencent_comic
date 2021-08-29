package main

import (
	"io"
	"log"
	"net/http"
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
	var html string
	{
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		savehtml(body)
		html = string(body)
	}
	log.Println(html)
	return nil
}
