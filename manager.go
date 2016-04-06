//
//  main
//

package main

import (
	"fmt"
	"gomore/global"
	"gomore/golog"
	"gomore/hub"
	//"net"
	"gomore/admin"
	"gomore/area"
	"gomore/connector"
	"gomore/lib/syncmap"
	"gomore/worker"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"runtime"
	"time"
	//z_type "gomore/type"
)

// 启动一个测试的php worker以处理业务流程
func stop_php_worker() {

	c := exec.Command("/bin/sh", "-c", `ps -ef |grep "worker/php/workers.php"  |awk \'{print $2}\' |xargs -i kill -9 {} `)
	d, _ := c.Output()

	golog.Info("Stop_php_worker: ", string(d))

	time.Sleep(time.Second * 1)

}

// 启动一个测试的php worker以处理业务流程
func start_php_worker() {

	stop_php_worker()
	wd, _ := os.Getwd()
	work_num, _ := global.ConfigJson.GetString("worker", "worker_num")
	argv := []string{fmt.Sprintf("%s/worker/php/workers.php", wd), "start", work_num}
	golog.Info("Argv:", argv)
	c := exec.Command("/usr/bin/php", argv...)
	d, _ := c.Output()
	golog.Info("Start_php_worker: ", string(d))

	time.Sleep(time.Second * 1)

}

// 初始化全局变量
func init_global() {

	global.SumConnections = 0
	global.Qps = 0
	/*
		//worker_nbrs, _ := global.ConfigJson.GetInt64("hub", "number")
		//golog.Error("worker_nbrs:", worker_nbrs)
		tmp := int(worker_nbrs)
		var i int

		for i = 1; i <= tmp; i++ {
			global.WorkerNbrs = append(global.WorkerNbrs, fmt.Sprintf("%d", i))

		}
	*/
	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.Channels = make(map[string]string)

	// global.RpcChannels  =  make(map[string] *z_type.ChannelRpcType )

	global.SyncUserConns = syncmap.New()
	global.SyncUserSessions = syncmap.New()
	global.SyncUserWebsocketConns = syncmap.New()
	global.SyncUserJoinedChannels = syncmap.New()
	global.SyncCrons = syncmap.New()

	global.PackSplitType = `bufferio`

	hub.RedisInit()

	hub.LoadSessionFromRedis()

}

/**
 * zeromore 框架启动
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		fmt.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	golog.InitLogger()
	global.InitConfig()

	//fmt.Println( global.Config.Connector.MaxConections )
	//fmt.Println( global.Config.Area.Init_area )
	//fmt.Println( global.Config.WorkerAgent.Host )
	//hub.CronTest()
	//return
	//appConfig := &AppConfig{}
	// 读取配置文件

	init_global()
	go connector.SocketConnector("", global.Config.Connector.SocketPort)
	//go connector.WebsocketConnector("", global.Config.Connector.WebsocketPort)

	// 开启worker代理
	go worker.Worker_agent()

	// 开启hub服务器
	go hub.HubServer()

	// 预创建多个场景
	for _, area_id := range global.Config.Area.Init_area {
		area.CreateChannel(area_id, area_id)
		global.Channels[area_id] = global.Config.Hub.Hub_host
	}

	// 启动worker
	//go start_php_worker()
	go worker.Start()
	golog.Info("Server started!")

	go admin.HttpServer()

	select {}

}
