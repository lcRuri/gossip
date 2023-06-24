package gnet

import (
	"fmt"
	"gossip"
	"net"
)

func udpWrite(nodeList *gossip.NodeList, addr string, port int, data []byte) {
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
	})
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}

	_, err = socket.Write(data) // 发送数据
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}

	defer func(socket *net.UDPConn) {
		err = socket.Close()
		if err != nil {
			nodeList.Println("[Error]:", err)
		}
	}(socket)
}

func udpListen(nodeList *gossip.NodeList, mq chan []byte) {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", nodeList.ListenAddr, nodeList.LocalNode.Port))
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		nodeList.Println("[Error]:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err = conn.Close()
		if err != nil {
			nodeList.Println("[Error]:", err)
		}
	}(conn)

	for {
		//接收数组
		bs := make([]byte, nodeList.Size)

		//从UDP监听中接收数据
		n, _, err := conn.ReadFromUDP(bs)
		if err != nil {
			nodeList.Println("[Error]:", err)
			continue
		}

		if n >= nodeList.Size {
			nodeList.Println("[Error]:", fmt.Sprintf("received data size (%v) exceeds the limit (%v)", n, nodeList.Size))
			continue
		}

		//获取有效数据
		b := bs[:n]

		//将数据放入缓冲队列，异步处理数据
		mq <- b
	}
}
