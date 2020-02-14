package jk_err

import "errors"

type JKErrorStatus int

const (
	Input JKErrorStatus = iota
	Channel
	Network
)

func JKErrInput(status JKErrorStatus) error  {

	switch status {
	case Input:
		return errors.New("Input Error")
	case Channel:
		return errors.New("Channel Error")
	default:
		return errors.New("Unknow Error")
	}
}