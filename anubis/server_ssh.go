package anubis

import (
    "Anubis/pkg/config"
    "Anubis/pkg/handler"
    "Anubis/pkg/logger"
    "Anubis/pkg/sshd"
    "github.com/gliderlabs/ssh"
    gossh "golang.org/x/crypto/ssh"
    "io/ioutil"
    "net"
)

type server struct {
}

func (s *server) GetSSHAddr() string {
    cf := config.GlobalConfig
    return net.JoinHostPort(cf.BindHost, cf.SSHPort)
}

func (s *server) GetSSHSigner() ssh.Signer {
    cf := config.GlobalConfig
    privateBytes, err := ioutil.ReadFile(cf.HostKeyFile)
    if err != nil {
        logger.Fatalf("Read host key failed: %s", err)
    }

    signer, err := gossh.ParsePrivateKey(privateBytes)
    if err != nil {
        logger.Fatalf("Parse private key failed: %s\n", err)
    }

    return signer
}

func (s *server) KeyboardInteractiveAuth(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
    return sshd.AuthFailed
}

func (s *server) PasswordAuth(ctx ssh.Context, password string) bool {
    username := ctx.User()
    if username == "root" && password == "root" {
        return sshd.AuthSuccessful
    }

    return sshd.AuthFailed
}

func (s *server) PublicKeyAuth(ctx ssh.Context, key ssh.PublicKey) bool {
    return sshd.AuthFailed
}

func (s *server) NextAuthMethodsHandler(ctx ssh.Context) []string {
    return []string{}
}

func (s *server) SessionHandler(session ssh.Session) {

    interactiveSrv := handler.NewInteractiveHandler(session)
    logger.Infof("User %s request pty", session.User())
    interactiveSrv.Dispatch()
    return
}
