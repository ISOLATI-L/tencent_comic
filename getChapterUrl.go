package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func getChaptersUrl(cookies []*http.Cookie, id string, chaptersPattern *regexp.Regexp) (chaptersUrl []chapter, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://ac.qq.com/Comic/comicInfo/id/%s",
			id,
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

	startPattern, err := regexp.Compile(
		`class\s*=\s*"chapter-page-all works-chapter-list"`,
	)
	if err != nil {
		return nil, err
	}
	startl := startPattern.FindStringIndex(html)
	if len(startl) == 0 {
		return nil, errors.New("can not find the start of chapter list")
	}
	start := startl[0]
	endPattern, err := regexp.Compile(
		`class\s*=\s*"chapter-page-new works-chapter-list"`,
	)
	if err != nil {
		return nil, err
	}
	endl := endPattern.FindStringIndex(html)
	if len(endl) == 0 {
		return nil, errors.New("can not find the end of chapter list")
	}
	end := endl[0]

	// matches := ChapterPattern.FindAllStringSubmatch(html, -1)
	indexes := chaptersPattern.FindAllStringIndex(html[start:end], -1)
	// matches := chapterPattern.FindStringSubmatch(html)
	chapterPattern, err := regexp.Compile(
		fmt.Sprintf(
			`"(/ComicView/index/id/%s/cid/\d*?)"`,
			id,
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
