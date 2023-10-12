package twitter_test

import (
	"flag"
	"testing"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/erbieio/web2-bridge/utils/twitter"
)

func TestGetTweetTimeLineMentions(t *testing.T) {
	configPath := flag.String("config_path", "../../", "config file")
	logicLogFile := flag.String("logic_log_file", "../../log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	config.LoadConf(*configPath)
	client, err := twitter.NewBearerTokenClient(config.GetTwitterConfig().Bearer)
	if err != nil {
		t.Error(err)
		return
	}
	metions, err := twitter.GetTweetTimeLineMentions(client, "VXNlcjoxNTgwMTQ2NTA3NDI0MzM3OTIw", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(metions)
}
