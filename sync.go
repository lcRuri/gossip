package gossip

import (
	"strconv"
	"time"
)

//定时进行心跳广播
func task(nodeList *NodeList) {
	for {
		//停止同步
		if nodeList.status.Load().(bool) {
			break
		}

		//将本地节点加入已经传染的节点列表
		var infected = make(map[string]bool)
		infected[nodeList.LocalNode.Addr+":"+strconv.Itoa(nodeList.LocalNode.Port)] = true

		//更新本地节点信息
		nodeList.Set(nodeList.LocalNode)

		//设置心跳数据包
		p := packet{
			Node:      nodeList.LocalNode,
			Infected:  infected,
			SecretKey: nodeList.SecretKey,
		}

		//广播心跳数据包
		broadcast(nodeList, p)

		//向集群中某个节点发起数据同步
		swapRequest(nodeList)

		if nodeList.IsPrint {
			nodeList.Println("[Listen]:", nodeList.ListenAddr+":"+strconv.Itoa(nodeList.LocalNode.Port), "/ [Node list]:", nodeList.Get())
		}

		//间隔时间
		time.Sleep(time.Duration(nodeList.Cycle) * time.Second)
	}
}

func swapRequest(nodeList *NodeList) {

}

func broadcast(nodeList *NodeList, p packet) {

}

func consume(nodeList *NodeList, mq chan []byte) {

}

func listener(nodeList *NodeList, mq chan []byte) {

}
