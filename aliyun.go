package aliyun

import (
	"errors"
	"github.com/raythorn/aliyun/sms"
)

var (
	ali *aliyun
)

func init() {
	ali = &aliyun{accessKey: "", accessKeySecret: ""}
}

type aliyun struct {
	accessKey       string
	accessKeySecret string
}

func (a *aliyun) init(accessKey, accessKeySecret string) {
	a.accessKey = accessKey
	a.accessKeySecret = accessKeySecret
}

func (a *aliyun) sendsms(alisms *sms.Alisms) error {
	return alisms.SendSMS(a.accessKey, a.accessKeySecret)
}

func Init(accessKey, accessKeySecret string) {
	ali.init(accessKey, accessKeySecret)
}

func SendSMS(alisms *sms.Alisms) error {
	if alisms == nil {
		return errors.New("Invalid parameter!")
	}

	return ali.sendsms(alisms)
}
