package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// var jsAddrPattern *regexp.Regexp
var DATAPattern *regexp.Regexp
var NONCEPattern *regexp.Regexp

var DecodePattern1 *regexp.Regexp
var DecodePattern2 *regexp.Regexp
var DecodePattern3 *regexp.Regexp

func init() {
	var err error
	DATAPattern, err = regexp.Compile(
		`var\s+DATA\s*=\s*'(.*?)'`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	NONCEPattern, err = regexp.Compile(
		`window\s*\[\s*"n.*?o.*?n.*?c.*?e"\s*\]\s*=\s*(.*?)\s*;`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}

	DecodePattern1, err = regexp.Compile(
		`(\d+)[a-zA-Z]+`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DecodePattern2, err = regexp.Compile(
		`\d+`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	DecodePattern3, err = regexp.Compile(
		`[^A-Za-z0-9\+\/\=]`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func downloadComic(cookies []*http.Cookie, url string) error {
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
	// cookies = resp.Cookies()
	var html string
	{
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		html = string(body)
	}
	// log.Println("html: ", html)
	// "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	// "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	// "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	matches := DATAPattern.FindStringSubmatch(html)
	var data string
	if len(matches) > 0 {
		data = matches[1]
	}
	matches2 := NONCEPattern.FindAllStringSubmatch(html, -1)
	var nonce string
	if len(matches) > 0 {
		nonces := matches2[len(matches)-1][1]
		jsFile, err := os.OpenFile(
			"command.js",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0777,
		)
		if err != nil {
			return err
		}
		jsFile.Write([]byte("console.log(" + nonces + ")"))
		jsFile.Close()
		output, err := exec.Command("node", "command.js").Output()
		if err != nil {
			return err
		}
		nonce = string(output)
	}
	// log.Println("data: ", data)
	// log.Println("nonce: ", nonce)
	info_str, err := decode(data, nonce)
	if err != nil {
		return err
	}
	log.Println(info_str)

	// req, err = http.NewRequest(
	// 	"POST",
	// 	"https://ac.qq.com/ComicView/getNextChapterPicture?id=623654&cid=1060",
	// 	nil,
	// )
	// if err != nil {
	// return err
	// }
	// for _, cookie := range cookies {
	// 	req.AddCookie(cookie)
	// }
	// resp, err = client.Do(req)
	// if err != nil {
	// return err
	// }

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// return err
	// }
	// log.Println(string(body))
	return nil
}

const _keyStr string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

func decode(data string, nonce string) (string, error) {
	var err error
	T := data
	N := nonce
	var locate int
	matches := DecodePattern1.FindAllStringSubmatch(N, -1)
	for length := len(matches) - 1; length >= 0; length-- {
		locate, err = strconv.Atoi(matches[length][1])
		if err != nil {
			return "", err
		}
		locate &= 255
		str := DecodePattern2.ReplaceAllString(matches[length][0], "")
		T = T[:locate] + T[(locate+len(str)):]
	}
	c := T
	a := ""
	b := 0
	d := 0
	h := 0
	f := 0
	g := 0
	e := 0
	c = DecodePattern3.ReplaceAllString(c, "")
	for e < len(c) {
		b = strings.Index(_keyStr, c[e:e+1])
		e++
		if e < len(c) {
			d = strings.Index(_keyStr, c[e:e+1])
			e++
		}
		if e < len(c) {
			f = strings.Index(_keyStr, c[e:e+1])
			e++
		}
		if e < len(c) {
			g = strings.Index(_keyStr, c[e:e+1])
			e++
		}
		b = b<<2 | d>>4
		d = (d&15)<<4 | f>>2
		h = (f&3)<<6 | g
		a += string(b)
		if f != 64 {
			a += string(d)
		}
		if g != 64 {
			a += string(h)
		}
	}
	// return tdecode(a)
	return utf8_decode([]byte(a)), nil
}

func utf8_decode(c []byte) string {
	a := ""
	b := 0
	c1 := 0
	d := 0
	for b < len(c) {
		d = int(c[b])
		if 128 > d {
			a += string(d)
			b++
		} else if 191 < d && 224 > d {
			c1 = int(c[b+1])
			a += string((d&31)<<6 | c1&63)
			b += 2

		} else {
			c1 = int(c[b+1])
			c2 := int(c[b+2])
			a += string((d&15)<<12 | (c1&63)<<6 | c2&63)
			b += 3
		}
	}
	return a
}
