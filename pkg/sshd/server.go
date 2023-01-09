package sshd

import (
    "Anubis/pkg/config"
    "Anubis/pkg/logger"
    "context"
    "github.com/gliderlabs/ssh"
    gossh "golang.org/x/crypto/ssh"
    "net"
    "time"
)

type AuthResult bool

const (
    AuthFailed     = false
    AuthSuccessful = true
)

type Server struct {
    Srv *ssh.Server
}

type SSHHandler interface {
    GetSSHAddr() string
    GetSSHSigner() ssh.Signer
    KeyboardInteractiveAuth(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool
    PasswordAuth(ctx ssh.Context, password string) bool
    PublicKeyAuth(ctx ssh.Context, key ssh.PublicKey) bool
    NextAuthMethodsHandler(ctx ssh.Context) []string
    SessionHandler(session ssh.Session)
}

func NewSSHServer(handler SSHHandler) *Server {
    srv := &ssh.Server{
        Addr:        handler.GetSSHAddr(),
        HostSigners: []ssh.Signer{handler.GetSSHSigner()},
        Version:     config.GlobalConfig.SSHVersion,
        KeyboardInteractiveHandler: func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
            return handler.KeyboardInteractiveAuth(ctx, challenger)
        },
        PasswordHandler: func(ctx ssh.Context, password string) bool {
            return handler.PasswordAuth(ctx, password)
        },
        PublicKeyHandler: func(ctx ssh.Context, key ssh.PublicKey) bool {
            return handler.PublicKeyAuth(ctx, key)
        },
        Handler: handler.SessionHandler,
    }

    return &Server{srv}
}

func (s *Server) Start() {
    logger.Infof("Start SSH server at %s", s.Srv.Addr)
    listen, err := net.Listen("tcp", s.Srv.Addr)
    if err != nil {
        logger.Fatal(err)
    }

    logger.Fatal(s.Srv.Serve(listen))
}

func (s *Server) Stop() {
    logger.Info("Stop SSH server")

    ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelFunc()
    logger.Fatal(s.Srv.Shutdown(ctx))
}
