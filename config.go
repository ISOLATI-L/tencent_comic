package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
)

type config struct {
	id             string
	chapterPattern string
	interval       int
}

const (
	DEFAULT_ID             string = "623654"
	DEFAULT_CHAPTERPATTERN string = ".*?"
	DEFAULT_INTERVAL       int    = 1000
)

func loadConfig() (cfg config) {
	var configFile *goconfig.ConfigFile
	var err error
	configFile, err = goconfig.LoadConfigFile("config.ini")
	for err != nil {
		if os.IsNotExist(err) {
			os.Create("config.ini")
			configFile, err = goconfig.LoadConfigFile("config.ini")
		} else {
			log.Fatalln(err.Error())
		}
	}
	defer func() {
		err = goconfig.SaveConfigFile(configFile, "config.ini")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	id, err := configFile.GetValue(goconfig.DEFAULT_SECTION, "id")
	if err != nil || id == "" {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"id",
			DEFAULT_ID,
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"id",
			fmt.Sprintf(
				"id为漫画id，如《名侦探柯南》的链接为https://ac.qq.com/Comic/comicInfo/id/%s，id即为%s",
				DEFAULT_ID, DEFAULT_ID,
			),
		)
		id = DEFAULT_ID
	}
	cfg.id = id

	chapterPattern, err := configFile.GetValue(goconfig.DEFAULT_SECTION, "chapterPattern")
	if err != nil || chapterPattern == "" {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			DEFAULT_CHAPTERPATTERN,
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			fmt.Sprintf(
				"chapterPattern为检索章节的正则表达式，默认为%s（即检索所有章节）",
				DEFAULT_CHAPTERPATTERN,
			),
		)
		chapterPattern = DEFAULT_CHAPTERPATTERN
	}
	cfg.chapterPattern = chapterPattern

	intervalStr, err1 := configFile.GetValue(goconfig.DEFAULT_SECTION, "interval")
	interval, err2 := strconv.Atoi(intervalStr)
	if err1 != nil || err2 != nil || interval < 0 {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"interval",
			fmt.Sprint(DEFAULT_INTERVAL),
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"interval",
			fmt.Sprintf(
				"interval为下载每张图片的间隔，单位为毫秒，默认为%d毫秒",
				DEFAULT_INTERVAL,
			),
		)
		interval = DEFAULT_INTERVAL
	}
	cfg.interval = interval

	return cfg
}
