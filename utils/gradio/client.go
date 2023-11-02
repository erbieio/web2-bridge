package gradio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/erbieio/web2-bridge/config"
)

func Image2Vedio(path string) (io.Reader, error) {
	errMsg := fmt.Sprintf("gradio can't find the file or file is invalid %s", path)
	if !isValidURL(path) {
		return nil, fmt.Errorf("gradio found a invalid URL: %s", path)
	}
	resp, err := http.Get(path)
	defer resp.Body.Close()
	if err != nil {
		return nil, errors.New(errMsg)
	}
	cType := resp.Header.Get("Content-Type")
	cType = strings.Replace(cType, "/", ".", -1)
	ext := getExt(cType)
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		return nil, errors.New(errMsg)
	}
	imageBytes, err := io.ReadAll(resp.Body)
	imageBase64 := "data:image/" + ext + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)

	params := ImageToVedioReq{
		Data:    []string{imageBase64},
		FnIndex: 9,
	}
	bodyJson, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", config.GetGradioConfig().Url+"api/predict", bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	//req.Header.Set("User-Agent", "gradio_client_go/1.0")
	req.Header.Set("Content-Type", "application/json")

	toVedioResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer toVedioResp.Body.Close()
	b, err := io.ReadAll(toVedioResp.Body)
	if err != nil {
		return nil, err
	}
	if toVedioResp.StatusCode != 200 {
		return nil, errors.New(string(b))
	}
	vBody := &ImageToVedio{}
	err = json.Unmarshal(b, vBody)
	if err != nil {
		return nil, err
	}
	if len(vBody.Data) == 0 && len(vBody.Data[0]) == 0 {
		return nil, errors.New("vedio info not exist")
	}
	vedioSource, err := http.Get(config.GetGradioConfig().Url + "file=" + vBody.Data[0][0].Name)
	if err != nil {
		return nil, err
	}
	return vedioSource.Body, nil
}

func getExt(p string) string {
	s := strings.Split(p, ".")
	ext := s[len(s)-1]
	if ext == "jpeg" || ext == "jpg" || ext == "png" || ext == "gif" {
		return ext
	}
	return ""
}

func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true
}
