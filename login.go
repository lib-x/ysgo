package ysgo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	loginPath              = "/nkj/csxx.aspx"
	logoutPath             = "/nkj/csxx.aspx"
	sessionPath            = "/nkj/csxx.aspx"
	loginCommand           = "yzglmm"
	logoutCommand          = "tcgl"
	sessionCommand         = "dq"
	spacePasswordCommand   = "yzdlmm"
	defaultDirectoryNumber = "1445856"
)

func (c *YSClient) InitSession() (*SessionResponse, error) {
	return c.InitSessionContext(context.Background())
}

func (c *YSClient) InitSessionContext(ctx context.Context) (*SessionResponse, error) {
	return c.initSessionContext(ctx, true)
}

func (c *YSClient) initSessionContext(ctx context.Context, allowAutoSpacePassword bool) (*SessionResponse, error) {
	req, err := c.newRequestWithContext(ctx, http.MethodPost, sessionPath+"?cz="+sessionCommand, url.Values{})
	if err != nil {
		return nil, err
	}

	var sessionResp SessionResponse
	if err := c.do(req, &sessionResp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("init session: %w", err)
	}

	if sessionResp.Token != "" {
		c.setAuthToken(sessionResp.Token)
	}

	if allowAutoSpacePassword && requiresSpacePassword(sessionResp.Message) && c.getSpacePassword() != "" {
		if err := c.VerifySpacePasswordContext(ctx, c.getSpacePassword()); err != nil {
			return nil, err
		}
		return c.initSessionContext(ctx, false)
	}

	if requiresSpacePassword(sessionResp.Message) {
		return &sessionResp, ErrSpacePasswordRequired
	}

	return &sessionResp, nil
}

func (c *YSClient) VerifySpacePassword(password string) error {
	return c.VerifySpacePasswordContext(context.Background(), password)
}

func (c *YSClient) VerifySpacePasswordContext(ctx context.Context, password string) error {
	form := url.Values{}
	form.Add("dlmm", password)

	req, err := c.newRequestWithContext(ctx, http.MethodPost, sessionPath+"?cz="+spacePasswordCommand, form)
	if err != nil {
		return err
	}

	resp, err := c.getHTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("verify space password: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("verify space password: read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("verify space password: %w", parseYSError(resp.StatusCode, body, http.StatusText(resp.StatusCode)))
	}
	value := strings.TrimSpace(string(body))
	if value == "" {
		return nil
	}
	if strings.HasPrefix(value, "ERR") || strings.Contains(value, "密码") {
		return parseYSError(httpStatusAlreadyReported, body, "verify space password failed")
	}
	c.setAuthToken(value)
	return nil
}

func (c *YSClient) Login() (*LoginResponse, error) {
	return c.LoginContext(context.Background())
}

func (c *YSClient) LoginContext(ctx context.Context) (*LoginResponse, error) {
	form := url.Values{}
	form.Add("glmm", c.managementPass)
	form.Add("mlbh", c.getManagementDirectory())

	req, err := c.newRequestWithContext(ctx, http.MethodPost, loginPath+"?cz="+loginCommand, form)
	if err != nil {
		return nil, err
	}

	var loginResp LoginResponse
	if err := c.do(req, &loginResp, http.StatusOK); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	return &loginResp, nil
}

func (c *YSClient) Logout() error {
	return c.LogoutContext(context.Background())
}

func (c *YSClient) LogoutContext(ctx context.Context) error {
	req, err := c.newRequestWithContext(ctx, http.MethodGet, logoutPath+"?cz="+logoutCommand, nil)
	if err != nil {
		return err
	}

	if err := c.do(req, nil, http.StatusOK); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	return nil
}

func requiresSpacePassword(message string) bool {
	message = strings.TrimSpace(message)
	return strings.Contains(message, "需要输入登陆密码") || strings.Contains(message, "限制访客登陆")
}
