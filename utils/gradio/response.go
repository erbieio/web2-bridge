package gradio

type ImageToVedio struct {
	AverageDuration     float64      `json:"average_duration"`
	Data                [][]FileBody `json:"data"`
	Duration            float64      `json:"duration"`
	IsGenerating        bool         `json:"is_generating"`
	ModelscopeRequestID string       `json:"modelscope_request_id"`
}

type FileBody struct {
	Data     interface{} `json:"data"`
	IsFile   bool        `json:"is_file"`
	Name     string      `json:"name"`
	OrigName string      `json:"orig_name"`
}
