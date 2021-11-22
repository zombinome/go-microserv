package hosting

import "github.com/zombinome/go-microserv/microserv"

type HttpHostConfiguration struct {
	Address string
	Logger  microserv.Logger
}
