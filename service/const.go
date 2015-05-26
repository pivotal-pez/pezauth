package pezauth

import "errors"

var (
	ErrCouldNotGetUserGUID = errors.New("query failed. unable to find matching user guid.")
)
