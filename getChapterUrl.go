package main

import (
	"io"
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = savehtml(body)
	if err != nil {
		return err
	}
	// html := string(body)
	return nil
}
