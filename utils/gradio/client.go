package gradio

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/erbieio/web2-bridge/config"
)

var nagativePrompts = "worst quality, normal quality, low quality, low res, blurry, text, watermark, logo, banner, extra digits, cropped, jpeg artifacts, signature, username, error, sketch ,duplicate, ugly, monochrome, horror, geometry, mutation, disgusting"

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
		FnIndex: 8,
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
		return nil, errors.New(toVedioResp.Status)
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

func Description2Prompts(desp string) (string, error) {
	params := DescriptionToPromptsReq{
		Data:    []string{desp},
		FnIndex: 3,
	}
	bodyJson, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", config.GetGradioConfig().Url+"api/predict", bytes.NewBuffer(bodyJson))
	if err != nil {
		return "", err
	}
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")

	toPromptResp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer toPromptResp.Body.Close()
	b, err := io.ReadAll(toPromptResp.Body)
	if err != nil {
		return "", err
	}
	if toPromptResp.StatusCode != 200 {
		return "", errors.New(string(b))
	}
	vBody := &DescriptionToPrompts{}
	err = json.Unmarshal(b, vBody)
	if err != nil {
		return "", err
	}
	if len(vBody.Data) == 0 {
		return "", errors.New("empty prompts generated")
	}
	return vBody.Data[0], nil

}

func Prompts2Image(prompts string) ([]byte, string, error) {
	seed := rand.Intn(1000000000000000000)
	params := PromptsToImageReq{
		Data:    []interface{}{prompts, nagativePrompts, 512, 512, 10, 50, seed},
		FnIndex: 5,
	}
	bodyJson, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", config.GetGradioConfig().Url+"api/predict", bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, "", err
	}
	client := &http.Client{}

	req.Header.Set("Content-Type", "application/json")

	toImageResp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer toImageResp.Body.Close()
	b, err := io.ReadAll(toImageResp.Body)
	if err != nil {
		return nil, "", err
	}
	if toImageResp.StatusCode != 200 {
		return nil, "", errors.New(string(b))
	}
	vBody := &PromptsToImage{}
	err = json.Unmarshal(b, vBody)
	if err != nil {
		return nil, "", err
	}
	if len(vBody.Data) == 0 {
		return nil, "", errors.New("empty image generated")
	}
	imageB64 := vBody.Data[0]
	prefix := "image/"
	suffix := ";base64"
	start := strings.Index(imageB64, prefix)
	if start < 0 {
		return nil, "", errors.New("unknown image format")
	}
	start += len(prefix)
	end := strings.Index(imageB64[start:], suffix)
	if end < 0 {
		return nil, "", errors.New("unknown image format")
	}
	ext := imageB64[start : start+end]
	imageBytes, err := base64.StdEncoding.DecodeString(imageB64[strings.IndexByte(imageB64, ',')+1:])

	return imageBytes, ext, err

}
