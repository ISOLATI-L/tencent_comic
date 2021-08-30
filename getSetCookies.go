package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

var cookieNamePattern *regexp.Regexp
var cookieValuePattern *regexp.Regexp

func init() {
	var err error
	cookieNamePattern, err = regexp.Compile(
		`^\s*(.*?)=`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	cookieValuePattern, err = regexp.Compile(
		`=(.*?)\s*$`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getSetCookies(resp *http.Response) []*http.Cookie {
	headers := resp.Header
	setCookies := make([]*http.Cookie, 0)
	for idx, header := range headers {
		if idx == "Set-Cookie" {
			for _, head := range header {
				cookies := strings.Split(head, ";")
				for _, cookie := range cookies {
					// log.Println(cookie)
					matches := cookieValuePattern.FindStringSubmatch(cookie)
					// log.Println(matches)
					if len(matches) > 0 {
						value := matches[1]
						if len(value) > 0 && value != "/" {
							name := cookieNamePattern.FindStringSubmatch(cookie)[1]
							setCookies = append(setCookies, &http.Cookie{
								Name:  name,
								Value: value,
							})
						}
					}
				}
			}
		}
	}
	return setCookies
}
