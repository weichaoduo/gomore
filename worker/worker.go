//
//  Load-balancing broker.
// Use of Reactor, and other higher level functions.
//

package worker

import (
	log "github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
	//"github.com/satori/go.uuid"
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
	"gomore/global"
)

 

func workerTaskWithBufferio(index string, host string, port int) {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idf := fmt.Sprintf("%d%d", r.Intn(9999), rand.Intn(99999))

	fmt.Println(" worker_task tcpAddr:  \n", fmt.Sprintf(host+":%d", port))

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(host+":%d", port))
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	//defer conn.Close()
	time.Sleep(10 * time.Millisecond)
	worker_ready := fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_WORKER_CONNECT,  "", idf, "")
	fmt.Println("worker_ready:", worker_ready)
	conn.Write([]byte(worker_ready + "\n"))

	cmd := ""
	client_idf := ""
	worker_idf := ""
	task_data := ""

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		//fmt.Println( "worker_task recive 5:", msg  )
		if err != nil {
			fmt.Println( "worker_task ", " connection error: ", err.Error())
			conn.Close()
			return
		}
		if msg == "" {
			continue
		}
		
		//fmt.Println("worker_task from ", worker_idf, " receive str:", string(msg))

		ret_arr := strings.Split(msg, "||")

		if len(ret_arr) < 4 {
			//fmt.Println("workeDispath data length error!")
			continue
		}

		//fmt.Println("worke task arr :", ret_arr)
		cmd = ret_arr[0]
		client_idf = ret_arr[1]
		worker_idf = ret_arr[2]
		task_data = ret_arr[3]
        //fmt.Println("worke task cmd :", cmd )
		if cmd == global.DATA_REQ_MSG {

			worker_json, errjson := jason.NewObjectFromBytes([]byte(task_data))
			checkError(errjson)

			//  do some thing
			cmd, _ = worker_json.GetString("cmd")
			//fmt.Printf(" worker_task logic cmd: %s\n", cmd)

			json := fmt.Sprintf(`{"cmd":"%s","data":"%s"}`, cmd, client_idf )
			data := fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_WORKER_REPLY,  client_idf, worker_idf, json)
			//fmt.Println(" post : ", data)
			_, err = conn.Write([]byte(data + "\n"))
			checkError(err)
			//time.Sleep(10 * time.Millisecond)

		}
		if cmd == global.DATA_WORKER_CONNECT {

			fmt.Printf(" worker %s  ready!\n", worker_idf)

		}

	}
}

//
func Start() {

	worker_nbrs, err := global.ConfigJson.GetInt64("worker", "worker_num")
	global.CheckError(err)
	worker_port, err := global.ConfigJson.GetInt64("worker", "port")
	global.CheckError(err)
	
	worker_host, err := global.ConfigJson.GetString("worker", "host")
	global.CheckError(err)

	log.Info("Hub broker ready")
	for worker_nbr := 0; worker_nbr < int(worker_nbrs); worker_nbr++ {

		go workerTaskWithBufferio(fmt.Sprintf("%d", worker_nbr), worker_host, int(worker_port))
	}
	log.Info("Hub worker  ready")

}
