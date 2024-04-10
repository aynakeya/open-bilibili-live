package openblive

import (
	"time"
)

type BLiveApp struct {
	AppID        int64
	AccessKey    string
	AccessSecret string
}

func (app *BLiveApp) CreateClient(code string) *BLiveClient {
	return &BLiveClient{
		AppID:            app.AppID,
		Code:             code,
		apiclient:        NewApiClient(app.AccessKey, app.AccessSecret),
		HearbeatInterval: 20 * time.Second,
	}
}

type BLiveClient struct {
	apiclient IApiClient
	AppID     int64
	Code      string // 主播身份码

	AppInfo *AppStartResult

	HearbeatInterval time.Duration

	running bool

	longConn BLiveLongConnection
}

func NewBliveClient(appID int64, code string, client IApiClient) *BLiveClient {
	return &BLiveClient{
		AppID:            appID,
		Code:             code,
		apiclient:        client,
		HearbeatInterval: 20 * time.Second,
	}
}

func (c *BLiveClient) Status() bool {
	return c.running
}

func (c *BLiveClient) Start() error {
	if c.running {
		return nil
	}
	resp, err := c.apiclient.AppStart(c.Code, c.AppID)
	if err != nil {
		return err
	}
	c.AppInfo = resp
	c.running = true
	// if GameID is empty. there is no need to send heartbeat
	if c.AppInfo.GameInfo.GameID == "" {
		return nil
	}
	go func() {
		for c.running {
			_ = c.apiclient.HearBeat(c.AppInfo.GameInfo.GameID)
			time.Sleep(c.HearbeatInterval)
		}
	}()
	return nil
}

func (c *BLiveClient) End() error {
	if !c.running {
		return nil
	}
	e := c.apiclient.AppEnd(c.AppID, c.AppInfo.GameInfo.GameID)
	c.AppInfo = nil
	c.running = false
	// public error to error
	if e == nil {
		return nil
	}
	return e
}

func (c *BLiveClient) GetLongConn() BLiveLongConnection {
	if c.longConn == nil {
		if c.AppInfo == nil {
			return nil
		}
		c.longConn = NewOpenBLiveLongConn(c.AppInfo.WebSocketInfo)
	}
	return c.longConn
}

func (c *BLiveClient) SetLongConn(conn BLiveLongConnection) {
	c.longConn = conn
}
