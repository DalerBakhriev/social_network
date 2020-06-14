package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("Record not found")

	// ErrFriendRequestWasAlreadySent ...
	ErrFriendRequestWasAlreadySent = errors.New("Friend request was already sent")
)
