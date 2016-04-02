package worker

import (
	"math/rand"
	"strings"
	"sync/atomic"
	"gomore/area"
	"gomore/global"
	"gomore/protocol"

	log "github.com/Sirupsen/logrus"
	//"github.com/antonholmquist/jason"
	//sync"
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// workeragent连接对象切片版
var WorkerAgentConnsIndexSlice = make([]string, 0, 1000)
var WorkerAgentConnsSlice = make([]*net.TCPConn, 0, 1000)

var ReqAgentConnsIndexSlice = make([]string, 0, 1000)
var ReqAgentConnsSlice = make([]*net.TCPConn, 0, 1000)

/**
 * 监听客户端连接
 */
func Worker_agent() {

	worker_agent_host := global.Config.WorkerAgent.Host

	worker_agent_port := global.Config.WorkerAgent.Port

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(worker_agent_host), int(worker_agent_port), ""})
	if err != nil {
		log.Error("ListenTCP Exception:", err.Error())
		return
	}

	log.Error("Worker agent  Server :", worker_agent_host, worker_agent_port)

	WorkerAgentListen(listen)
}

/**
 *  处理客户端连接
 */
func WorkerAgentListen(listen *net.TCPListener) {

	//go stat_kick()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Error("AcceptTCP Exception::", err.Error(), time.Now().UnixNano())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		defer conn.Close()
		//conn.SetNoDelay(false)
		log.Info("RemoteAddr:", conn.RemoteAddr().String())

		//go handleWorkerWithJson( conn  )
		go handleWorkerWithBufferio(conn)

	} //end for {

}

// D:\php7\php7.exe  D:\gopath\src\zeromore\worker\php\workers.php 1

func workerPingKick(conn *net.TCPConn, worker_idf string) {

	timer := time.Tick(time.Second * 2)
	for _ = range timer {
		if worker_idf == "" {
			continue
		}
		ping := fmt.Sprintf(`worker.ping||%s||%s||%d`, "", worker_idf, time.Now().Unix())
		fmt.Println(ping)
		_, err := conn.Write([]byte(ping + "\n"))
		if err != nil {
			break
		}
	}
}

func handleWorkerWithBufferio(conn *net.TCPConn) {

	var worker_idf string
	var client_idf string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	worker_idf = fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(99999))
	client_idf = ""

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {

		msg, err := reader.ReadString('\n')

		if err != nil {
			//fmt.Println( "handleWorker connection error: ", err.Error())
			// 超时处理
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {

			}
			closeWorker(conn, client_idf, worker_idf)
			break

		}
		if msg == "" {
			continue
		}
		//fmt.Println("worker agent from ", worker_idf, " receive 3:", msg)

		go workeDispath(msg, conn, &client_idf, &worker_idf)

	}

}

