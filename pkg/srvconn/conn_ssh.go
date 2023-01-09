package srvconn

import (
	"errors"
	gossh "golang.org/x/crypto/ssh"
	"io"
)

type SSHOptions struct {
	charset string
	win     Windows
	term    string

	isLoginToSu  bool
	sudoCommand  string
	sudoUsername string
	sudoPassword string
}

type SSHConnection struct {
	session *gossh.Session
	stdin   io.Writer
	stdout  io.Reader
	options *SSHOptions
}

func NewSSHConnection(sess *gossh.Session, opts ...SSHOption) (*SSHConnection, error) {
	if sess == nil {
		return nil, errors.New("ssh session is nil")
	}
	options := &SSHOptions{
		charset: "utf8",
		win: Windows{
			Width:  80,
			Height: 120,
		},
		term: "xterm",
	}
	for _, setter := range opts {
		setter(options)
	}
	modes := gossh.TerminalModes{
		gossh.ECHO:          1,     // enable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4 kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4 kbaud
	}

	err := sess.RequestPty(options.term, options.win.Height, options.win.Width, modes)
	if err != nil {
		return nil, err
	}
	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := sess.StdoutPipe()
	if err != nil {
		return nil, err
	}
	conn := &SSHConnection{
		session: sess,
		stdin:   stdin,
		stdout:  stdout,
		options: options,
	}

	err = sess.Shell()
	if err != nil {
		_ = sess.Close()
		return nil, err
	}
	return conn, nil
}

func (sc *SSHConnection) SetWinSize(w, h int) error {
	return sc.session.WindowChange(h, w)
}

func (sc *SSHConnection) Read(p []byte) (n int, err error) {
	return sc.stdout.Read(p)
}

func (sc *SSHConnection) Write(p []byte) (n int, err error) {
	return sc.stdin.Write(p)
}

func (sc *SSHConnection) Close() (err error) {
	return sc.session.Close()
}

func (sc *SSHConnection) KeepAlive() error {
	_, err := sc.session.SendRequest("keepalive@openssh.com", false, nil)
	return err
}

type SSHOption func(*SSHOptions)

func SSHCharset(charset string) SSHOption {
	return func(opt *SSHOptions) {
		opt.charset = charset
	}
}

func SSHPtyWin(win Windows) SSHOption {
	return func(opt *SSHOptions) {
		opt.win = win
	}
}

func SSHTerm(termType string) SSHOption {
	return func(opt *SSHOptions) {
		opt.term = termType
	}
}

func SSHLoginToSudo(ok bool) SSHOption {
	return func(opt *SSHOptions) {
		opt.isLoginToSu = ok
	}
}

func SSHSudoCommand(cmd string) SSHOption {
	return func(opt *SSHOptions) {
		opt.sudoCommand = cmd
	}
}

func SSHSudoUsername(username string) SSHOption {
	return func(opt *SSHOptions) {
		opt.sudoUsername = username
	}
}

func SSHSudoPassword(password string) SSHOption {
	return func(opt *SSHOptions) {
		opt.sudoPassword = password
	}
}
