package commands

import "errors"

var (
	ErrEnvFileNotFound     = errors.New(".env file not found")
	ErrInternalServerError = errors.New("internal server error")
	ErrLocationNotFound    = errors.New("location not found")
)
