package gradio

type ImageToVedioReq struct {
	Data    []string `json:"data"`
	FnIndex int      `json:"fn_index"`
}
