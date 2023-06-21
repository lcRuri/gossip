package gossip

import (
	"sync"
	"sync/atomic"
)

// Node 节点
type Node struct {
	Addr        string //ip地址
	Port        string //端口
	Name        string //节点名称
	PrivateData string //节点私有数据
}

// NodeList 节点列表
type NodeList struct {
	nodes sync.Map //节点集合（key为Node结构体，value为节点最近更新的秒级时间戳）

	Amount  int   //每次给多少个节点发生信息
	Cycle   int64 //同步时间周期（每隔多少秒向其他节点发送一次列表同步信息）
	Buffer  int   //UDP/TCP接收缓冲区大小（决定UDP/TCP监听服务可以异步处理多少个请求）
	Size    int   //单个UDP/TCP心跳数据包的最大容量（单位：字节）
	TimeOut int64 //单个节点的过期删除界限（多少秒后删除）

	SecretKey string //集群密钥，同一个集群持有相同的密钥

	localNode Node //本地节点信息

	Protocol   string //集群连接使用的网络协议，UDP或TCP，默认UDP
	ListenAddr string //本地UDP/TCP监听地址，用这个监听地址接收其他节点发来的心跳包（一般填0.0.0.0即可）

	status atomic.Value //本地节点列表更新状态（true：正常运行，false：停止发布心跳）

	IsPrint bool //是否打印列表同步信息到控制台

	metadata atomic.Value //元数据，集群中各个节点的元数据内容一致，相当于集群的公共数据（可存储一些公共配置信息），可以通过广播更新各个节点的元数据内容
}
