package gossip

import "gossip/utils"

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

}
