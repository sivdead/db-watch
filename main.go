package main

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
	"db-watch/config"
	"db-watch/job"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/json-iterator/go"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"path"
	"github.com/pkg/errors"
)

const (
	cronStr    = "@every 3s"
	version    = "version: 0.0.1"
	connectURL = "%s:%s@(%s)/%s?charset=utf8"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	conf = make([]config.Config, 0)

	executor = cron.New()
)

func init() {
	configLocalFilesystemLogger("", "db-watch", 240*time.Hour, 4*time.Hour)
	initParameters()
}

func main() {

	for _, c := range conf {

		url := fmt.Sprintf(connectURL, c.Username, c.Password, c.Addr, c.Schema)
		//fmt.Println(url)
		engine, err := xorm.NewEngine("mysql", url)
		if err != nil {
			log.Panicf("error occurs while connecting to database %s: %s \n", c.Name, err.Error())
		}
		err = engine.Ping()
		if err != nil {
			log.Panicf("error occurs while ping database %s: %s", c.Name, err.Error())
		}
		log.Infof("database :%s, connection established!", c.Name)

		//create logger file
		{
			/*f, err := os.OpenFile(c.Name+".db.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Warnf("error occurs while creating database logger file: %s \n", err.Error())
			}*/
			//engine.ShowSQL(true)
			//engine.ShowExecTime(true)
			//engine.SetLogger(xorm.NewSimpleLogger2(f, "["+c.Name+"]", sysLog.Ldate|sysLog.Ltime))
			//engine.Logger().SetLevel(core.LOG_DEBUG)
		}
		if c.Cron == "" {
			c.Cron = cronStr
		}
		executor.AddJob(c.Cron, &job.Job{Engine: engine, Config: c})
	}
	executor.Entries()
	//create scheduler
	executor.Start()
	log.Infof("starting timeout scheduler %s ", time.Now().Format("2006-01-02 15:04:05"))

	//wait for quit command
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
	defer executor.Stop()
	defer log.Println("stop timeout scheduler")
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{})
	log.AddHook(lfHook)
}

func initParameters() {

	//args := os.Args
	//fmt.Println(args)
	v := flag.BoolP("version", "v", false, "show version")
	h := flag.BoolP("help", "h", false, "show help")

	configPath := flag.StringP("config-file", "c", "", "config file path")
	flag.Parse()

	if *h {
		flag.Usage()
		os.Exit(0)
	}

	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
	if *configPath != "" {
		data, err := ioutil.ReadFile(*configPath)
		if err != nil {
			fmt.Printf("error occurs while reading config file: %s \n", err.Error())
			os.Exit(1)
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			fmt.Printf("error occures while parsing config file: %s \n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("please specify config file path!")
		os.Exit(1)
	}
}
