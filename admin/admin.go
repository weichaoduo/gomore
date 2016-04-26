package admin

import (
	"encoding/json"
	"fmt"
	"gomore/global"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	cpu "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type MongoLog struct {
	Name    string
	Level   string
	File    string
	Line    int
	Message string
	Time    int
}

type Logs struct {
	All []MongoLog
}

/**
 * 建立web server
 */
func HttpServer() {

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/admin/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	http.HandleFunc("/stats", statsTask)
	http.HandleFunc("/lastlogs", lastLogsTask)
	go func() {
		http.ListenAndServe(":"+global.Config.Admin.HttpPort, nil)
	}()
}

func statsTask(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("statsTask is running...")
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	cpuf, _ := cpu.Percent(1, true)
	//fmt.Println("Percent:", cpuf)

	now := time.Now().Format("15:04:05")
	str := fmt.Sprintf(`{"time":"%s","conns":%d,"qps":%d,"cpu_per":%v,"mem_total":"%v","mem_free":"%v" , "mem_use_per":"%f"}`,
		now, global.SumConnections, global.Qps, cpuf, v.Total, v.Free, v.UsedPercent)
	w.Write([]byte(str))

}

func lastLogsTask(w http.ResponseWriter, req *http.Request) {
	fmt.Println("lastLogsTask is running...")

	session, err := mgo.Dial(global.Config.Log.MongodbHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("gomore") //数据库名称
	collection := db.C("logs")

	ms := []MongoLog{}
	err = collection.Find(bson.M{}).All(&ms)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("All:", ms)
	if b, err := json.Marshal(ms); err == nil {
		fmt.Println("================struct 到json str==")
		fmt.Println(string(b))
		w.Write(b)
	} else {
		w.Write([]byte(`[]`))
	}

}
