package comfyui

import (
	"errors"
	"fmt"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/richinsley/comfy2go/client"
)

func Prompts2Image(prompts string) ([]byte, string, error) {
	clientaddr := config.GetComfyuiConfig().Host
	clientport := config.GetComfyuiConfig().Port
	workflow := "txt2img_workflow.json"

	// create a new ComgyGo client
	c := client.NewComfyClient(clientaddr, clientport, nil)

	// the ComgyGo client needs to be in an initialized state before
	// we can create and queue graphs
	if !c.IsInitialized() {
		logger.Logrus.Info(fmt.Sprintf("Initialize Client with ID: %s\n", c.ClientID()))
		err := c.Init()
		if err != nil {
			return nil, "", err
		}
	}

	// create a graph from the png file
	graph, _, err := c.NewGraphFromJsonFile(workflow)
	if err != nil {
		return nil, "", err
	}
	promptNode := graph.GetNodeById(6)
	promptNode.WidgetValues = []interface{}{
		prompts,
	}
	// queue the prompt and get the resulting image =
	item, err := c.QueuePrompt(graph)
	if err != nil {
		return nil, "", err
	}

	// continuously read messages from the QueuedItem until we get the "stopped" message type
	for continueLoop := true; continueLoop; {
		msg := <-item.Messages
		switch msg.Type {
		case "stopped":
			// if we were stopped for an exception, display the exception message
			qm := msg.ToPromptMessageStopped()
			if qm.Exception != nil {
				return nil, "", err
			}
			continueLoop = false
		case "data":
			qm := msg.ToPromptMessageData()
			for _, v := range qm.Images {
				// retrieve the image from ComfyUI
				img_data, err := c.GetImage(v)
				if err != nil {
					return nil, "", err
				}
				return *img_data, "png", err
			}
		}
	}
	return nil, "", errors.New("failed gen image")
}
