package ysgo

import "context"

// PrepareSession initializes the ysepan session.
func (c *YSClient) PrepareSession() (*SessionResponse, error) {
	return c.PrepareSessionContext(context.Background())
}

// PrepareSessionContext initializes the ysepan session with context.
func (c *YSClient) PrepareSessionContext(ctx context.Context) (*SessionResponse, error) {
	return c.InitSessionContext(ctx)
}

// PrepareAdminSession initializes the session and performs admin login.
func (c *YSClient) PrepareAdminSession() (*LoginResponse, error) {
	return c.PrepareAdminSessionContext(context.Background())
}

// PrepareAdminSessionContext initializes the session and performs admin login with context.
func (c *YSClient) PrepareAdminSessionContext(ctx context.Context) (*LoginResponse, error) {
	if _, err := c.InitSessionContext(ctx); err != nil {
		return nil, err
	}
	return c.LoginContext(ctx)
}
