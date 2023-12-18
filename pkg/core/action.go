package core

import (
	"errors"
	"strings"
)

type Action int

//go:generate stringer -type=Action
const (
	Unknown Action = iota
	Allow
	Deny
)

var ErrUnsupported = errors.New("unsupported action type")

var mapEnumStringToAction = func() map[string]Action {
	m := make(map[string]Action)
	for i := Allow; i <= Deny; i++ {
		m[strings.ToLower(i.String())] = i
	}
	return m
}()

func (t *Action) UnmarshalText(text []byte) error {
	s := string(text)
	if val, ok := mapEnumStringToAction[strings.ToLower(s)]; ok {
		*t = val
		return nil
	}

	return ErrUnsupported
}
