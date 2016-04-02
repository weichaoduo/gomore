package global

import (
    "fmt"
   "github.com/BurntSushi/toml"  
)

type configType struct {
  Name       string
  Enable     bool
  Status     string
  Version    string
  Loglevel   string
  RpcType    string 
  PackType   string 
  Connector  connector  `toml:"connector"`
  Object  object  `toml:"object"`
  Worker         worker  `toml:"worker"`
  WorkerAgent    workerAgent `toml:"worker_agent"`
  Hub            hub `toml:"hub"`
  Area           area `toml:"area"`
  
}

type connector struct {
    
    WebsocketPort 	  int `toml:"websocket_port"` 
    SocketPort    	  int  `toml:"socket_port"` 
    MaxConections 	  int  `toml:"max_conections"` 
    MaxConntionsIp    int    `toml:"max_conntions_ip"` 
    MaxPacketRate	  int      `toml:"max_packet_rate"` 
    MaxPacketRateUnit int `toml:"max_packet_rate_unit"`
}


type object struct {
        
    DataType 	  string  `toml:"data_type"`   
	RedisHost 	  string  `toml:"redis_host"` 
	RedisPort 	  string     `toml:"redis_port"` 
	RedisPassword 	  string  `toml:"redis_password"` 
	MonogoHost    string  `toml:"monogo_host"` 
	MonogoPort    int     `toml:"3306"` 
}

type worker struct {
    WorkerLanguage  string      `toml:"worker_language"` 
    PhpBinPath 	    string      `toml:"php_bin_path"` 
	WorkerNum 		int         `toml:"worker_num"` 
	AgentHost 		string      `toml:"agent_host"` 
	AgentPort 		int         `toml:"agent_port"` 
	HubHost 		string      `toml:"hub_host"` 
	HubPort 		int         `toml:"hub_port"` 
	
	
}

type	workerAgent struct {
	Host	string      
	Port	int    
}


type   hub struct {
    	Hub_host string
    	Hub_port int
 }

type	area struct {
		Init_area []string
		
}

var Config configType

func InitConfig(){
    
    
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		fmt.Println( "toml.DecodeFile error:", err )
		return
	}

}