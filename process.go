package iloveapigolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProcessRequest struct {
	Task  string  `json:"task"`
	Tool  string  `json:"tool"`
	Files []Files `json:"files"`
	Meta  Meta    `json:"meta"`
	CompressOptions
}

type Files struct {
	ServerFileName string `json:"server_filename"`
	FileName       string `json:"filename"`
}

type Meta struct {
	Title        string `json:"Title"`
	Author       string `json:"Author"`
	Subject      string `json:"Subject"`
	Keywords     string `json:"Keywords"`
	Creator      string `json:"Creator"`
	Producer     string `json:"Producer"`
	CreationDate string `json:"CreationDate"`
	ModDate      string `json:"ModDate"`
	Trapped      string `json:"Trapped"`
}

type CompressOptions struct {
	CompressionLevel string `json:"compression_level"`
}

type ProcessResponse struct {
	DownloadFilename string `json:"download_filename"`
	Filesize         int    `json:"filesize"`
	OutputFilesize   int    `json:"output_filesize"`
	OutputFilenumber int    `json:"output_filenumber"`
	OutputExtensions string `json:"output_extensions"`
	Timer            string `json:"timer"`
	Status           string `json:"status"`
}

func (c *Client) Process(server string, params ProcessRequest) (ProcessResponse, error) {
	processUrl := fmt.Sprintf(processURL, server)

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return ProcessResponse{}, fmt.Errorf("error encoding request:\n%v", err)
	}
	body := bytes.NewReader(jsonBody)

	req, err := http.NewRequest("POST", processUrl, body)
	if err != nil {
		return ProcessResponse{}, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return ProcessResponse{}, fmt.Errorf("error sending request:\n%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return ProcessResponse{}, handleError(res)
	}

	var response ProcessResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return ProcessResponse{}, fmt.Errorf("error decoding response:\n%v", err)
	}

	return response, nil
}
