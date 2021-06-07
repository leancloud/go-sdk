package leancloud

import "fmt"

type SMS struct {
	c *Client
}

func (ref *SMS) RequestSMSCode(number string) error {
	path := "/1.1/requestSmsCode"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
	}

	_, err := ref.c.request(methodPost, path, options)
	if err != nil {
		return err
	}

	return nil
}

func (ref *SMS) VerifySMSCode(number, smsCode string) error {
	path := fmt.Sprintf("/1.1/verifySmsCode/%s", smsCode)
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
	}

	_, err := ref.c.request(methodPost, path, options)
	if err != nil {
		return err
	}

	return nil
}
