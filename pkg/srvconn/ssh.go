package srvconn

import (
	"Anubis/pkg/logger"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type SSHClientOption func(conf *SSHClientOptions)

type SSHClientOptions struct {
	Host       string
	Port       string
	Username   string
	Password   string
	PrivateKey string
	//Passphrase   string
	Timeout int
	//keyboardAuth gossh.KeyboardInteractiveChallenge
	PrivateAuth gossh.Signer

	//proxySSHClientOptions []SSHClientOptions
}

type SSHClient struct {
	*gossh.Client
	Cfg *SSHClientOptions
	//ProxyClient *SSHClient

	sync.Mutex

	traceSessionMap map[*gossh.Session]time.Time

	refCount int32
}

func NewSSHClient(opts ...SSHClientOption) (*SSHClient, error) {
	cfg := &SSHClientOptions{
		Host: "127.0.0.1",
		Port: "22",
	}
	for _, setter := range opts {
		setter(cfg)
	}

	return NewSSHClientWithCfg(cfg)
}

func NewSSHClientWithCfg(cfg *SSHClientOptions) (*SSHClient, error) {
	gosshCfg := gossh.ClientConfig{
		User:            cfg.Username,
		Auth:            cfg.AuthMethods(),
		Timeout:         time.Duration(cfg.Timeout) * time.Second,
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Config:          createSSHConfig(),
	}

	destAddr := net.JoinHostPort(cfg.Host, cfg.Port)
	gosshClient, err := gossh.Dial("tcp", destAddr, &gosshCfg)
	if err != nil {
		return nil, err
	}

	return &SSHClient{Client: gosshClient, Cfg: cfg,
		traceSessionMap: make(map[*gossh.Session]time.Time)}, nil
}

func (cfg *SSHClientOptions) AuthMethods() []gossh.AuthMethod {
	authMethods := make([]gossh.AuthMethod, 0, 1)
	if cfg.Password != "" {
		authMethods = append(authMethods, gossh.Password(cfg.Password))
	}

	return authMethods
}

func createSSHConfig() gossh.Config {
	var cfg gossh.Config
	cfg.SetDefaults()
	cfg.Ciphers = append(cfg.Ciphers, notRecommendCiphers...)
	cfg.KeyExchanges = append(cfg.KeyExchanges, notRecommendKeyExchanges...)
	return cfg
}

var (
	notRecommendCiphers = []string{
		"arcfour256", "arcfour128", "arcfour",
		"aes128-cbc", "3des-cbc",
	}

	notRecommendKeyExchanges = []string{
		"diffie-hellman-group1-sha1", "diffie-hellman-group-exchange-sha1",
		"diffie-hellman-group-exchange-sha256",
	}
)

func (s *SSHClient) AcquireSession() (*gossh.Session, error) {
	atomic.AddInt32(&s.refCount, 1)
	sess, err := s.Client.NewSession()
	if err != nil {
		atomic.AddInt32(&s.refCount, -1)
		return nil, err
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.traceSessionMap[sess] = time.Now()
	logger.Infof("SSHClient(%s) session add one ", s)
	return sess, nil
}

func (s *SSHClient) ReleaseSession(sess *gossh.Session) {
	atomic.AddInt32(&s.refCount, -1)
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.traceSessionMap, sess)
	logger.Infof("SSHClient(%s) release one session remain %d", s, len(s.traceSessionMap))
}

func SSHClientUsername(username string) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.Username = username
	}
}

func SSHClientPassword(password string) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.Password = password
	}
}

func SSHClientPrivateKey(privateKey string) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.PrivateKey = privateKey
	}
}

//func SSHClientPassphrase(passphrase string) SSHClientOption {
//    return func(args *SSHClientOptions) {
//        args.Passphrase = passphrase
//    }
//}

func SSHClientHost(host string) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.Host = host
	}
}

func SSHClientPort(port int) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.Port = strconv.Itoa(port)
	}
}

func SSHClientTimeout(timeout int) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.Timeout = timeout
	}
}

func SSHClientPrivateAuth(privateAuth gossh.Signer) SSHClientOption {
	return func(args *SSHClientOptions) {
		args.PrivateAuth = privateAuth
	}
}

//func SSHClientProxyClient(proxyArgs ...SSHClientOptions) SSHClientOption {
//    return func(args *SSHClientOptions) {
//        args.proxySSHClientOptions = proxyArgs
//    }
//}
//
//func SSHClientKeyboardAuth(keyboardAuth gossh.KeyboardInteractiveChallenge) SSHClientOption {
//    return func(conf *SSHClientOptions) {
//        conf.keyboardAuth = keyboardAuth
//    }
//}
