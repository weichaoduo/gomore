package admin

import (
	"fmt"
	"gomore/global"
	"net/http"
	"os"

	cpu "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

/**
 * 建立web server
 */
func HttpServer() {

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/admin/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	http.HandleFunc("/stats", statsTask)
	go func() {
		http.ListenAndServe(":"+global.Config.Admin.HttpPort, nil)
	}()
}

func statsTask(w http.ResponseWriter, req *http.Request) {
	fmt.Println("statsTask is running...")
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	cpuf, _ := cpu.Percent(1, true)
	fmt.Println("Percent:", cpuf)
	str := fmt.Sprintf(`{"conns":%d,"qps":%d,"cpu_per":%v,"mem_total":"%v","mem_free":"%v" , "mem_use_per":"%f"}`, global.SumConnections, global.Qps, cpuf, v.Total, v.Free, v.UsedPercent)
	w.Write([]byte(str))

}
