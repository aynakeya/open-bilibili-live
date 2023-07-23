package openblive

import "encoding/json"

type GameInfo struct {
	GameID string `json:"game_id"`
}

type WebSocketInfo struct {
	AuthBody string   `json:"auth_body"`
	WssLink  []string `json:"wss_link"` //
}

type AnchorInfo struct {
	RoomID int    `json:"room_id"`
	Uname  string `json:"uname"`
	Uface  string `json:"uface"`
	UID    int    `json:"uid"`
}

type AppStartResult struct {
	GameInfo      GameInfo      `json:"game_info"`
	WebSocketInfo WebSocketInfo `json:"websocket_info"`
	AnchorInfo    AnchorInfo    `json:"anchor_info"`
}

type CommonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}
