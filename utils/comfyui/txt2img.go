package comfyui

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/richinsley/comfy2go/client"
)

func Prompts2Image(c *client.ComfyClient, prompts string) ([]byte, string, error) {

	// the ComgyGo client needs to be in an initialized state before
	// we can create and queue graphs
	if !c.IsInitialized() {
		logger.Logrus.Info(fmt.Sprintf("Initialize Client with ID: %s\n", c.ClientID()))
		err := c.Init()
		if err != nil {
			return nil, "", err
		}
	}

	graph, _, err := c.NewGraphFromJsonReader(bytes.NewBufferString(workflow))
	if err != nil {
		return nil, "", err
	}
	promptNode := graph.GetNodeById(6)
	promptNode.WidgetValues = []interface{}{
		prompts,
	}
	sampleNode := graph.GetNodeById(13)
	sampleNode.WidgetValues = []interface{}{
		true,
		rand.Uint32(),
		"randomize",
		1,
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
