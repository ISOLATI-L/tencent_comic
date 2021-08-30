package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var picStartPattern *regexp.Regexp
var picEndPattern *regexp.Regexp
var picPattern *regexp.Regexp
var suffixPattern *regexp.Regexp

func init() {
	var err error
	picStartPattern, err = regexp.Compile(
		`"picture"\s*:\s*\[`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	picEndPattern, err = regexp.Compile(
		`"ads"\s*:\s*{`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	picPattern, err = regexp.Compile(
		`"url"\s*:\s*"(.*?)"`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	suffixPattern, err = regexp.Compile(
		`(\.[(a-z)|(A-Z)|(0-9)]*)[/$]`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func downloadCommic(chaptersUrl []chapter) error {
	var wg sync.WaitGroup
	for _, chapterUrl := range chaptersUrl {
		_, err := os.Stat(chapterUrl.name)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(chapterUrl.name, 0777)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		info_str, err := getImgUrl(chapterUrl.url)
		if err != nil {
			return err
		}
		startl := picStartPattern.FindStringIndex(info_str)
		if len(startl) == 0 {
			continue
		}
		start := startl[0]
		endl := picEndPattern.FindStringIndex(info_str)
		if len(endl) == 0 {
			continue
		}
		end := endl[0]
		matches := picPattern.FindAllStringSubmatch(info_str[start:end], -1)
		wg.Add(len(matches))
		for idx, match := range matches {
			go func(url string, title string) {
				defer wg.Done()
				// 出现错误最多尝试5次
				for i := 0; i < 5; i++ {
					req, err := http.NewRequest(
						"GET",
						url,
						nil,
					)
					if err != nil {
						log.Println("error: ", err.Error())
						continue
					}
					resp, err := client.Do(req)
					if err != nil {
						log.Println("error: ", err.Error())
						continue
					}
					err = saveImg(resp.Body, title, url, idx)
					if err != nil {
						log.Println("error: ", err.Error())
						continue
					}
					break
				}
			}(strings.ReplaceAll(match[1], "\\", ""), chapterUrl.name)
			time.Sleep(time.Duration(CFG.interval) * time.Millisecond)
		}
	}
	wg.Wait()
	return nil
}

func saveImg(img io.ReadCloser, title string, url string, idx int) error {
	defer img.Close()
	suffix := ""
	matches := suffixPattern.FindAllStringSubmatch(url, -1)
	if len(matches) > 0 {
		suffix = matches[len(matches)-1][1]
	}
	fileName := fmt.Sprintf("%s/%02d%s", title, idx, suffix)
	log.Println("正在下载", fileName)
	imgFile, err := os.OpenFile(
		fileName,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	_, err = io.Copy(imgFile, img)
	if err != nil {
		return err
	}
	return nil
}
