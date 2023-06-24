package gossip

import (
	"fmt"
	"net"
)

func tcpWrite(nodeList *NodeList, addr string, port int, data []byte) {

	tcpAddr := fmt.Sprintf("%s:%v", addr, port)

	server, err := net.ResolveTCPAddr("tcp4", tcpAddr)

	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}

	//与服务器建立连接
	conn, err := net.DialTCP("tcp", nil, server)
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}

	//向服务器发送信息
	_, err = conn.Write(data)
	if err != nil {
		nodeList.Println(err)

	}

	defer func(conn *net.TCPConn) {
		err = conn.Close()
		if err != nil {
			nodeList.Println("[Error]:", err)
		}
	}(conn)

}

func tcpListen(nodeList *NodeList, mq chan []byte) {
	//节点列表监听的地址和本地节点的端口
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%v", nodeList.ListenAddr, nodeList.LocalNode.Port))
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}

	defer func(server net.Listener) {
		err = server.Close()
		if err != nil {
			nodeList.Println("[Error]:", err)
		}
	}(server)

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}

		go func() {

			//接收数组
			bs := make([]byte, nodeList.Size)
			n, err := conn.Read(bs)
			if err != nil {
				nodeList.Println("[Error]:", err)
				return
			}

			//如果接受的字节数大于心跳数据包的最大容量，error
			if n >= nodeList.Size {
				nodeList.Println("[Error]:", fmt.Sprintf("received data size (%v) exceeds the limit (%v)", n, nodeList.Size))
				return
			}

			//获取有效数据
			b := bs[:n]

			//将数据放去缓冲队列，异步处理数据
			mq <- b
		}()
	}
}
