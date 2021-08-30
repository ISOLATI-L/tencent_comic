package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var UrlPattern *regexp.Regexp
var ClientIdPattern *regexp.Regexp
var LoginTypePattern *regexp.Regexp
var UinPattern *regexp.Regexp
var CodePattern *regexp.Regexp

func init() {
	var err error
	UrlPattern, err = regexp.Compile(
		`'(https.*?)'`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	ClientIdPattern, err = regexp.Compile(
		`pt_3rd_aid=(\d+)'`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	LoginTypePattern, err = regexp.Compile(
		`pt_login_type=(\d+)`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	UinPattern, err = regexp.Compile(
		`uin=(\d+)`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
	CodePattern, err = regexp.Compile(
		`code=(.+?)$`,
	)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func Get_login_sig() (pt_login_sig string, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://xui.ptlogin2.qq.com/cgi-bin/xlogin?appid=716027609&login_text=授权并登录&hide_title_bar=1&hide_border=1&target=self&s_url=https://graph.qq.com/oauth2.0/login_jump&pt_3rd_aid=101483258&pt_feedback_link=https://support.qq.com/products/77942?customInfo=.appid101483258",
		nil,
	)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	found := false
	cookies := getSetCookies(resp)
	for _, cookie := range cookies {
		if cookie.Name == "pt_login_sig" {
			found = true
			pt_login_sig = cookie.Value
			break
		}
	}
	if !found {
		return "", errors.New("can not find pt_login_sig")
	}
	return pt_login_sig, nil
}

const QR_SIZE int = 37

func ShowQRcode() (ptqrtoken uint64, err error) {
	randFloat := rand.Float64()
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"https://ssl.ptlogin2.qq.com/ptqrshow?appid=716027609&e=2&l=M&s=3&d=72&v=4&t=%f&daid=383&pt_3rd_aid=101483258",
			randFloat,
		),
		nil,
	)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	QRcode, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	reader := bytes.NewReader(QRcode)
	image, err := png.Decode(reader)
	if err != nil {
		return 0, err
	}
	binary := binaryImage(image, QR_SIZE, QR_SIZE)
	for y := 0; y < QR_SIZE; y++ {
		for x := 0; x < QR_SIZE; x++ {
			if binary[x][y] {
				textbackground(0xFF)
				fmt.Print("　")
				// fmt.Print("■")
			} else {
				textbackground(0x00)
				fmt.Print("　")
			}
		}
		resettextbackground()
		fmt.Println()
	}
	resettextbackground()

	found := false
	cookies := getSetCookies(resp)
	var qrsig string
	for _, cookie := range cookies {
		if cookie.Name == "qrsig" {
			found = true
			qrsig = cookie.Value
			break
		}
	}
	if !found {
		return 0, errors.New("can not find qrsig")
	}
	for _, c := range []byte(qrsig) {
		ptqrtoken += ptqrtoken<<5 + uint64(c)
	}
	ptqrtoken &= 0x7FFFFFFF

	return ptqrtoken, nil
}

func binaryImage(m image.Image, w int, h int) [][]bool {
	bounds := m.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	var dx, dy int
	if w > 0 {
		dx = width / w
	} else {
		dx = 1
	}
	if w > 0 {
		dy = height / h
	} else {
		dy = 1
	}
	result := make([][]bool, w)
	var wg sync.WaitGroup
	wg.Add(w)
	for x := 0; x < w; x++ {
		result[x] = make([]bool, h)
		go func(x int) {
			for y := 0; y < h; y++ {
				_, g, _, _ := m.At(x*dx+(dx-1)/2, y*dy+(dy-1)/2).RGBA()
				result[x][y] = g>>15 > 0
			}
			wg.Done()
		}(x)
	}
	wg.Wait()
	return result
}

func GetAction() (action string) {
	action = fmt.Sprintf("0-0-%d", time.Now().UnixNano()/int64(time.Millisecond))
	return action
}

func Login() (err error) {
	for {
		login_sig, err := Get_login_sig()
		if err != nil {
			return err
		}
		for {
			ptqrtoken, err := ShowQRcode()
			if err != nil {
				break
			}
			action := GetAction()
			// log.Println("pt_login_sig: ", login_sig)
			// log.Println("ptqrtoken: ", ptqrtoken)
			// log.Println("action: ", action)
			textbackground(0x04)
			fmt.Println("请扫码登陆")
			resettextbackground()
			var req *http.Request
			req, err = http.NewRequest(
				"GET",
				fmt.Sprintf(
					"https://ssl.ptlogin2.qq.com/ptqrlogin?u1=https://graph.qq.com/oauth2.0/login_jump&ptqrtoken=%d&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=%s&js_ver=21082415&js_type=1&login_sig=%s&pt_uistyle=40&aid=716027609&daid=383&pt_3rd_aid=101483258&",
					ptqrtoken,
					action,
					login_sig,
				),
				nil,
			)
			if err != nil {
				break
			}
			ch := make(chan error)
			go func() {
				for {
					resp, err := client.Do(req)
					if err != nil {
						ch <- err
						return
					}
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						ch <- err
						return
					}
					html := string(body)
					// log.Println(html)
					if strings.Contains(html, "登录成功") {
						matches := LoginTypePattern.FindStringSubmatch(html)
						if len(matches) == 0 {
							ch <- errors.New("can not find pt_login_type")
						}
						err = getUserInfo(html)
						ch <- err
						return
					} else if strings.Contains(html, "二维码已经失效") {
						textbackground(0x04)
						fmt.Println("二维码已经失效")
						resettextbackground()
						ch <- errors.New("time out")
						return
					}
					time.Sleep(1 * time.Second)
				}
			}()
			err = <-ch
			if err == nil {
				break
			}
		}
		if err == nil {
			break
		}
	}
	return nil
}

func getUserInfo(loginedHTML string) error {
	matches := UrlPattern.FindStringSubmatch(loginedHTML)
	if len(matches) == 0 {
		return errors.New("can not find url")
	}
	url := matches[1]
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	cookies := getSetCookies(resp)
	var p_skey string
	for _, cookie := range cookies {
		if cookie.Name == "p_skey" && len(cookie.Value) > 0 {
			p_skey = cookie.Value
			break
		}
	}

	matches = ClientIdPattern.FindStringSubmatch(loginedHTML)
	if len(matches) == 0 {
		return errors.New("can not find client id")
	}
	clientId := matches[1]
	g_tk := getGTK(p_skey)
	reqBody := fmt.Sprintf(
		"response_type=code&client_id=%s&redirect_uri=https://ac.qq.com/loginSuccess.html?url=https://ac.qq.com/Comic/comicInfo/id/%s?auth=1&scope=&state=&switch=&from_ptlogin=1&src=1&update_auth=1&openapi=80901010&g_tk=%d&auth_time=%d&ui=C783C362-54B4-4154-9DDF-4BB1757CAE80",
		clientId,
		CFG.id,
		g_tk,
		time.Now().UnixNano()/int64(time.Millisecond),
	)
	log.Println("p_skey: ", p_skey)
	fmt.Println()
	log.Println("reqBody: ", reqBody)
	fmt.Println()
	req, err = http.NewRequest(
		"POST",
		"https://graph.qq.com/oauth2.0/authorize",
		bytes.NewReader([]byte(reqBody)),
	)
	if err != nil {
		return err
	}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	location := resp.Header.Get("Location")
	log.Println("location: ", location)
	fmt.Println()
	matches = CodePattern.FindStringSubmatch(location)
	if len(matches) == 0 {
		return errors.New("can not find code")
	}
	code := matches[1]
	log.Println("code: ", code)
	fmt.Println()

	reqBody = "code=" + code
	rannum := rand.Float64()
	log.Println("url: ", fmt.Sprintf(
		"https://ac.qq.com/User/qqInfo?%f",
		rannum,
	))
	fmt.Println()
	req, err = http.NewRequest(
		"POST",
		fmt.Sprintf(
			"https://ac.qq.com/User/qqInfo?%f",
			rannum,
		),
		bytes.NewReader([]byte(reqBody)),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	html := string(body)
	log.Println(html)
	fmt.Println()
	return nil
}

func getGTK(skey string) int {
	var hash = 5381
	length := len(skey)
	for i := 0; i < length; i++ {
		hash += (hash << 5) + int(byte(skey[i]))
	}
	return hash & 0x7fffffff
}
