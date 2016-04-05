package connector

import (
	"bufio"
	"fmt"
	"gomore/area"
	"gomore/global"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
	//"gomore/protocol"
	"encoding/json"
	"gomore/worker"

	log "gomore/lib/Sirupsen/logrus"
	//"strings"
	//"io"
	//sync"
)

/**
 * 监听客户端连接
 */
func SocketConnector(ip string, port int) {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		log.Error("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化
	log.Debug("Game Connetor Server :", ip, port)

	listenAcceptTCP(listen)
}

/**
 *  处理客户端连接
 */
func listenAcceptTCP(listen *net.TCPListener) {

	//go stat_kick()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Error("AcceptTCP Exception::", err.Error())
			continue
		}

		//defer conn.Close()
		atomic.AddInt32(&global.SumConnections, 1)
		conn.SetNoDelay(false)
		max_conns := int32(global.Config.Connector.MaxConections)
		if max_conns > 0 && global.SumConnections > max_conns {
			conn.Write([]byte(global.ERROR_MAX_CONNECTIONS + "\n"))
			conn.Close()
			continue
		}

		// 校验ip地址
		conn.SetKeepAlive(true)
		log.Info("RemoteAddr:", conn.RemoteAddr().String())

		//remoteAddr :=conn.RemoteAddr()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))

		req_host := global.Config.WorkerAgent.Host
		req_port := string(global.Config.WorkerAgent.Port)
		fmt.Println("tcpAddr:", req_host+":"+req_port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", req_host+":"+req_port)
		if err != nil {
			fmt.Println("req_conn tcpAddr :", err.Error())
			return
		}

		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		//defer req_conn.Close()
		if err != nil {
			fmt.Println("req_conn net.DialTCP :", err.Error())
			return
		}
		worker_idf := GetRandWorkerIdf()

		fmt.Println("RemoteAddr:", conn.RemoteAddr().String(), "sid:", sid, " worker_idf:", worker_idf)

		// 接收worker返回的数据
		go ReqWorkerAgentWithBufferio(conn, req_conn, sid, worker_idf)

		go handleConnWithBufferio(conn, req_conn, sid, worker_idf)
		// 发送数据给worker

	} //end for {

}

func ReqWorkerAgentWithBufferio(conn *net.TCPConn, req_conn *net.TCPConn, sid string, worker_idf string) {

	worker.AddReqConn(sid, req_conn)
	area.ConnRegister(conn, sid)
	//req_ready := fmt.Sprintf( `{"cmd":"req.connect", "client_idf":"%s" }`, sid    )
	req_ready := fmt.Sprintf(`%s||%s||%s||`, global.DATA_REQ_CONNECT, sid, worker_idf)

	req_conn.Write([]byte(req_ready + "\n"))

}

func handleConnWithBufferio(conn *net.TCPConn, req_conn *net.TCPConn, sid string, worker_idf string) {

	// 发包频率判断
	range_count := 1
	limit_date := global.Config.Connector.MaxPacketRate
	var now int64
	var start_time int64
	var range_times int64
	start_time = time.Now().Unix()
	range_times = int64(global.Config.Connector.MaxPacketRateUnit)

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	for {
		if !global.Config.Enable {
			conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		// 区间范围内的计数
		if limit_date > 0 {
			now = time.Now().Unix()
			if (now - start_time) <= range_times {
				range_count++
			} else {
				start_time = now
				range_count = 1
			}
			// 判断发包频率是否超过限制
			if range_count > limit_date {
				conn.Write([]byte(global.ERROR_PACKET_RATES + "\n"))
				conn.Close()
				break
			}
		}

		msg, err := reader.ReadString('\n')
		//fmt.Println(  "handleConn ReadString: ", string(msg) )
		if err != nil {
			FreeConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		if msg == "" {
			continue
		}
		go func(sid string, msg string, req_conn *net.TCPConn) {

			// fmt.Println(conn.RemoteAddr().String(), "receive str:", string(msg) )
			worker_idf := GetRandWorkerIdf()
			worker_data := fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_REQ_MSG, sid, worker_idf, msg)
			//fmt.Println("req push worker_data 2:", worker_data)
			req_conn.Write([]byte(worker_data + "\n"))

		}(sid, msg, req_conn)

	}

}

func handleConnWithJson(conn *net.TCPConn, req_conn *net.TCPConn, sid string, worker_idf string) {

	d := json.NewDecoder(conn)
	for {
		if global.Config.Enable {
			conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		var msg interface{}
		err := d.Decode(&msg)
		if err != nil {
			log.Info(conn.RemoteAddr().String(), " connection error: ", err.Error(), " , sid:", sid)
			FreeConn(conn, sid)
			break
		}

		//log.Info(conn.RemoteAddr().String(), "receive data:", msg)
		json_encode, err_encode := json.Marshal(msg)
		if err_encode != nil {
			log.Error("json.Marshal error:", err_encode.Error())
			conn.Write([]byte{'E', 'O', 'F'})
			conn.Close()
			atomic.AddInt32(&global.SumConnections, -1)
			break
		}
		str := string(json_encode)
		log.Info(conn.RemoteAddr().String(), "receive str:", str)

		if str != "" {

			worker_data := fmt.Sprintf(`req.msg||%s||%s||%s`, sid, worker_idf, str)
			_, err_req := req_conn.Write([]byte(worker_data))
			if err_req != nil {
				log.Error(" req_conn.Write  error:", err_req.Error())

			}

		}

	}

}

func handleConnWithProtobuf(conn *net.TCPConn, req_conn *net.TCPConn, sid string, worker_idf string) {

}
