package model

import "fmt"

type Asset struct {
	ID       string `json:"id"          gorm:"column:id"`
	Name     string `json:"name"        gorm:"column:name"`
	Address  string `json:"address"     gorm:"column:address"`
	Port     string `json:"port"        gorm:"column:port"`
	Protocol string `json:"protocol"    gorm:"column:protocol"`
	Platform string `json:"platform"    gorm:"column:platform"`
	IsActive bool   `json:"is_active"   gorm:"column:is_active"`
}

func (asset *Asset) String() string {
	return fmt.Sprintf("%s(%s:%s)", asset.Name, asset.Address, asset.Port)
}

const (
	ProtocolSSH   = "ssh"
	ProtocolMysql = "mysql"
)
