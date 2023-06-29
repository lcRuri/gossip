package test

import (
	"fmt"
	"gossip/connect"
	"testing"
	"time"
)

// 完整测试用例-在本地启动四个节点，构成一个Gossip集群（UDP）
func TestUDPCluster(t *testing.T) {

	fmt.Println("---- Start a gossip cluster (UDP) ----")

	//使用UDP连接集群节点
	protocol = "UDP"

	//先启动节点A（初始节点）
	nodeA()
	//启动节点B
	nodeB()
	//启动节点C
	nodeC()
	//启动节点D
	nodeD()

	//延迟10秒
	time.Sleep(10 * time.Second)

	//结束测试
	fmt.Println("---- End ----")
}

var protocol string

// 运行节点A（初始节点）
func nodeA() *connect.NodeList {
	//配置节点A的本地节点列表nodeList参数
	nodeList := &connect.NodeList{
		Protocol:   protocol,
		SecretKey:  "test_key",
		IsPrint:    true,
		ListenAddr: "0.0.0.0",
	}

	//创建节点A及其本地节点列表
	nodeList.New(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8000,
		Name:        "A-server",
		PrivateData: "test-data",
	})

	//因为是第一个启动的节点，所以不需要用Set函数添加其他节点

	//本地节点加入Gossip集群，本地节点列表与集群中的各个节点所存储的节点列表进行数据同步
	nodeList.Join()

	//延迟3秒
	time.Sleep(3 * time.Second)

	return nodeList
}

// 运行节点B
func nodeB() *connect.NodeList {
	//配置节点B的本地节点列表nodeList参数
	nodeList := &connect.NodeList{
		Protocol:   protocol,
		SecretKey:  "test_key",
		IsPrint:    true,
		ListenAddr: "0.0.0.0",
	}

	//创建节点B及其本地节点列表
	nodeList.New(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8001,
		Name:        "B-server",
		PrivateData: "test-data",
	})

	//将初始节点A的信息加入到B节点的本地节点列表当中
	nodeList.Set(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8000,
		Name:        "A-server",
		PrivateData: "test-data",
	})

	//调用Join后，节点B会自动与节点A进行数据同步
	nodeList.Join()

	nodeList.Publish([]byte("test metadata B"))

	metadata := nodeList.Read()
	fmt.Println("Metadata:", string(metadata)) //打印元数据信息
	//延迟10秒
	time.Sleep(10 * time.Second)

	return nodeList

}

// 运行节点C
func nodeC() *connect.NodeList {
	//配置节点C的本地节点列表nodeList参数
	nodeList := &connect.NodeList{
		Protocol:   protocol,
		SecretKey:  "test_key",
		IsPrint:    true,
		ListenAddr: "0.0.0.0",
	}

	//创建节点C及其本地节点列表
	nodeList.New(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8002,
		Name:        "C-server",
		PrivateData: "test-data",
	})

	//nodeList.Set(gossip.Node{
	//	Addr:        "0.0.0.0",
	//	Port:        8000,
	//	Name:        "A-server",
	//	PrivateData: "test-data",
	//})

	nodeList.Set(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8001,
		Name:        "B-server",
		PrivateData: "test-data",
	})

	//在加入集群后，节点C将会与上面的节点A及节点B进行数据同步
	nodeList.Join()

	//延迟10秒
	time.Sleep(10 * time.Second)

	//获取本地节点列表
	list := nodeList.Get()
	fmt.Println("Node list::", list) //打印节点列表

	//在集群中发布新的元数据信息
	nodeList.Publish([]byte("test metadata"))

	//读取本地元数据信息
	metadata := nodeList.Read()
	fmt.Println("Metadata:", string(metadata)) //打印元数据信息

	//停止节点C的心跳广播服务（节点C暂时下线）
	nodeList.Stop()

	//延迟30秒
	time.Sleep(30 * time.Second)

	//因为之前节点C下线，C的本地节点列表无法接收到各节点的心跳数据包，列表被清空
	//所以要先往C的本地节点列表中添加一些集群节点，再调用Start()重启节点D的同步工作
	nodeList.Set(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8001,
		Name:        "B-server",
		PrivateData: "test-data",
	})

	//重启节点C的心跳广播服务（节点C重新上线）
	nodeList.Start()

	//读取本地元数据信息
	metadata = nodeList.Read()
	fmt.Println("Metadata:", string(metadata)) //打印元数据信息

	return nodeList
}

// 运行节点D
func nodeD() *connect.NodeList {
	//配置节点D的本地节点列表nodeList参数
	nodeList := &connect.NodeList{
		Protocol:   protocol,
		SecretKey:  "test_key",
		IsPrint:    true,
		ListenAddr: "0.0.0.0",
	}

	//创建节点D及其本地节点列表
	nodeList.New(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8003,
		Name:        "D-server",
		PrivateData: "test-data",
	})

	nodeList.Set(connect.Node{
		Addr:        "0.0.0.0",
		Port:        8000,
		Name:        "A-server",
		PrivateData: "test-data",
	})

	//调用Join后，节点D会自动与节点A进行数据同步
	nodeList.Join()

	//延迟5秒
	time.Sleep(5 * time.Second)

	//读取本地元数据信息
	metadata := nodeList.Read()

	fmt.Println("Metadata Byte:", metadata)
	fmt.Println("Metadata:", string(metadata)) //打印元数据信息

	return nodeList
}
