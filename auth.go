package ysgo

import (
	"fmt"
	"strconv"
	"time"
)

func (c *YSClient) generateAuthToken() AuthToken {
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	return AuthToken{
		Username:  c.userEndpoint,
		Timestamp: timestamp,
	}
}

func (c *YSClient) setCommonHeaders(req map[string]string) {
	token := c.generateAuthToken()
	req["Authorization"] = token.String()
	req["Content-Type"] = "application/x-www-form-urlencoded"
	req["DNT"] = "1"
	req["Pragma"] = "no-cache"
	req["Cache-Control"] = "no-cache"
	req["Accept"] = "application/json"
}

func (c *YSClient) buildDownloadURL(downloadToken, fileToken string) string {
	return fmt.Sprintf("http://ys-j.ysepan.com/wap/%s/%s/%s",
		c.userEndpoint, downloadToken, fileToken)
}
