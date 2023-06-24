package gossip

// Write  发送数据
func Write(nodeList *NodeList, addr string, port int, data []byte) {
	if nodeList.Protocol != "TCP" {
		udpWrite(nodeList, addr, port, data)
	} else {
		tcpWrite(nodeList, addr, port, data)
	}
}

// Listen 服务端监听
func Listen(nodeList *NodeList, mq chan []byte) {
	if nodeList.Protocol != "TCP" {
		udpListen(nodeList, mq)
	} else {
		tcpListen(nodeList, mq)
	}
}
