package job

import (
	"github.com/go-xorm/xorm"
	"github.com/json-iterator/go"
	"timeout-kill/model"
	"github.com/lunny/log"
	"timeout-kill/config"
	"strings"
	"github.com/go-xorm/builder"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	querySQL = "db = ? and command in (?) and time > ? "
	killSQL  = "kill query ?"
)

type Job struct {
	Engine *xorm.Engine
	Config *config.Config
}

func (job Job) Run() {
	//strings.Join(job.Config.Command, ",")
	processList := make([]model.Process, 0)
	//err := job.Engine.Where(querySQL, job.Config.Schema, "Query", 3).Find(&processList)
	err := job.Engine.Select("*").
		Where("db = ? and time > ? ", job.Config.Schema, job.Config.Timeout).
		And(builder.In("command", job.Config.Command)).
		Find(&processList)
	if err != nil {
		log.Errorf("error occurs while querying database: %s \n", err)
	}

	if len(processList) != 0 {
		log.Warnf("database: %s, find %d timeout query process: \n", job.Config.Name, len(processList))
		if strings.Compare(config.Kill, strings.ToLower(job.Config.Action)) != 0 {
			return
		}
		for index, process := range processList {
			str, err := json.Marshal(&process)
			if err != nil {
				log.Errorf("err occurs while marshaling json: %s", err.Error())
			}
			log.Warnf("database: %s, index: %d, process: %s",job.Config.Name, index, string(str))
			job.Engine.Exec(killSQL, process.Id)
		}
	}
}
