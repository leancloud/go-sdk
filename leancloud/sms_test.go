package leancloud

import (
	"os"
	"testing"
)

func TestSMSCodeRequest(t *testing.T) {
	client = NewEnvClient()
	sms := &SMS{c}
	mobile := os.Getenv("TEST_SMS_REQUEST_MOBILE")
	err := sms.RequestSMSCode(mobile)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSMSCodeVerify(t *testing.T) {
	client = NewEnvClient()
	sms := &SMS{c}
	mobile := os.Getenv("TEST_SMS_REQUEST_MOBILE")
	code := os.Getenv("TEST_SMS_REQUEST_MOBILE_CODE")
	err := sms.VerifySMSCode(mobile, code)
	if err != nil {
		t.Error(err)
	}
}
