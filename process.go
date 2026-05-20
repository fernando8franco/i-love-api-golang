package iloveapigolang

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProcessParams struct {
	Server string
	Task   string
	Tool   string
	Files  []Files
	Meta
	CompressOptions
}

type Files struct {
	ServerFileName string
	FileName       string
}

type Meta struct {
	Title        string
	Author       string
	Subject      string
	Keywords     string
	Creator      string
	Producer     string
	CreationDate string
	ModDate      string
	Trapped      string
}

type CompressOptions struct {
	CompressionLevel string
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

func (c *Client) Process(ctx context.Context, params ProcessParams) (ProcessResponse, error) {
	processUrl := fmt.Sprintf(processURL, params.Server)

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return ProcessResponse{}, fmt.Errorf("error encoding request:\n%v", err)
	}
	body := bytes.NewReader(jsonBody)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		processUrl,
		body,
	)
	if err != nil {
		return ProcessResponse{}, fmt.Errorf("error creating request:\n%v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.GetToken())

	res, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return ProcessResponse{}, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		}
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
