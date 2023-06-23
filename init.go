package gossip

import (
	"gossip/utils"
	"time"
)

// New 初始化
func (nodeList *NodeList) New(localNode Node) {
	//Addr 缺省值：0.0.0.0
	if localNode.Addr == "" {
		localNode.Addr = "0.0.0.0"
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
	nodeList.nodes.Store(localNode, time.Now().Unix()) //将本地节点信息添加进节点集合
	nodeList.LocalNode = localNode                     //初始化本地节点信息
	nodeList.status.Store(true)                        //初始化节点服务状态

	//设置元数据信息
	md := metadata{
		Data:   []byte(""), //元数据内容
		Update: 0,          //元数据更新时间戳
	}
	nodeList.metadata.Store(md) //初始化元数据信息
}

// Join 加入集群
func (nodeList *NodeList) Join() {
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return
	}

	//定时光广播本地节点信息
	go task(nodeList)

	//监听队列(UDP监听缓冲区)
	var mq = make(chan []byte, nodeList.Buffer)

	//监听来自于其他节点的信息，并且将信息放入mq队列
	go listener(nodeList, mq)

	//消费mq中的信息
	go consume(nodeList, mq)

	nodeList.Println("[Join]:", nodeList.LocalNode)
}

//Set 向本地节点列表加入其他节点
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

	nodeList.nodes.Store(node, time.Now().Unix())
}

// Get 获取本地节点列表
func (nodeList *NodeList) Get() []Node {
	//如果本地节点还未初始化
	if len(nodeList.LocalNode.Addr) == 0 {
		nodeList.Println("[Error]:", "Please use the New() function first")
		//直接返回
		return nil
	}

	var nodes []Node
	// 遍历所有sync.Map中的键值对
	nodeList.nodes.Range(func(k, v any) bool {
		//如果该节点超过一段数据没有更新
		if v.(int64)+nodeList.Timeout < time.Now().Unix() {
			nodeList.nodes.Delete(k)
		} else {
			nodes = append(nodes, k.(Node))
		}
		return true
	})

	return nodes
}
