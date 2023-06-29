package connect

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

// 定时进行心跳广播
func task(nodeList *NodeList) {
	for {
		//停止同步
		//一开始少了！
		if !nodeList.Status.Load().(bool) {
			//fmt.Println("停止同步")
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
		//fmt.Println("------broadcast------")
		broadcast(nodeList, p)

		//向集群中某个节点发起数据同步
		//fmt.Println("------swapRequest------")
		swapRequest(nodeList)

		if nodeList.IsPrint {
			//nodeList.Println("[Listen]:", nodeList.ListenAddr+":"+strconv.Itoa(nodeList.LocalNode.Port), "/ [Node list]:", nodeList.Get())
		}

		//间隔时间
		time.Sleep(time.Duration(nodeList.Cycle) * time.Second)
	}
}

// 消费信息
func consume(nodeList *NodeList, mq chan []byte) {
	for {
		//从监听队列中取出信息
		bs := <-mq
		var p packet
		err := json.Unmarshal(bs, &p)

		if err != nil {
			nodeList.Println("[Error]:", err)
			continue
		}

		//如果数据包密钥与当前节点不匹配
		if p.SecretKey != nodeList.SecretKey {
			nodeList.Println("[Error]:", "The secretKey do not match")
			//跳过，不处理该数据包
			continue
		}

		//如果该数据包是两节点间的元数据交换数据包
		if p.IsSwap != 0 {
			//.Printf("two node swap\n")
			//如果数据包中的元数据版本要比本地存储的元数据版本新
			//则进行更新
			if p.Metadata.Update > nodeList.Metadata.Load().(metadata).Update {
				//更新本地节点中存储的元数据信息
				log.Printf("[Updata]:%v update metadata", nodeList.LocalNode.Name)
				nodeList.Metadata.Store(p.Metadata)
				md := nodeList.Metadata.Load().(metadata)
				log.Printf("[%v] metadata:%v\n", nodeList.LocalNode.Name, string(md.Data))
				//跳过，不广播，不回应发起方
				continue
			}

			//如果数据包中的元数据版本要比本地存储的元数据版本旧，说明发起方的元数据版本较旧，发起方需要更新
			if p.Metadata.Update < nodeList.Metadata.Load().(metadata).Update {
				//如果发起方发出的数据交换请求
				if p.IsSwap == 1 {
					//回应发起方，向发起方发送最新的元数据信息，完成交换流程
					swapResponse(nodeList, p.Node)
				}
			}

			//跳过，不广播
			continue

		}

		//从心跳包里面取出节点信息
		node := p.Node

		//更新本地列表
		nodeList.Set(node)

		//如果该数据包是元数据更新数据包，且版本更新
		if p.IsUpdate && p.Metadata.Update > nodeList.Metadata.Load().(metadata).Update {
			//fmt.Printf("[%v] consume metadata \n", nodeList.LocalNode.Name)
			//更新本地节点中存储的元数据信息
			log.Printf("[Updata]:%v update metadata", nodeList.LocalNode.Name)
			nodeList.Metadata.Store(p.Metadata)
			md := nodeList.Metadata.Load().(metadata)
			log.Printf("[%v] metadata:%v\n", nodeList.LocalNode.Name, string(md.Data))
		}

		//广播推送该节点信息
		broadcast(nodeList, p)
	}
}

// 广播推送信息
func broadcast(nodeList *NodeList, p packet) {

	//取出所有未过期的节点
	nodes := nodeList.Get()
	//fmt.Printf("[%v] all active nodes:%v\n", nodeList.LocalNode.Name, nodes)
	//本次广播的目标节点列表
	var targetNodes []Node

	//选取部分未被传染的节点
	i := 0

	for _, v := range nodes {

		//如果已经超过最大推送数量
		if i > nodeList.Amount {
			//结束广播
			break
		}

		//如果该节点已经被传染过了
		if p.Infected[v.Addr+":"+strconv.Itoa(v.Port)] {
			//跳过该节点
			continue
		}

		p.Infected[v.Addr+":"+strconv.Itoa(v.Port)] = true //标记该节点为已传染状态

		//设置发送目标节点
		targetNode := Node{
			Addr: v.Addr, //设置发送目标地址
			Port: v.Port, //设置发送目标端口
		}

		//将该节点添加进广播列表
		targetNodes = append(targetNodes, targetNode)
		i++
	}

	//fmt.Printf("[%v] all target nodes:%v\n", nodeList.LocalNode.Name, targetNodes)
	//向这些未被传染的节点广播传染数据
	for _, v := range targetNodes {
		bs, err := json.Marshal(p)
		if err != nil {
			nodeList.Println("[Error]:", err)
		}
		//fmt.Printf("%v send to %v\n", nodeList.LocalNode.Name, v.Port)
		//发送
		Write(nodeList, v.Addr, v.Port, bs)
	}
}

// 监听其他节点发来的同步信息
func listener(nodeList *NodeList, mq chan []byte) {
	//监听协程
	Listen(nodeList, mq)
}

// 发起两节点数据交换请求
func swapRequest(nodeList *NodeList) {
	//设置为数据交换数据包
	p := packet{
		//将本地节点信息存入数据包，接收方根据这个信息回复请求
		Node:      nodeList.LocalNode,
		Infected:  make(map[string]bool),
		IsSwap:    1,
		Metadata:  nodeList.Metadata.Load().(metadata),
		SecretKey: nodeList.SecretKey,
	}

	//取出所有未过期的节点
	nodes := nodeList.Get()

	bs, err := json.Marshal(p)
	if err != nil {
		nodeList.Println("[Error]:", err)
	}

	//在节点列表中随机选取一个节点，发起数据交换请求
	for i := 0; i < len(nodes); i++ {
		//如果遍历到自己，则跳过
		if nodes[i].Addr == nodeList.LocalNode.Addr && nodes[i].Port == nodeList.LocalNode.Port {
			continue
		}

		//fmt.Printf("[%v] Join the cluster\n", nodeList.LocalNode.Name)
		//fmt.Printf("[%v] send request to [%v]", nodeList.LocalNode.Port, nodes[i].Port)

		//发送请求
		Write(nodeList, nodes[i].Addr, nodes[i].Port, bs)

		if nodeList.IsPrint {
			//nodeList.Println("[Swap Request]:", nodeList.LocalNode.Addr+":"+strconv.Itoa(nodeList.LocalNode.Port), "->", nodes[i].Addr+":"+strconv.Itoa(nodes[i].Port))
		}
		break
	}
}

// 接收数据交换请求并回应发送方，完成交换工作
func swapResponse(nodeList *NodeList, node Node) {
	//设置为数据交换数据包
	p := packet{
		Node:      nodeList.LocalNode,
		Infected:  make(map[string]bool),
		IsSwap:    2,
		Metadata:  nodeList.Metadata.Load().(metadata),
		SecretKey: nodeList.SecretKey,
	}

	bs, err := json.Marshal(p)
	if err != nil {
		nodeList.Println("[Error]:", err)
	}

	//回应发起节点
	Write(nodeList, node.Addr, node.Port, bs)

	if nodeList.IsPrint {
		nodeList.Println("[Swap Response]:", node.Addr+":"+strconv.Itoa(node.Port), "<-", nodeList.LocalNode.Addr+":"+strconv.Itoa(nodeList.LocalNode.Port))
	}
}
