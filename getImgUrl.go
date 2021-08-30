package main

import (
	"fmt"
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

func getImgUrl(url string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return "", err
	}
	// for _, cookie := range cookies {
	// 	// if len(cookie.Value) > 0 && cookie.Value != "/" && cookie.Value != "/;" && cookie.Value != ";" {
	// 	req.AddCookie(cookie)
	// 	// }
	// }
	// log.Println(cookies)
	// fmt.Println()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	// cookies = resp.Cookies()
	var html string
	{
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		html = string(body)
	}
	matches := DATAPattern.FindStringSubmatch(html)
	var data string
	if len(matches) > 0 {
		data = matches[1]
	}
	matches2 := NONCEPattern.FindAllStringSubmatch(html, -1)
	var nonce string
	if len(matches) > 0 {
		nonces := matches2[len(matches)-1][1]
		nonces = strings.Replace(nonces, "!document.children", "false", -1)
		nonces = strings.Replace(nonces, "!window.Array", "false", -1)
		nonces = strings.Replace(nonces, "!document.getElementsByTagName('html')", "false", -1)
		jsFile, err := os.OpenFile(
			"command.js",
			os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
			0777,
		)
		if err != nil {
			return "", err
		}
		jsFile.Write([]byte("console.log(" + nonces + ")"))
		jsFile.Close()
		output, err := exec.Command("node", "command.js").Output()
		if err != nil {
			log.Println("nonces: ", nonces)
			fmt.Println()
			return "", err
		}
		nonce = string(output)
	}
	// log.Println("data: ", data)
	// log.Println("nonce: ", nonce)
	info_str, err := decode(data, nonce)
	if err != nil {
		return "", err
	}
	// log.Println(info_str)
	// fmt.Println()
	return info_str, err
}

const _keyStr string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="

func decode(data string, nonce string) (string, error) {
	var err error
	T := data
	N := nonce
	var locate uint64
	matches := DecodePattern1.FindAllStringSubmatch(N, -1)
	for length := len(matches) - 1; length >= 0; length-- {
		locate, err = strconv.ParseUint(matches[length][1], 10, 64)
		if err != nil {
			return "", err
		}
		locate &= 255
		str := DecodePattern2.ReplaceAllString(matches[length][0], "")
		T = T[:locate] + T[(int(locate)+len(str)):]
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
		d = strings.Index(_keyStr, c[e:e+1])
		e++
		f = strings.Index(_keyStr, c[e:e+1])
		e++
		g = strings.Index(_keyStr, c[e:e+1])
		e++
		b = b<<2 | d>>4
		d = (d&15)<<4 | f>>2
		h = (f&3)<<6 | g
		a += string(rune(b)) // 直接string(b)也完全一样，但有警告看着很烦
		if f != 64 {
			a += string(rune(d))
		}
		if g != 64 {
			a += string(rune(h))
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
			a += string(rune(d))
			b++
		} else if 191 < d && 224 > d {
			c1 = int(c[b+1])
			a += string(rune((d&31)<<6 | c1&63))
			b += 2

		} else {
			c1 = int(c[b+1])
			c2 := int(c[b+2])
			a += string(rune((d&15)<<12 | (c1&63)<<6 | c2&63))
			b += 3
		}
	}
	return a
}
