package connector

import (
	"fmt"
	"gomore/area"
	"gomore/global"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"time"
	//"gomore/protocol"
	"gomore/lib/websocket"
	"gomore/worker"

	"gomore/golog"
	//"strings"
	//"io"
	"sync"
)

var Glock *sync.Mutex
var ConnMlock *sync.RWMutex
var ChannelMlock *sync.RWMutex
var SessionMlock *sync.RWMutex
var UserChannelsMlock *sync.RWMutex

func CloseWsConn(sid string) {

	_, conn_exist := global.SyncUserWebsocketConns.Get(sid)
	if conn_exist {
		global.SyncUserWebsocketConns.Delete(sid)
	}

}

func CloseConn(sid string) {

	_, conn_exist := global.SyncUserConns.Get(sid)
	if conn_exist {
		global.SyncUserConns.Delete(sid)
	}

}

func CloseSession(sid string) {

	_, session_exist := global.SyncUserSessions.Get(sid)
	if session_exist {
		global.SyncUserSessions.Delete(sid)
	}

}

func CloseUserChannel(sid string) {

	global.SyncUserJoinedChannels.Delete(sid)

}

func FreeConn(conn *net.TCPConn, sid string) {

	conn.Write([]byte{'E', 'O', 'F'})
	conn.Close()
	golog.Warn("Sid closing:", sid)
	CloseConn(sid)
	CloseSession(sid)
	CloseUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	area.UserUnSubscribeChannel(sid)
	golog.Info("UserConns length:", len(global.UserConns))

}

func FreeWsConn(ws *websocket.Conn, sid string) {

	ws.Write([]byte{'E', 'O', 'F'})
	ws.Close()
	golog.Warn("Sid closing:", sid)
	CloseWsConn(sid)
	CloseSession(sid)
	CloseUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	area.UserUnSubscribeChannel(sid)
	golog.Info("UserConns length:", len(global.UserConns))

}

func GetRandWorkerIdf() string {

	worker_idf := ""
	worker_len := int(len(worker.WorkerAgentConnsIndexSlice))
	rand_num := int(rand.Intn(worker_len))
	for i := 0; i < worker_len; i++ {
		if rand_num == i {
			worker_idf = worker.WorkerAgentConnsIndexSlice[i]
			break
		}
	}
	return worker_idf
}

func checkError(err error) {
	if err != nil {
		golog.Error(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func stat_kick() {

	timer := time.Tick(1000 * time.Millisecond)
	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		fmt.Println(time.Now().Unix(), " Connections: ", global.SumConnections, "  Qps: ", global.Qps)
	}
}

func user_kick(conn *net.TCPConn) {

	timer := time.Tick(5000 * time.Millisecond)
	for _ = range timer {
		ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }`, time.Now().Unix())
		go conn.Write([]byte(fmt.Sprintf("%s\r\n", ping)))
	}
}
