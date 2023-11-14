package gradio

type ImageToVedioReq struct {
	Data    []string `json:"data"`
	FnIndex int      `json:"fn_index"`
}

type DescriptionToPromptsReq struct {
	Data    []string `json:"data"`
	FnIndex int      `json:"fn_index"`
}

type PromptsToImageReq struct {
	Data    []interface{} `json:"data"`
	FnIndex int           `json:"fn_index"`
}
