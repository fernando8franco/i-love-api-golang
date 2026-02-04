package iloveapigolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type UploadRequest struct {
	TaskID       string
	File         io.Reader
	FileName     string
	CloudFileURL string
}

type UploadResponse struct {
	ServerFilename string `json:"server_filename"`
}

func (c *Client) Upload(server string, params UploadRequest) (UploadResponse, error) {
	url := fmt.Sprintf(uploadURL, server)

	var body io.Reader
	var contentType string
	var err error

	if params.CloudFileURL != "" {
		body, contentType, err = prepareCloudBody(params)
	} else {
		body, contentType, err = prepareLocalBody(params)
	}

	if err != nil {
		return UploadResponse{}, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("error sending request:\n%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return UploadResponse{}, handleError(res)
	}

	var response UploadResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return UploadResponse{}, fmt.Errorf("error decoding response:\n%v", err)
	}

	return response, nil
}

func prepareLocalBody(params UploadRequest) (io.Reader, string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		part, err := writer.CreateFormFile("file", params.FileName)
		if err != nil {
			pw.CloseWithError(err)
			return
		}

		if _, err := io.Copy(part, params.File); err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	return pr, writer.FormDataContentType(), nil
}

func prepareCloudBody(params UploadRequest) (io.Reader, string, error) {
	body := struct {
		Task      string `json:"Task"`
		CloudFile string `json:"CloudFile"`
	}{
		Task:      params.TaskID,
		CloudFile: params.CloudFileURL,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, "", fmt.Errorf("error encoding request:\n%v", err)
	}

	return bytes.NewReader(jsonBody), "application/json", nil
}
