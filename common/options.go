package common

type Config struct {
	NodeListConfig  NodeListConfig  `ini:"nodeListConfig""`
	NodeConfig      NodeConfig      `ini:"nodeConfig""`
	ClusterNodeInfo ClusterNodeInfo `ini:"clusterNodeInfo"`
}

type NodeListConfig struct {
	Protocol   string `ini:"protocol"`
	SecretKey  string `ini:"secretKey"`
	IsPrint    bool   `ini:"isPrint"`
	ListenAddr string `ini:"ListenAddr"`
}

type NodeConfig struct {
	Addr        string `ini:"addr"`
	Port        int    `ini:"port"`
	Name        string `ini:"name"`
	PrivateData string `ini:"privateData"`
}

type ClusterNodeInfo struct {
	Addr        string `ini:"addr"`
	Port        int    `ini:"port"`
	Name        string `ini:"name"`
	PrivateData string `ini:"privateData"`
}
