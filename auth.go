package ysgo

import "fmt"

func (c *YSClient) generateAuthToken() AuthToken {
	return AuthToken{
		Username: c.GetUserEndpoint(),
		Token:    c.GetAuthToken(),
	}
}

func (c *YSClient) buildDownloadURL(downloadToken, fileToken string) string {
	return fmt.Sprintf("https://ys-j.ysepan.com/wap/%s/%s/%s",
		c.GetUserEndpoint(), downloadToken, fileToken)
}
