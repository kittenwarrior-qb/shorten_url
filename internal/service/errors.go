package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrLinkNotFound       = errors.New("link not found")
	ErrInvalidURL         = errors.New("invalid URL")
	ErrInvalidAlias       = errors.New("invalid alias")
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrLinkExpired        = errors.New("link has expired")
	ErrInvalidToken       = errors.New("invalid token")
)
