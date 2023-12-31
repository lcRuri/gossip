package connect

import (
	"gossip/utils"
	"strconv"
	"time"
)

// New 初始化
func (nodeList *NodeList) New(LocalNode Node) {
	//Addr 缺省值：0.0.0.0
	if LocalNode.Addr == "" {
		LocalNode.Addr = "0.0.0.0"
	}

	//Protocol 缺省值：UDP
	if nodeList.Protocol != "TCP" {
		nodeList.Protocol = "UDP"
	}

	//ListenAddr 缺省值：0.0.0.0
	if nodeList.ListenAddr == "" {
		nodeList.ListenAddr = "0.0.0.0"
	}

	//Amount 缺省值：3
	if nodeList.Amount == 0 {
		nodeList.Amount = 3
	}

	//Cycle 缺省值：6
	if nodeList.Cycle == 0 {
		nodeList.Cycle = 6
	}

	//Buffer 缺省值：不填则默认等于Amount乘3
	if nodeList.Buffer == 0 {
		nodeList.Buffer = nodeList.Amount * 3
	}

	//Size 缺省值：16384
	if nodeList.Size == 0 {
		nodeList.Size = 16384
	}

	//Timeout 缺省值：如果当前Timeout小于或等于Cycle，则自动扩大Timeout的值
	if nodeList.Timeout <= nodeList.Cycle {
		nodeList.Timeout = nodeList.Cycle*3 + 2
	}

	//如果密钥设置不为空，则对密钥进行md5加密
	if nodeList.SecretKey != "" {
		nodeList.SecretKey = utils.Md5Sign(nodeList.SecretKey)
	}

	//初始化本地节点列表的基础数据
	nodeList.Nodes.Store(LocalNode, time.Now().Unix()) //将本地节点信息添加进节点集合
	nodeList.LocalNode = LocalNode                     //初始化本地节点信息
	nodeList.Status.Store(true)                        //初始化节点服务状态

	//设置元数据信息
	md := metadata{
		Data:   []byte(""), //元数据内容
		Update: 0,          //元数据更新时间戳
	}
	nodeList.Metadata.Store(md) //初始化元数据信息
}

// Join 加入集群
func (nodeList *NodeList) Join() {
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	//定时广播本地节点信息
	go task(nodeList)

	//监听队列(UDP监听缓冲区)
	var mq = make(chan []byte, nodeList.Buffer)

	//监听来自于其他节点的信息，并且将信息放入mq队列
	go listener(nodeList, mq)

	//消费mq中的信息
	go consume(nodeList, mq)

	nodeList.Println("[Join]:", nodeList.LocalNode)
}

// Set 向本地节点列表加入其他节点
func (nodeList *NodeList) Set(node Node) {
	//先进行校验
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	//Addr 缺省值：0.0.0.0
	if node.Addr == "" {
		node.Addr = "0.0.0.0"
	}

	nodeList.Nodes.Store(node, time.Now().Unix())
}

// Get 获取本地节点列表
func (nodeList *NodeList) Get() []Node {
	//如果本地节点还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return nil
	}

	var Nodes []Node
	// 遍历所有sync.Map中的键值对
	nodeList.Nodes.Range(func(k, v any) bool {
		//如果该节点超过一段数据没有更新
		if v.(int64)+nodeList.Timeout < time.Now().Unix() {
			nodeList.Nodes.Delete(k)
		} else {
			Nodes = append(Nodes, k.(Node))
		}
		return true
	})

	return Nodes
}

// Stop 停止广播心跳
func (nodeList *NodeList) Stop() {
	//如果该节点的本地节点列表还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	nodeList.Println("[Stop]:", nodeList.LocalNode)
	nodeList.Status.Store(false)
}

// Start 重新开始广播心
func (nodeList *NodeList) Start() {
	//如果该节点的本地节点列表还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	//如果当前心跳服务正常
	if nodeList.Status.Load().(bool) {
		//返回
		return
	}
	nodeList.Println("[Start]:", nodeList.LocalNode)
	nodeList.Status.Store(true)
	//定时广播本地节点信息
	go task(nodeList)
}

// Read 读取本地节点列表的元数据信息
func (nodeList *NodeList) Read() []byte {

	//如果该节点的本地节点列表还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return nil
	}

	return nodeList.Metadata.Load().(metadata).Data
}

// Publish 在集群中发布新的元数据信息
func (nodeList *NodeList) Publish(newMetadata []byte) {
	//如果该节点的本地节点列表还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	nodeList.Println("[Publish]:", nodeList.LocalNode, "/ [Metadata]:", string(newMetadata))

	//将本地节点加入已传染的节点列表infected
	var infected = make(map[string]bool)
	infected[nodeList.LocalNode.Addr+":"+strconv.Itoa(nodeList.LocalNode.Port)] = true

	//更新本地节点信息
	nodeList.Set(nodeList.LocalNode)

	//设置新的元数据信息
	md := metadata{
		Data:   newMetadata,           //元数据内容
		Update: time.Now().UnixNano(), //元数据更新时间戳
	}

	//更新本地节点的元数据信息
	nodeList.Metadata.Store(md)

	//设置心跳数据包
	p := packet{
		Node:     nodeList.LocalNode,
		Infected: infected,

		//将数据包设为元数据更新数据包
		Metadata: md,
		IsUpdate: true,

		SecretKey: nodeList.SecretKey,
	}

	//在集群中广播数据包
	broadcast(nodeList, p)
}
