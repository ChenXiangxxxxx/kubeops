package config

import "time"

const (
	ListenAddr     = "0.0.0.0:9090"
	Kubeconfigs    = `{"TST-1":"D:\\gocode\\cmdb\\kube-backend\\kubeconfig","TST-2":"D:\\gocode\\cmdb\\kube-backend\\kubeconfig"}`
	PodLogTailLine = 5000 //查看容器日志时显示的tail行数
	//数据库配置
	DbType = "mysql"
	DbHost = "192.168.1.11"
	DbPort = 3306
	DbName = "kubeops"
	DbUser = "root"
	DbPwd  = "123456"
	//打印mysql debug sql日志
	LogMode = false
	//连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
	//helm上传路径
	UploadPath = "D:\\custom\\project\\123\\实战\\"
	//账号密码
	AdminUser = "admin"
	AdminPwd  = "123456"
)
