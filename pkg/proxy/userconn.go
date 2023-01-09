package proxy

import (
	"context"
	"github.com/gliderlabs/ssh"
	"io"
)

type UserConnection interface {
	io.ReadWriteCloser
	ID() string
	WinCh() <-chan ssh.Window
	LoginFrom() string
	RemoteAddr() string
	Pty() ssh.Pty
	Context() context.Context
}
