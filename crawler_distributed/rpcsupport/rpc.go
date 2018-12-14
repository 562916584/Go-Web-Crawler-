package rpcsupport

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 起服务器
//  传入方法 发布在服务上
func ServeRpc(host string, service interface{}) error {
	// 将方法发布到默认服务器上
	rpc.Register(service)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s", host)
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

// 用于连接上 server
// 返回连接上服务的 client
func NewClient(host string) (*rpc.Client, error) {
	// 返回连接客户端
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	// 返回jsonRpc客户端
	client := jsonrpc.NewClient(conn)
	return client, err
}
