package lanzou

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

var (
	headers       map[string]string
	last_headers  map[string]string
	client        *req.Client
	url_cache     string
	default_error error
	pub           int
)

func init() {
	client = req.C()
	client.SetRedirectPolicy(req.NoRedirectPolicy())
	default_error = errors.New("Get url error")
	url_cache = ""
	pub = 0
}

func requestUrl(url string) *req.Response {
	resp, err := client.R().SetHeaders(headers).
		Get(url)
	if err != nil {
		// do nothing
		fmt.Println("Error: Request ", err)
	}

	return resp
}

func Lanzou(url string, pwd string, args ...string) (string, error) {
	if len(args) > 0 {
		pub = 1
	}
	res, err := getLanzouUrl(url, &pwd)

	if err != nil {
		res, err = getLanzouUrl(replaceMainUrl(url), &pwd)

		if err != nil {
			return "", err
		}
	}

	return res, nil
}

func getLanzouUrl(url string, pwd *string) (string, error) {
	if checkExpiredTime(url_cache) {
		return url_cache, nil
	}

	generateRandHeader()
	mainUrl := string(regexp.MustCompile(`(^http.*?\.com)`).Find([]byte(url)))

	if mainUrl == "" {
		return "", errors.New("地址有误。")
	}

	firstRes := requestUrl(url)
	firstResString := firstRes.String()

	var signData []string
	util_url := ""
	var err error
	postData := map[string]string{
		"action": "downprocess",
		"signs":  "?ctdf",
	}
	if strings.ContainsAny(firstResString, "输入密码") {
		signData, err = getSignDataWithPwd(firstResString)
		if *pwd == "" {
			if pub != 0 {
				fmt.Println("该文件需要分享密码，请输入密码:")
				fmt.Scanln(pwd)
			} else {
				return "", errors.New("需要分享密码。")
			}
		}
		postData["p"] = *pwd
	} else {
		signData, util_url, err = getSignData(firstResString, mainUrl)
	}

	if err != nil {
		return "", err
	}

	postData["sign"] = signData[1]

	ajax_json := map[string]interface{}{}
	res, err := client.R().SetHeaders(headers).SetHeader("Referer", mainUrl+util_url).SetResult(&ajax_json).
		SetFormData(postData).Post(mainUrl + "/ajaxm.php")

	if err != nil || res.IsError() {
		return "", default_error
	}

	if ajax_json["zt"].(float64) != 1 {
		return "", errors.New(ajax_json["inf"].(string))
	}
	direct_url := ajax_json["dom"].(string) + "/file/" + ajax_json["url"].(string)

	lastRes, err := client.R().SetHeaders(last_headers).Get(direct_url)

	if err != nil || lastRes.StatusCode != 302 {
		return direct_url, nil
	}

	redirect_url, ok := lastRes.Header["Location"]

	if !ok {
		url_cache = ""
		return direct_url, nil
	}

	url_cache = redirect_url[0]
	return redirect_url[0], nil
}

func getSignData(resString string, mainUrl string) ([]string, string, error) {
	var signData []string
	reg := regexp.MustCompile(`iframe.+?src=\"([^\"]{20,}?)\"`)

	iframeSrc := reg.FindStringSubmatch(resString)

	if len(iframeSrc) == 0 {
		return signData, "", default_error
	}

	secRes := requestUrl(mainUrl + iframeSrc[1])
	reg = regexp.MustCompile(`'sign':'(.*?)',`)
	signData = reg.FindStringSubmatch(secRes.String())

	if len(signData) == 0 {
		return signData, "", default_error
	}

	return signData, iframeSrc[1], nil
}

func getSignDataWithPwd(resString string) ([]string, error) {
	var signData []string
	reg := regexp.MustCompile(`data.*?sign=([^']{20,})&p='`)
	signData = reg.FindStringSubmatch(resString)
	if len(signData) == 0 {
		return signData, default_error
	}

	return signData, nil
}

func randIP() string {
	buf := make([]byte, 4)
	ip := rand.Uint32()

	binary.LittleEndian.PutUint32(buf, ip)
	return net.IP(buf).String()
}

func checkExpiredTime(urls string) bool {
	if urls == "" {
		return false
	}

	e := getExpiredTime(urls)
	a, err := strconv.Atoi(e)
	if err != nil {
		return false
	}

	if a-int(time.Now().Unix())-120 > 0 {
		return true
	}

	return false
}

func getExpiredTime(urls string) string {
	u, err := url.Parse(urls)

	if err != nil {
		return ""
	}

	res, ok := u.Query()["e"]

	if !ok {
		return ""
	}

	return res[0]
}

func generateRandHeader() {
	rand_ip := randIP()

	headers = map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.9",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36 Edg/104.0.1293.70",
		"X-FORWARDED-FOR": rand_ip,
		"CLIENT-ip":       rand_ip,
	}
	last_headers = map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
		"Cookie":                    "down_ip=1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36 Edg/104.0.1293.70",
		"Referer":                   "https://developer.lanzoug.com",
		"X-FORWARDED-FOR":           rand_ip,
		"CLIENT-ip":                 rand_ip,
	}
}

func replaceMainUrl(url string) string {
	re := regexp.MustCompile(`(^http.*?\.com)`)
	res := re.ReplaceAllString(url, "https://www.lanzoui.com")
	return res
}
