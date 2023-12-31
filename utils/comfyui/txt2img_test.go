package comfyui_test

import (
	"flag"
	"os"
	"testing"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/comfyui"
	"github.com/erbieio/web2-bridge/utils/logger"
)

func TestPrompts2Image(t *testing.T) {
	configPath := flag.String("config_path", "../../", "config file")
	logicLogFile := flag.String("logic_log_file", "../../log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	config.LoadConf(*configPath)
	imageBytes, ext, err := comfyui.Prompts2Image("A boy")
	if err != nil {
		t.Error(err)
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	err = os.WriteFile("output."+ext, imageBytes, 0644)
	if err != nil {
		t.Error(err)
	}
}
