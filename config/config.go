package config

type Config struct {
	Name string `json:"name"`
	Addr string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema string `json:"schema"`
	Command []string `json:"command"`
	Timeout int32 `json:"timeout"`
	Action string `json:"action"`
}

//timeout Action
const (
	Log = "log"
	Kill = "kill"
)