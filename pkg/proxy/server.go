package proxy

import (
	"Anubis/pkg/common"
	"Anubis/pkg/logger"
	"Anubis/pkg/model"
	"Anubis/pkg/srvconn"
	"context"
	"errors"
)

var (
	ErrMissClient      = errors.New("the protocol client has not installed")
	ErrUnMatchProtocol = errors.New("the protocols are not matched")
	ErrAPIFailed       = errors.New("api failed")
	ErrPermission      = errors.New("no permission")
	ErrNoAuthInfo      = errors.New("no auth info")
)

type Server struct {
	ID          string
	UserConn    UserConnection
	connOpts    *ConnectionOptions
	account     *model.Account
	sessionInfo *model.Session
}

func NewServer(conn UserConnection, opts ...ConnectionOption) (*Server, error) {
	connOpts := &ConnectionOptions{}
	for _, setter := range opts {
		setter(connOpts)
	}

	account := connOpts.predefinedAccount
	assetName := "ssh"

	apiSession := &model.Session{
		ID:         common.UUID(),
		User:       connOpts.user,
		Account:    account.String(),
		LoginForm:  conn.LoginFrom(),
		RemoteAddr: conn.RemoteAddr(),
		Protocol:   connOpts.Protocol,
		Asset:      assetName,
	}

	return &Server{
		ID:          apiSession.ID,
		UserConn:    conn,
		connOpts:    connOpts,
		account:     account,
		sessionInfo: apiSession,
	}, nil
}

func (s *Server) Proxy() {
	ctx, cancel := context.WithCancel(context.Background())
	sw := SwitchSession{
		ID:            s.ID,
		MaxIdleTime:   60,
		keepAliveTime: 60,
		ctx:           ctx,
		cancel:        cancel,
		p:             s,
	}

	srvConn, err := s.getServerConn()
	if err != nil {
		logger.Error(err)
		return
	}
	defer srvConn.Close()

	logger.Infof("Conn[%s] create session %s success", s.UserConn.ID(), s.ID)
	if err = sw.Bridge(s.UserConn, srvConn); err != nil {
		logger.Error(err)
	}
}

func (s *Server) getServerConn() (srvconn.ServerConnection, error) {
	switch s.connOpts.Protocol {
	case srvconn.ProtocolSSH:
		return s.getSSHConn()
	default:
		return nil, ErrUnMatchProtocol
	}
}

func (s *Server) getSSHConn() (srvConn *srvconn.SSHConnection, err error) {
	loginAccount := s.account

	sshAuthOpts := make([]srvconn.SSHClientOption, 0, 6)
	sshAuthOpts = append(sshAuthOpts, srvconn.SSHClientUsername(loginAccount.Username))
	sshAuthOpts = append(sshAuthOpts, srvconn.SSHClientHost("121.40.251.109"))
	sshAuthOpts = append(sshAuthOpts, srvconn.SSHClientPort(22))
	sshAuthOpts = append(sshAuthOpts, srvconn.SSHClientTimeout(60))
	sshAuthOpts = append(sshAuthOpts, srvconn.SSHClientPassword(loginAccount.Secret))

	sshClient, err := srvconn.NewSSHClient(sshAuthOpts...)
	if err != nil {
		logger.Errorf("Get new ssh client err: %s", err)
		return nil, err
	}

	sess, err := sshClient.AcquireSession()
	if err != nil {
		logger.Errorf("SSH client(%s) start session err %s", sshClient, err)
		return nil, err
	}

	pty := s.UserConn.Pty()
	sshConnectOpts := make([]srvconn.SSHOption, 0, 6)
	sshConnectOpts = append(sshConnectOpts, srvconn.SSHCharset(""))
	sshConnectOpts = append(sshConnectOpts, srvconn.SSHTerm(pty.Term))
	sshConnectOpts = append(sshConnectOpts, srvconn.SSHPtyWin(srvconn.Windows{
		Width:  pty.Window.Width,
		Height: pty.Window.Height,
	}))

	sshConn, err := srvconn.NewSSHConnection(sess, sshConnectOpts...)
	if err != nil {
		_ = sess.Close()
		sshClient.ReleaseSession(sess)
		return nil, err
	}

	go func() {
		_ = sess.Wait()
		sshClient.ReleaseSession(sess)
		logger.Infof("SSH client(%s) shell connection release", sshClient)
	}()
	return sshConn, nil
}

//func (s *Server) sendConnectErrorMsg(err error) {
//    msg := fmt.Sprintf("%s error: %s", s.connOpts.ConnectMsg(),
//        s.ConvertErrorToReadableMsg(err))
//    utils.IgnoreErrWriteString(s.UserConn, msg)
//    utils.IgnoreErrWriteString(s.UserConn, utils.CharNewLine)
//    logger.Error(msg)
//    password := s.account.Secret
//    if password != "" {
//        passwordLen := len(s.account.Secret)
//        showLen := passwordLen / 2
//        hiddenLen := passwordLen - showLen
//        var msg2 string
//        if s.connOpts.Protocol == srvconn.ProtocolK8s {
//            msg2 = fmt.Sprintf("Try token: %s", password[:showLen]+strings.Repeat("*", hiddenLen))
//        } else {
//            msg2 = fmt.Sprintf("Try password: %s", password[:showLen]+strings.Repeat("*", hiddenLen))
//        }
//        logger.Error(msg2)
//    }
//
//}
