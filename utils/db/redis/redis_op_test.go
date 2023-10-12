package redis

import (
	"flag"
	"fmt"
	"testing"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/logger"

	"github.com/sirupsen/logrus"
)

func TestScanAndDelete(t *testing.T) {
	configPath := flag.String("config_path", "./", "config file")
	logicLogFile := flag.String("logic_log_file", "./log/sdk.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	err := config.LoadConf(*configPath)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"ErrMsg": err.Error()}).Error("load config failed")
		return
	}

	err = InitRedis()
	if err != nil {
		logger.Logrus.Error("init redis failed")
		return
	}

	matckkey := "nft.treasure.address" + "*"
	keys, err := Scan(matckkey)
	if err != nil {
		fmt.Printf("DeleteFTTreasureCache scan failed, %v", err)
		return
	}

	err = BatchDelete(keys)
	if err != nil {
		fmt.Printf("DeleteFTTreasureCache batch del failed, %v", err)
		return
	}
}
