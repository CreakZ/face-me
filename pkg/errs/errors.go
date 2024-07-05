package custom_errors

import (
	"errors"
)

var (
	ErrNoMatches     = errors.New("no matches played")
	ErrNotAllMatches = errors.New("info about certain matches failed to fetch due to external reasons")
	ErrWrongNickname = errors.New("wrong nickname")
)

func NoMatches() error {
	return ErrNoMatches
}

// maybe no need in this function
func NotAllMatchesFetched() error {
	return ErrNotAllMatches
}

func WrongNickname() error {
	return ErrWrongNickname
}
