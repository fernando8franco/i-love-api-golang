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

func (ac ApiCredentials) Upload(server string, params UploadRequest) (UploadResponse, error) {
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
	req.Header.Set("Authorization", "Bearer "+ac.AuthToken)

	res, err := ac.APIClient.Do(req)
	if err != nil {
		return UploadResponse{}, fmt.Errorf("error sending request:\n%v", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return UploadResponse{}, handleError(res)
	}

	var response UploadResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return UploadResponse{}, fmt.Errorf("error decoding response body:\n%v", err)
	}

	return response, nil
}

func prepareLocalBody(params UploadRequest) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	formFieldName := "file"
	part, err := writer.CreateFormFile(formFieldName, params.FileName)
	if err != nil {
		return nil, "", fmt.Errorf("error creating form file:\n%v", err)
	}

	_, err = io.Copy(part, params.File)
	if err != nil {
		return nil, "", fmt.Errorf("error copying file data:\n%v", err)
	}

	writer.WriteField("task", params.TaskID)

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("error closing writer:\n%v", err)
	}

	return body, writer.FormDataContentType(), nil
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
		return nil, "", err
	}

	return bytes.NewReader(jsonBody), "application/json", nil
}
