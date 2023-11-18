package ingest

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.k6.io/k6/js/modules"
	"io"
	"math/rand"
	"net/http"
	url2 "net/url"
	"strings"
)

const (
	encoding = "JPEG_LOSSY"
)

var (
	client = http.DefaultClient
)

type BasicProperties struct {
	OriginalName string `json:"OriginalName"`
	Width        int32  `json:"Width"`
	Height       int32  `json:"Height"`

	MicronPerPixelWidth   float64 `json:"MicronPerPixelWidth"`
	MicronPerPixelHeight  float64 `json:"MicronPerPixelHeight"`
	Size                  int32   `json:"Size"`
	MagnificationLevelMax int32   `json:"MagnificationLevelMax"`

	TileXMinNative int `json:"TileXMinNative"`
	TileXMaxNative int `json:"TileXMaxNative"`
	TileYMinNative int `json:"TileYMinNative"`
	TileYMaxNative int `json:"TileYMaxNative"`
}

func init() {
	modules.Register("k6/x/sisi", &SISI{})
}

type SISI struct{}

func (*SISI) GetRandomTile(token string, props BasicProperties) error {

	x := int(rand.Float64()*float64(props.TileXMaxNative-props.TileXMinNative)) + props.TileXMinNative
	y := int(rand.Float64()*float64(props.TileYMaxNative-props.TileYMinNative)) + props.TileYMinNative
	_, err := getTile(token, x, y, x, y, 0)
	return err
}

func (*SISI) GetBasicProperties(slideToken string) (*BasicProperties, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:5120/api/simpleslideinterface/v1/slide/%s/base_properties", slideToken))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Wrong status code : %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response BasicProperties
	err = json.Unmarshal(body, &response)
	return &response, err
}

func (*SISI) GetSlideToken(url string, rootPath string, path string) (string, error) {
	u := fmt.Sprintf(url, url2.PathEscape(fmt.Sprintf(rootPath, path)))
	resp, err := http.Post(u, "application/json", nil)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return resp.Status, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(body), "\"", ""), nil
}

func (*SISI) GetTile(slideToken string, x1, y1, magnification int) ([]byte, error) {
	return getTile(slideToken, x1, y1, x1, y1, magnification)
}

func getTile(slideToken string, x1, y1, x2, y2, magnification int) ([]byte, error) {
	_url := fmt.Sprintf("http://localhost:5120/api/simpleslideinterface/v1/slide/%s/tile", slideToken)
	req, err := http.NewRequest("GET", _url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Set("encoding", encoding)
	q.Set("x1", fmt.Sprintf("%d", x1))
	q.Set("y1", fmt.Sprintf("%d", y1))
	q.Set("x2", fmt.Sprintf("%d", x2))
	q.Set("y2", fmt.Sprintf("%d", y2))
	q.Set("magnification", fmt.Sprintf("%d", magnification))

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Wrong status code : %d", resp.StatusCode))
	}
	return io.ReadAll(resp.Body)
}
