package gradio_test

import (
	"flag"
	"io"
	"os"
	"testing"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/gradio"
	"github.com/erbieio/web2-bridge/utils/logger"
)

func TestImage2Vedio(t *testing.T) {
	configPath := flag.String("config_path", "../../", "config file")
	logicLogFile := flag.String("logic_log_file", "../../log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	config.LoadConf(*configPath)
	reader, err := gradio.Image2Vedio("https://www.erbiescan.io/ipfs/QmXr6A42GrMzCQE8i6GdsEVmmv6gcTV5WMs2hF174BaoDE")
	if err != nil {
		t.Error(err)
		return
	}
	vedioBytes, _ := io.ReadAll(reader)
	err = os.WriteFile("output.mp4", vedioBytes, 0644)
	if err != nil {
		t.Error(err)
	}
}

func TestDescription2Prompts(t *testing.T) {
	configPath := flag.String("config_path", "../../", "config file")
	logicLogFile := flag.String("logic_log_file", "../../log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	config.LoadConf(*configPath)
	prompts, err := gradio.Description2Prompts("一个女人坐在沙滩上")
	if err != nil {
		t.Error(err)
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(prompts)
}

func TestPrompts2Image(t *testing.T) {
	configPath := flag.String("config_path", "../../", "config file")
	logicLogFile := flag.String("logic_log_file", "../../log/bridge.log", "logic log file")
	flag.Parse()

	//init logic logger
	logger.Init(*logicLogFile)

	//load config
	config.LoadConf(*configPath)
	imageBytes, ext, err := gradio.Prompts2Image("A woman sitting on the beach, detailed realism, soft lighting, 4k, with a calm and peaceful expression on her face, her long hair flowing in the ocean breeze")
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
