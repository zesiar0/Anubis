package proxy

import (
	"Anubis/pkg/logger"
	"Anubis/pkg/srvconn"
	"bytes"
	"context"
	"time"
)

type SwitchSession struct {
	ID string

	MaxIdleTime   int
	keepAliveTime int

	ctx    context.Context
	cancel context.CancelFunc

	p *Server
}

func (s *SwitchSession) Bridge(userConn UserConnection, srvConn srvconn.ServerConnection) (err error) {
	srvInChan := make(chan []byte, 1)
	usrInChan := make(chan []byte, 1)
	done := make(chan struct{})

	usrOutputChan, srvOutputChan := parse(usrInChan, srvInChan)
	go func() {
		exitFlag := false

		for {
			buf := make([]byte, 1024)
			nr, err := srvConn.Read(buf)
			if nr > 0 {
				select {
				case srvInChan <- buf[:nr]:
				case <-done:
					exitFlag = true
					logger.Infof("Session[%s] done", s.ID)
				}
				if exitFlag {
					break
				}
			}

			if err != nil {
				logger.Errorf("Session[%s] srv read err: %s", s.ID, err)
				break
			}
		}
		logger.Infof("Session[%s] srv read end", s.ID)
		close(srvInChan)
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			nr, err := userConn.Read(buf)
			logger.Infof("user write message: %s", buf)
			if nr > 0 {
				index := bytes.IndexFunc(buf[:nr], func(r rune) bool {
					return r == '\r' || r == '\n'
				})
				if index < 0 {
					select {
					case <-done:
					case usrInChan <- buf[:nr]:
					}
				} else {
					select {
					case <-done:
					case usrInChan <- buf[:index]:
					}
					time.Sleep(time.Millisecond * 1000)
					select {
					case <-done:
					case usrInChan <- buf[index:nr]:
					}
				}
				if err != nil {
					logger.Errorf("Session[%s] user read err: %s", s.ID, err)
					break
				}
			}
			logger.Infof("Session[%s] user read end", s.ID)
		}
	}()

	for {
		select {
		case msg, ok := <-srvOutputChan:
			if !ok {
				return
			}
			if _, err := userConn.Write(msg); err != nil {
				logger.Errorf("Session usrConn write err: %s", err)
			}
		case msg, ok := <-usrOutputChan:
			if !ok {
				return
			}
			if _, err := srvConn.Write(msg); err != nil {
				logger.Errorf("Session srvConn write err: %s", err)
			}
		case <-userConn.Context().Done():
			logger.Infof("Session[%s]: user conn context done", s.ID)
			return nil
		}
	}
}

func parse(usrInChan <-chan []byte, srvInChan <-chan []byte) (usrOut, srvOut <-chan []byte) {
	usrOutputChan := make(chan []byte, 1)
	srvOutputChan := make(chan []byte, 1)

	go func() {
		defer func() {
			close(usrOutputChan)
			close(srvOutputChan)

			logger.Info("Session done")
		}()

		for {
			select {
			// 用户输入
			case msg, ok := <-usrInChan:
				if !ok {
					return
				}
				select {
				case usrOutputChan <- msg:
				}
			case msg, ok := <-srvInChan:
				if !ok {
					return
				}
				select {
				case srvOutputChan <- msg:
				}

			}
		}
	}()
	return usrOutputChan, srvOutputChan
}