func workeDispath(msg string, conn *net.TCPConn, client_idf *string, worker_idf *string) (int, error) {

	//ret_json, _ := jason.NewObjectFromBytes(msg)
	//cmd, _ := ret_json.GetString("cmd")

	ret_arr := strings.Split(msg, "||")
	//fmt.Println("workeDispath ret_arr 4:", ret_arr)
	if len(ret_arr) < 4 {
		//fmt.Println("workeDispath data length error!")
		return 0, nil
	} else {
		if len(ret_arr) > 4 {

			for i := 4; i < int(len(ret_arr)); i++ {

				ret_arr[3] = ret_arr[3] + ret_arr[i]
			}

		}
	}

	cmd := ret_arr[0]
	*client_idf = ret_arr[1]
	*worker_idf = ret_arr[2]
	worker_data := ret_arr[3]
	data := ""

	// 前端连接到代理
	if cmd == global.DATA_REQ_CONNECT {

		//data_byte, _ = protocol.Packet(`{"cmd":"req.connect","ret":200,"msg":"ok"}`)
		data := fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_REQ_CONNECT, "", "", "ok")
		return conn.Write([]byte(data))

	}
	// 前端发数据过来
	if cmd == global.DATA_REQ_MSG {

		//fmt.Println( "workeDispath idf  :", *worker_idf )
		// @todo应该自动分配worker

		worker_conn := GetWorkerConn(*worker_idf)
		log.Info(*worker_idf, worker_conn)
		if worker_conn == nil {
			data = fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_REQ_MSG, *client_idf, *worker_idf, "worker no found!")
			fmt.Println("worker_conn ", *worker_idf, " no found!")
			return conn.Write([]byte(data + "\n"))
		}
		data = fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_REQ_MSG, *client_idf, *worker_idf, worker_data)
		log.Info("workeDispath worker_data :", data)
		return worker_conn.Write([]byte(data + "\n"))
	}
	// worker连接到代理
	if cmd == global.DATA_WORKER_CONNECT {
		AddWorkerConn(*worker_idf, conn)
		data = fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_WORKER_CONNECT, *client_idf, *worker_idf, "ok")
		fmt.Println("worker.connect : ", *worker_idf)
		conn.Write([]byte(data + "\n"))
		go workerPingKick(conn, *worker_idf)

	}
	// worker返回处理结果
	if cmd == global.DATA_WORKER_REPLY {

		//fmt.Println( "worker.reply data 7 :" , ret_arr );
		req_conn := area.GetConn(*client_idf)
		if req_conn == nil {
			fmt.Println("req_conn :", *client_idf, " no found")
		}

		//如果没有找到客户端连接对象则报错并返回
		atomic.AddInt64(&global.Qps, 1)
		if req_conn != nil {
			worker_data_byte, _ := protocol.Packet(worker_data)
			log.Info(worker_data_byte)
			//fmt.Println( "worker.reply  worker_data :" ,  worker_data  );
			return req_conn.Write([]byte(worker_data + "\n"))
		}

		data = fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_WORKER_REPLY, *client_idf, *worker_idf, "req conn no found!")
		return conn.Write([]byte(data + "\n"))
	}

	return 1, nil

}

func closeWorker(conn *net.TCPConn, client_idf string, worker_idf string) {

	conn.Close()
	DeleteWorkerConn(worker_idf)
	DeleteReqConn(client_idf)

}

func GetReqConn(idf string) *net.TCPConn {

	for i := 0; i < int(len(ReqAgentConnsIndexSlice)); i++ {

		if ReqAgentConnsIndexSlice[i] == idf {
			return ReqAgentConnsSlice[i]
		}
	}
	return nil

}

func GetWorkerConn(worker_idf string) *net.TCPConn {

	for i := 0; i < int(len(WorkerAgentConnsIndexSlice)); i++ {
		if WorkerAgentConnsIndexSlice[i] == worker_idf {
			return WorkerAgentConnsSlice[i]
		}
	}
	return nil

}

func AddReqConn(idf string, conn *net.TCPConn) {

	exist := false
	for i := 0; i < int(len(ReqAgentConnsIndexSlice)); i++ {
		if ReqAgentConnsSlice[i] == nil {
			ReqAgentConnsIndexSlice[i] = idf
			ReqAgentConnsSlice[i] = conn
			exist = true
			break
		}
	}
	if !exist {
		ReqAgentConnsIndexSlice = append(ReqAgentConnsIndexSlice, idf)
		ReqAgentConnsSlice = append(ReqAgentConnsSlice, conn)
	}

}

func AddWorkerConn(worker_idf string, conn *net.TCPConn) {

	exist := false
	for i := 0; i < int(len(WorkerAgentConnsIndexSlice)); i++ {
		if WorkerAgentConnsSlice[i] == nil {
			WorkerAgentConnsIndexSlice[i] = worker_idf
			WorkerAgentConnsSlice[i] = conn
			exist = true
			break
		}
	}
	if !exist {
		WorkerAgentConnsIndexSlice = append(WorkerAgentConnsIndexSlice, worker_idf)
		WorkerAgentConnsSlice = append(WorkerAgentConnsSlice, conn)
	}

}

func DeleteReqConn(idf string) {

	for i := 0; i < int(len(ReqAgentConnsIndexSlice)); i++ {
		if ReqAgentConnsIndexSlice[i] == idf {
			ReqAgentConnsIndexSlice[i] = ""
			ReqAgentConnsSlice[i] = nil
		}
	}

}

func DeleteWorkerConn(worker_idf string) {

	for i := 0; i < int(len(WorkerAgentConnsIndexSlice)); i++ {
		if WorkerAgentConnsIndexSlice[i] == worker_idf {
			WorkerAgentConnsIndexSlice[i] = ""
			WorkerAgentConnsSlice[i] = nil
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Error(os.Stderr, "Fatal error: %s", err.Error())
	}
}
