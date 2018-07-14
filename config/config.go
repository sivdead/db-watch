package config

type Config struct {
	Name     string   `json:"name"`
	Addr     string   `json:"addr"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Schema   string   `json:"schema"`
	Command  []string `json:"command"`
	Timeout  int32    `json:"timeout"`
	Action   string   `json:"action"`
	Cron     string   `json:"cron"`
}

//timeout Action
const (
	Log  = "log"  //just log
	Kill = "kill" // log and kill the process
)
