package client_notify

import (
	"context"
	"errors"
)

func ensure() error {
	if Default == nil {
		return errors.New("client_notify: not initialized; call Init(url) first")
	}
	return nil
}

func GetStatus(ctx context.Context) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.GetStatus(ctx)
}

func SendWelcomeEmail(ctx context.Context, to, userName string) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.SendWelcomeEmail(ctx, to, userName)
}

func SendSecurityCodeEmail(ctx context.Context, to, code string) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.SendSecurityCodeEmail(ctx, to, code)
}

func SendServiceAlertEmail(ctx context.Context, to, title, message string) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.SendServiceAlertEmail(ctx, to, title, message)
}

func SendLoginLinkEmail(ctx context.Context, to, link, name string) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.SendLoginLinkEmail(ctx, to, link, name)
}

func SendGenericEmail(ctx context.Context, to, subject, message string) (*EmailResponse, error) {
	if err := ensure(); err != nil {
		return nil, err
	}
	return Default.SendGenericEmail(ctx, to, subject, message)
}
