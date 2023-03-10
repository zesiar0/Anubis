package srvconn

import "io"

type ServerConnection interface {
	io.ReadWriteCloser
	SetWinSize(width, height int) error
	KeepAlive() error
}

type Windows struct {
	Width  int
	Height int
}

const (
	ProtocolSSH    = "ssh"
	ProtocolTELNET = "telnet"
	ProtocolK8s    = "k8s"

	ProtocolMySQL      = "mysql"
	ProtocolMariadb    = "mariadb"
	ProtocolSQLServer  = "sqlserver"
	ProtocolRedis      = "redis"
	ProtocolMongoDB    = "mongodb"
	ProtocolPostgreSQL = "postgresql"
	ProtocolClickHouse = "clickhouse"
)
