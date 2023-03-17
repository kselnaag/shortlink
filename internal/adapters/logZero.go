package adapters

import "shortlink/internal/ports"

var _ ports.Ilog = (*LogZero)(nil)

type LogZero struct {
	//
}

func NewLogZero() LogZero {
	//
	return LogZero{}
}
