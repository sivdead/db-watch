package job

import (
	"github.com/go-xorm/xorm"
	"timeout-kill/model"
	"github.com/lunny/log"
	"timeout-kill/config"
	"strings"
	"fmt"
)

const (
	querySQL = "db = ? and command in (?) and time > ? "
	killSQL  = "kill query %d"
)

type Job struct {
	Engine     *xorm.Engine
	Config     config.Config
	commandMap map[string]bool
	initialized bool
}

func (job *Job) init() {
	job.commandMap = make(map[string]bool)
	for _, c := range job.Config.Command {
		job.commandMap[strings.ToLower(c)] = true
	}
	job.initialized = true
}

func (job *Job) Run() {

	if !job.initialized{
		job.init()
	}

	processList := make([]model.Process, 0)
	err := job.Engine.SQL("show processlist").Find(&processList)
	if err != nil {
		log.Errorf("error occurs while querying database: %s \n", err)
	}

	if len(processList) != 0 {
		count := 0
		for _, process := range processList {
			//fmt.Printf("schema: %s", job.Config.Schema)
			if int32(process.Time) > job.Config.Timeout &&
				process.Db == job.Config.Schema &&
				job.commandMap[strings.ToLower(process.Command)] {

				log.Warnf("database: %s, id: %d, process: %+v", job.Config.Name, process.Id, process)
				count ++

				if strings.Compare(config.Kill, strings.ToLower(job.Config.Action)) == 0 {
					_,err := job.Engine.Exec(fmt.Sprintf(killSQL,process.Id))
					if err != nil{
						log.Errorf("error occurs while killing process %d: %+v", process.Id, err)
					}
				}
			}
		}
		if count > 0 {
			log.Warnf("database: %s, find %d timeout query process \n", job.Config.Name, count)
		}
	}
}
