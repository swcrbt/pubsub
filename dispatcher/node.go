package dispatcher

import (
	"net"
	"time"
)

type Node struct {
	Address    string    `json:"address"`
	Host       string    `json:"host"`
	Port       string    `json:"port"`
	CreateTime time.Time `json:"createtime"`
}

func NewNode(address string) *Node {
	host, port, _ := net.SplitHostPort(address)

	return &Node{
		Address:    address,
		Host:       host,
		Port:       port,
		CreateTime: time.Now(),
	}
}
