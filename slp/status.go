package slp

import (
	"encoding/json"
	"io"
)

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        interface{} `json:"description"`
	Favicon            string      `json:"favicon"`
	PreviewsChat       bool        `json:"previewsChat"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type Players struct {
	Max    int      `json:"max"`
	Online int      `json:"online"`
	Sample []Sample `json:"sample"`
}

type Sample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (sr StatusResponse) Write(w io.Writer) (n int, err error) {
	jsonResponseBytes, err := json.Marshal(sr)
	if err != nil {
		return
	}

	if n, err = WriteJsonBytes(w, jsonResponseBytes); err != nil {
		return
	}

	return
}

func (sr *StatusResponse) Read(r io.Reader) (err error) {
	if err = ReadStatusRequest(r); err != nil {
		return
	}

	return ReadJsonBytes(r, sr)
}
