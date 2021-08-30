package main

import (
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const (
	DEFAULT_DISPLAYERHRIGHT uint   = 108
	DEFAULT_SLICEWIDTH      uint   = 96
	DEFAULT_SLICEHEIGHT     uint   = 100
	DEFAULT_OFFSETX         int    = 0
	DEFAULT_OFFSETY         int    = 0
	DEFAULT_MAXGOROUTINENUM uint   = 20
	DEFAULT_NEEDELECTRIC    uint8  = 0
	DEFAULT_OUTPUT          string = "output"
)

func loadConfig() config {
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

	var result config

	id, err := configFile.GetValue(goconfig.DEFAULT_SECTION, "id")
	if err != nil || id == "" {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"id",
			"623654",
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"id",
			"id为漫画id，如《名侦探柯南》的链接为https://ac.qq.com/Comic/comicInfo/id/623654，id即为623654",
		)
		id = "623654"
	}
	result.id = id

	chapterPattern, err := configFile.GetValue(goconfig.DEFAULT_SECTION, "chapterPattern")
	if err != nil || id == "" {
		configFile.SetValue(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			".*?",
		)
		configFile.SetKeyComments(
			goconfig.DEFAULT_SECTION,
			"chapterPattern",
			"chapterPattern为检索章节的正则表达式，默认为.*?（即检索所有章节）",
		)
		chapterPattern = ".*?"
	}
	result.chapterPattern = chapterPattern

	return result
}
