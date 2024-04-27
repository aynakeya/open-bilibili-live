package openblive

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
)

type DanmuHandler func(data DanmakuData)
type GiftHandler func(data GiftData)
type SuperChatHandler func(data SuperChatData)
type SuperChatDelHandler func(data SuperChatDelData)
type NewGuardHandler func(data NewGuardData)
type LikeHandler func(data LikeData)

type DisconnectCallback func(conn BLiveLongConnection)

type ConnErrorHandler func(err error)

type bLiveClientHandlers struct {
	Danmu        []DanmuHandler
	Gift         []GiftHandler
	SuperChat    []SuperChatHandler
	SuperChatDel []SuperChatDelHandler
	NewGuard     []NewGuardHandler
	Like         []LikeHandler
}

type BLiveLongConnection interface {
	Status() bool
	EstablishConnection(ctx context.Context) error
	CloseConnection() error
	OnDanmu(handlers ...DanmuHandler)
	OnGift(handlers ...GiftHandler)
	OnSuperChat(handlers ...SuperChatHandler)
	OnSuperChatDel(handlers ...SuperChatDelHandler)
	OnNewGuard(handlers ...NewGuardHandler)
	OnLike(handlers ...LikeHandler)
	OnDisconnect(callback DisconnectCallback)
	OnError(callback ConnErrorHandler)
}

type openBLiveLongConn struct {
	wssInfo WebSocketInfo

	HearbeatInterval time.Duration
	wsConn           *websocket.Conn
	status           bool

	Handlers bLiveClientHandlers
	cancel   context.CancelFunc

	diconnHandler DisconnectCallback
	errHandler    ConnErrorHandler
}

func (c *openBLiveLongConn) Status() bool {
	return c.status
}

func (c *openBLiveLongConn) OnDanmu(handlers ...DanmuHandler) {
	c.Handlers.Danmu = append(c.Handlers.Danmu, handlers...)
}

func (c *openBLiveLongConn) OnGift(handlers ...GiftHandler) {
	c.Handlers.Gift = append(c.Handlers.Gift, handlers...)
}

func (c *openBLiveLongConn) OnSuperChat(handlers ...SuperChatHandler) {
	c.Handlers.SuperChat = append(c.Handlers.SuperChat, handlers...)
}

func (c *openBLiveLongConn) OnSuperChatDel(handlers ...SuperChatDelHandler) {
	c.Handlers.SuperChatDel = append(c.Handlers.SuperChatDel, handlers...)
}

func (c *openBLiveLongConn) OnNewGuard(handlers ...NewGuardHandler) {
	c.Handlers.NewGuard = append(c.Handlers.NewGuard, handlers...)
}

func (c *openBLiveLongConn) OnLike(handlers ...LikeHandler) {
	c.Handlers.Like = append(c.Handlers.Like, handlers...)
}

func (c *openBLiveLongConn) OnDisconnect(callback DisconnectCallback) {
	c.diconnHandler = callback
}

func (c *openBLiveLongConn) OnError(callback ConnErrorHandler) {
	c.errHandler = callback
}

func (c *openBLiveLongConn) doErrCallBack(err error) {
	if c.errHandler != nil {
		c.errHandler(err)
	}
}

func NewOpenBLiveLongConn(
	wssInfo WebSocketInfo) BLiveLongConnection {
	return &openBLiveLongConn{
		wssInfo:          wssInfo,
		HearbeatInterval: 20 * time.Second,
		Handlers: bLiveClientHandlers{
			Danmu:        []DanmuHandler{},
			Gift:         []GiftHandler{},
			SuperChat:    []SuperChatHandler{},
			SuperChatDel: []SuperChatDelHandler{},
			NewGuard:     []NewGuardHandler{},
			Like:         []LikeHandler{},
		},
	}
}

func (c *openBLiveLongConn) sendAuth() error {
	err := c.wsConn.WriteMessage(
		websocket.BinaryMessage,
		MakeWSPacket(OpAuth, []byte(c.wssInfo.AuthBody)))
	if err != nil {
		return err
	}
	return nil
}

func (c *openBLiveLongConn) EstablishConnection(ctx context.Context) error {
	if len(c.wssInfo.WssLink) == 0 {
		return nil
	}
	ctx, c.cancel = context.WithCancel(ctx)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx,
		c.wssInfo.WssLink[0], nil,
	)
	if err != nil {
		return err
	}
	c.wsConn = conn

	err = c.sendAuth()
	if err != nil {
		return err
	}

	go c.eventLoop(ctx)
	c.status = true
	return nil
}

func (c *openBLiveLongConn) CloseConnection() error {
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	if c.wsConn != nil {
		_ = c.wsConn.Close()
		// dont reset wsConn, otherwise, it might trigger race condition when sending heartbeat
		// keep wsConn for sending heartbeat. if failed it will trigger an error instead of panic
		//c.wsConn = nil
	}
	c.status = false
	return nil
}

func (c *openBLiveLongConn) eventLoop(ctx context.Context) {
	ticker := time.NewTicker(c.HearbeatInterval)
	packetChan := make(chan WsPacket, 16)
	disconnedChan := make(chan int)
	go func() {
		for {
			messageType, message, err := c.wsConn.ReadMessage()
			if err != nil {
				c.doErrCallBack(err)
				break
			}
			if messageType != websocket.BinaryMessage {
				continue
			}
			packet, ok := ResolveWSPacket(message)
			if !ok {
				continue
			}
			if packet.Header.ProtocolVersion == 2 {
				if datas, err := ZlibDeCompress(packet.Data); err == nil {
					offset := 0
					for offset < len(datas) {
						subPacket, k := ResolveWSPacket(datas[offset:])
						if !k {
							break
						}
						packetChan <- subPacket
						offset += int(subPacket.Header.PacketLength)
					}
				}
			} else {
				packetChan <- packet
			}
		}
		// if for loop breaks, connection was broke
		disconnedChan <- 1
		c.status = false
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-disconnedChan:
			if c.diconnHandler != nil {
				c.diconnHandler(c)
			}
		case <-ticker.C:
			c.doErrCallBack(c.wsConn.WriteMessage(websocket.BinaryMessage, MakeWSPacket(OpHeartbeat, []byte("miao"))))
		case packet := <-packetChan:
			c.handleConnMessage(packet)
		}
	}
}

func (c *openBLiveLongConn) handleConnMessage(packet WsPacket) {
	switch int(packet.Header.Operation) {
	case OpSendMsg, OpSendMsgReply:
		c.handleCommand(packet.Data)
	default:
		return
	}
}

func (c *openBLiveLongConn) handleCommand(data []byte) {
	var cmdData CmdData
	if json.Unmarshal(data, &cmdData) != nil {
		return
	}
	switch cmdData.Cmd {
	case CmdDanmu:
		for _, handler := range c.Handlers.Danmu {
			handler(cmdData.ToDanmu())
		}
	case CmdGift:
		for _, handler := range c.Handlers.Gift {
			handler(cmdData.ToGift())
		}
	case CmdSuperChat:
		for _, handler := range c.Handlers.SuperChat {
			handler(cmdData.ToSuperChat())
		}
	case CmdSuperChatDel:
		for _, handler := range c.Handlers.SuperChatDel {
			handler(cmdData.ToSuperChatDel())
		}
	case CmdNewGuard:
		for _, handler := range c.Handlers.NewGuard {
			handler(cmdData.ToNewGuard())
		}
	case CmdLike:
		for _, handler := range c.Handlers.Like {
			handler(cmdData.ToLike())
		}
	default:
		return
	}
}
