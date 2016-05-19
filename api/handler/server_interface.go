package handler

import (
	"gopush/lib"
)

type Server interface {
	GetPool() *lib.Pool

	GetEnv() lib.EnvInfo
}