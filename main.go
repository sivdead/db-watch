package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
	"io"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/robfig/cron"
	"timeout-kill/job"
	"fmt"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"github.com/json-iterator/go"
	"timeout-kill/config"
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
	initLogger()
	initParameters()
}

func main() {

	for _, c := range conf {

		url := fmt.Sprintf(connectURL, c.Username, c.Password, c.Addr, "information_schema")
		engine, err := xorm.NewEngine("mysql", url)
		if err != nil {
			log.Errorf("error occurs while connecting to database %s: %s \n", c.Name, err.Error())
			os.Exit(1)
		}
		defer engine.Close()
		log.Infof("database %s connection establish success!", c.Name)

		//create logger file
		{
			f, err := os.OpenFile(c.Name+".db.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Errorf("error occurs while creating database logger file: %s \n", err.Error())
			}
			defer f.Close()
			engine.ShowSQL(true)
			engine.ShowExecTime(true)
			engine.SetLogger(xorm.NewSimpleLogger(f))
			engine.Logger().SetLevel(core.LOG_DEBUG)
		}

		executor.AddJob(cronStr, job.Job{engine, &c})
	}
	//create scheduler
	executor.Run()
	log.Infof("starting timeout scheduler %s ", time.Now().Format("2006-01-02 15:04:05"))

	//wait for quit command
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
	defer executor.Stop()
	defer log.Println("stop timeout scheduler")
}

func initLogger() {
	now := time.Now().Format("2006-01-02")
	f, err := os.OpenFile(now+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("error occuors while creating sys logger file:  %s\n", err.Error())
	}
	//defer f.Close()
	writers := []io.Writer{
		os.Stdout,
		f,
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(io.MultiWriter(writers...))
}

func initParameters() {

	//args := os.Args
	//fmt.Println(args)
	v := flag.BoolP("version", "v", false, "show version")
	h := flag.BoolP("help", "h", false, "show help")

	configPath := flag.StringP("config-file", "c", "", "config file path")
	flag.Parse()

	if *h{
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
