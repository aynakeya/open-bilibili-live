package openblive

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"math/rand"
	"time"
)

var ApiBase = "https://live-open.biliapi.com"

const (
	UrlAppStart          = "/v2/app/start"
	UrlAppStop           = "/v2/app/end"
	UrlAppHeartBeat      = "/v2/app/heartbeat"
	UrlAppBatchHeartBeat = "/v2/app/batchHeartbeat"
)

type apiHeader struct {
	TimeStamp        string `json:"x-bili-timestamp"`
	SignatureMethod  string `json:"x-bili-signature-method"`
	SignatureNonce   string `json:"x-bili-signature-nonce"`
	AccessKey        string `json:"x-bili-accesskeyid"`
	SignatureVersion string `json:"x-bili-signature-version"`
	ContentMD5       string `json:"x-bili-content-md5"`
}

var apiHeaderSeq = []string{"x-bili-accesskeyid", "x-bili-content-md5", "x-bili-signature-method", "x-bili-signature-nonce", "x-bili-signature-version", "x-bili-timestamp"}

func (h *apiHeader) ToMap() map[string]string {
	return map[string]string{
		"x-bili-timestamp":         h.TimeStamp,
		"x-bili-signature-method":  h.SignatureMethod,
		"x-bili-signature-nonce":   h.SignatureNonce,
		"x-bili-accesskeyid":       h.AccessKey,
		"x-bili-signature-version": h.SignatureVersion,
		"x-bili-content-md5":       h.ContentMD5,
	}
}

func (h *apiHeader) ToHeaderStr() string {
	httpHeader := h.ToMap()
	headerStr := ""
	for _, val := range apiHeaderSeq {
		headerStr += val + ":" + httpHeader[val] + "\n"
	}
	return headerStr[:len(headerStr)-1]
}

type IApiClient interface {
	AppStart(code string, appId int64) (*AppStartResult, *PublicError)
	AppEnd(appId int64, gameId string) *PublicError
	HearBeat(gameId string) *PublicError
}

type ApiClient struct {
	client       *resty.Client
	accessKey    string
	accessSecret string
}

func NewApiClient(accessKey string, accessSecret string) *ApiClient {
	c := &ApiClient{
		client:       resty.New(),
		accessKey:    accessKey,
		accessSecret: accessSecret,
	}
	c.client.SetBaseURL(ApiBase)
	return c
}

func (c *ApiClient) post(url string, data map[string]interface{}) (*resty.Response, error) {
	var err error
	content, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	hash := md5.Sum(content)
	nounce := rand.Int63() + time.Now().Unix()
	header := apiHeader{
		TimeStamp:        cast.ToString(time.Now().Unix()),
		SignatureMethod:  "HMAC-SHA256",
		SignatureNonce:   cast.ToString(nounce),
		AccessKey:        c.accessKey,
		SignatureVersion: "1.0",
		ContentMD5:       hex.EncodeToString(hash[:]),
	}

	httpHeader := header.ToMap()

	authEncoder := hmac.New(sha256.New, []byte(c.accessSecret))
	authEncoder.Write([]byte(header.ToHeaderStr()))
	httpHeader["Authorization"] = hex.EncodeToString(authEncoder.Sum(nil))
	httpHeader["Content-Type"] = "application/json"
	httpHeader["Accept"] = "application/json"
	return c.client.R().
		SetHeaders(httpHeader).
		SetBody(data).
		Post(url)
}

func parseResponse[T any](data []byte) (*T, *PublicError) {
	var commonResp CommonResponse = CommonResponse{Code: -1}
	var err error
	err = json.Unmarshal(data, &commonResp)
	if err != nil {
		return nil, ErrUnknown.WithDetail(err)
	}
	if commonResp.Code != 0 {
		return nil, GetErrorFromCode(commonResp.Code)
	}
	var result T
	err = json.Unmarshal(commonResp.Data, &result)
	if err == nil {
		return &result, nil
	}
	return nil, ErrUnknown.WithDetail(err)
}

func (c *ApiClient) AppStart(code string, appId int64) (*AppStartResult, *PublicError) {
	resp, err := c.post(UrlAppStart, map[string]interface{}{
		"code":   code,
		"app_id": appId,
	})
	if err != nil {
		return nil, ErrUnknown.WithDetail(err)
	}
	return parseResponse[AppStartResult](resp.Body())
}

func (c *ApiClient) AppEnd(appId int64, gameId string) *PublicError {
	resp, err := c.post(UrlAppStop, map[string]interface{}{
		"app_id":  appId,
		"game_id": gameId,
	})
	if err != nil {
		return ErrUnknown.WithDetail(err)
	}
	_, pe := parseResponse[map[string]interface{}](resp.Body())
	return pe
}

func (c *ApiClient) HearBeat(gameId string) *PublicError {
	resp, err := c.post(UrlAppHeartBeat, map[string]interface{}{
		"game_id": gameId,
	})
	if err != nil {
		return ErrUnknown.WithDetail(err)
	}
	_, pe := parseResponse[map[string]interface{}](resp.Body())
	return pe
}
