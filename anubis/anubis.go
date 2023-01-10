package anubis

import (
	"Anubis/pkg/config"
	"Anubis/pkg/db"
	"Anubis/pkg/logger"
	"Anubis/pkg/sshd"
	"os"
	"os/signal"
	"syscall"
)

type Anubis struct {
	sshSrv *sshd.Server
}

func Run(configPath string) {
	config.Setup(configPath)
	bootstrap()
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	srv := &server{}
	sshSrv := sshd.NewSSHServer(srv)
	sshSrv.Start()
	<-gracefulStop
	sshSrv.Stop()
}

func bootstrap() {
	logger.Initial()
	db.Initial()
}
