package rpcsupport

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 起服务器
//  传入方法 发布再服务上
func ServeRpc(host string, service interface{}) error {
	// 将方法发布到默认服务器上
	rpc.Register(service)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error: %s", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
	return nil
}

// 连接上 server
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	client := jsonrpc.NewClient(conn)
	return client, err
}
