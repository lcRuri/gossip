package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"gossip/common"
	"gossip/connect"
	"log"
	"os"
)

func main() {
	fmt.Println("请输入配置文件位置:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	cfg := new(common.Config)
	//指定配置文件，控制台读入，否则默认
	err := ini.MapTo(cfg, input.Text())
	log.Println("reading config from " + input.Text())

	if err != nil {
		log.Println("load config failed,err:%v", err)
		return
	}

	nodeList := &connect.NodeList{
		Protocol:   cfg.NodeListConfig.Protocol,
		SecretKey:  cfg.NodeListConfig.SecretKey,
		IsPrint:    cfg.NodeListConfig.IsPrint,
		ListenAddr: cfg.NodeListConfig.ListenAddr,
	}

	//创建节点A及其本地节点列表
	nodeList.New(connect.Node{
		Addr:        cfg.NodeConfig.Addr,
		Port:        cfg.NodeConfig.Port,
		Name:        cfg.NodeConfig.Name,
		PrivateData: cfg.NodeConfig.PrivateData,
	})

	if len(cfg.ClusterNodeInfo.Addr) != 0 {
		nodeList.Set(connect.Node{
			Addr:        cfg.ClusterNodeInfo.Addr,
			Port:        cfg.ClusterNodeInfo.Port,
			Name:        cfg.ClusterNodeInfo.Name,
			PrivateData: cfg.ClusterNodeInfo.PrivateData,
		})
	}

	nodeList.Join()

	nodeList.Publish([]byte("test metadata conn"))

	select {}
}
