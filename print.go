package gossip

import "log"

//打印信息到控制台
func (nodeList *NodeList) Println(a ...interface{}) {
	//输出错误信息,即使不同步打印信息到控制台
	if a[0] == "[Error]:" && !nodeList.IsPrint {
		log.Println(a)
	}

	if nodeList.IsPrint {
		log.Println(a)
	}
}
