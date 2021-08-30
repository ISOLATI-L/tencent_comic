package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

var startPattern *regexp.Regexp
var endPattern *regexp.Regexp

func init() {
	var err error
	startPattern, err = regexp.Compile(
		`class\s*=\s*"chapter-page-all works-chapter-list"`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	endPattern, err = regexp.Compile(
		`class\s*=\s*"chapter-page-new works-chapter-list"`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func getChaptersUrl(cookies []*http.Cookie, cfg config) (chaptersUrl []chapter, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://ac.qq.com/Comic/comicInfo/id/%s",
			cfg.id,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var html string
	{
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		savehtml(body)
		html = string(body)
	}

	startl := startPattern.FindStringIndex(html)
	if len(startl) == 0 {
		return nil, errors.New("can not find the start of chapter list")
	}
	start := startl[0]
	endl := endPattern.FindStringIndex(html)
	if len(endl) == 0 {
		return nil, errors.New("can not find the end of chapter list")
	}
	end := endl[0]

	chaptersPattern, err := regexp.Compile(
		fmt.Sprintf(
			`<a\s+target\s*=\s*"_blank"\s+title\s*=\s*"(%s)"\s*href=`,
			cfg.chapterPattern,
		),
	)
	if err != nil {
		return nil, err
	}
	indexes := chaptersPattern.FindAllStringIndex(html[start:end], -1)
	// matches := chapterPattern.FindStringSubmatch(html)
	chapterPattern, err := regexp.Compile(
		fmt.Sprintf(
			`"(/ComicView/index/id/%s/cid/\d*?)"`,
			cfg.id,
		),
	)
	if err != nil {
		return nil, err
	}
	chaptersUrl = make([]chapter, len(indexes))
	for i, index := range indexes {
		name := chaptersPattern.FindStringSubmatch(html[start+index[0] : start+endl[0]])
		url := chapterPattern.FindStringSubmatch(html[start+index[0] : start+endl[0]])
		chaptersUrl[i] = chapter{
			name: name[1],
			url:  "https://ac.qq.com" + url[1],
		}
		// log.Println(chaptersUrl[i])
	}
	// log.Println(len(indexes))
	return chaptersUrl, nil
}
