package main

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"gossip/common"
	"gossip/connect"
	"log"
	"net/http"
	"os"
	"strconv"
)

var nodeList *connect.NodeList

func init() {
restart:
	fmt.Println("请输入配置文件位置:")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	cfg := new(common.Config)
	//指定配置文件，控制台读入，否则默认
	err := ini.MapTo(cfg, input.Text())
	log.Println("reading config from " + input.Text())

	if err != nil {
		log.Printf("load config failed,err:%v", err)
		goto restart
	}

	nodeList = &connect.NodeList{
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

}

func publish(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := request.URL.Query().Get("publish")
	nodeList.Publish([]byte(key))
	return
}

// 开启一层http服务，可以通过http来调用start stop get set publish read等方法
func main() {
	http.HandleFunc("/gossip/publish", publish)

	//启动http服务
	_ = http.ListenAndServe("localhost:"+strconv.Itoa(nodeList.LocalNode.Port+100), nil)
}
