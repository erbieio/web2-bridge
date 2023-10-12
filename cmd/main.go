package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/internal/bot"
	"github.com/erbieio/web2-bridge/internal/chain"
	_ "github.com/erbieio/web2-bridge/utils/db/mysql"
	"github.com/erbieio/web2-bridge/utils/logger"

	"github.com/urfave/cli"
)

func main() {
	local := []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "",
			Action: func(cctx *cli.Context) error {
				run(cctx)
				return nil
			},
		},
	}
	app := &cli.App{
		Name:  "erbio web2 brdige server",
		Usage: "erbio web2 brdige server",

		Commands: local,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cctx *cli.Context) {
	configPath := flag.String("config_path", "./", "config file")
	logicLogFile := flag.String("logic_log_file", "./log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	err := config.LoadConf(*configPath)
	if err != nil {
		log.Fatal("load config failed:", err)
	}
	serverConf := config.GetServerConfig()
	if serverConf.LogOutStdout() {
		logger.Logrus.Out = os.Stdout
	}

	//set log level
	logger.SetLogLevel(serverConf.RunMode)

	/* 	db := mysql.GetDB()
	   	if db == nil {
	   		logger.Logrus.Error("init db failed")
	   		return
	   	}

	   	err = redis.InitRedis()
	   	if err != nil {
	   		logger.Logrus.Error("init redis failed")
	   		return
	   	} */

	botFactory := bot.GetFacotory()
	botFactory.Register(&bot.DiscordBot{Handler: chain.MessageHandler})
	botFactory.Do()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logrus.Info("Server exiting")
}
