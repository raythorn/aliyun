package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	AliSMSUrl         = "http://sms.aliyuncs.com"
	AliSMSAction      = "SingleSendSms"
	AliSMSSignMethod  = "HMAC-SHA1"
	AliSMSSignVersion = "1.0"
	AliSMSVersion     = "2016-09-27"
)

const (
	SMS_RAND_KIND_DIGIT = iota
	SMS_RAND_KIND_LOWER
	SMS_RAND_KIND_UPPER
	SMS_RAND_KIND_ALL
)

type Alisms struct {
	Phone           string
	SignName        string
	TemplateCode    string
	Params          string
	Format          string
	accessKey       string
	accessKeySecret string
}

func (sms *Alisms) SendSMS(key, secret string) error {
	if key == "" ||
		secret == "" ||
		sms.Phone == "" ||
		sms.SignName == "" ||
		sms.TemplateCode == "" {
		return errors.New("Invalid parameter!")
	}

	if sms.Format == "" {
		sms.Format = "JSON"
	}

	sms.accessKey = key
	sms.accessKeySecret = secret

	param := sms.parameter()

	paramstr := param.Encode()
	signature := sms.signature(paramstr)

	params, err := url.QueryUnescape(paramstr)
	if err != nil {
		return err
	}

	log.Println(params)

	body := ioutil.NopCloser(strings.NewReader("Signature=" + signature + "&" + params))
	log.Println(body)

	client := &http.Client{}
	request, err := http.NewRequest("POST", AliSMSUrl, body)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var resp map[string]string
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	log.Println(resp)
	if msg, exist := resp["Message"]; exist {
		return errors.New(msg)
	}

	return nil
}

func (sms *Alisms) Random(size, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	if kind < 0 || kind > SMS_RAND_KIND_ALL {
		ikind = 0
	}

	all := ikind == SMS_RAND_KIND_ALL

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if all {
			ikind = rand.Intn(SMS_RAND_KIND_ALL)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = byte(base + rand.Intn(scope))
	}

	return string(result)
}

func (sms *Alisms) signature(param string) string {

	param = "POST&%2F&" + url.QueryEscape(param)

	secret := sms.accessKeySecret + "&"
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(param))

	signstr := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(signstr)
}

func (sms *Alisms) parameter() *url.Values {

	param := &url.Values{}

	param.Add("Action", AliSMSAction)
	param.Add("SignName", sms.SignName)
	param.Add("TemplateCode", sms.TemplateCode)
	param.Add("RecNum", sms.Phone)
	param.Add("ParamString", sms.Params)
	param.Add("Format", sms.Format)
	param.Add("Version", AliSMSVersion)
	param.Add("AccessKeyId", sms.accessKey)
	param.Add("SignatureMethod", AliSMSSignMethod)
	param.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	param.Add("SignatureVersion", AliSMSSignVersion)
	param.Add("SignatureNonce", sms.nonce())

	return param
}

func (sms *Alisms) nonce() string {
	nonce := sms.Random(8, SMS_RAND_KIND_ALL) + "-" + sms.Random(4, SMS_RAND_KIND_ALL) + "-" + sms.Random(4, SMS_RAND_KIND_ALL) + "-" + sms.Random(4, SMS_RAND_KIND_ALL) + "-" + sms.Random(12, SMS_RAND_KIND_ALL)
	return strings.ToLower(nonce)
}
