package openblive

import "encoding/json"

const (
	CmdDanmu        string = "LIVE_OPEN_PLATFORM_DM"
	CmdGift                = "LIVE_OPEN_PLATFORM_SEND_GIFT"
	CmdSuperChat           = "LIVE_OPEN_PLATFORM_SUPER_CHAT"
	CmdSuperChatDel        = "LIVE_OPEN_PLATFORM_SUPER_CHAT_DEL"
	CmdNewGuard            = "LIVE_OPEN_PLATFORM_GUARD"
	CmdLike                = "LIVE_OPEN_PLATFORM_LIKE"
)

type CmdData struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
}

func (d *CmdData) ToDanmu() DanmakuData {
	var danmu DanmakuData
	_ = json.Unmarshal(d.Data, &danmu)
	return danmu
}

func (d *CmdData) ToGift() GiftData {
	var gift GiftData
	_ = json.Unmarshal(d.Data, &gift)
	return gift
}

func (d *CmdData) ToSuperChat() SuperChatData {
	var superChat SuperChatData
	_ = json.Unmarshal(d.Data, &superChat)
	return superChat
}

func (d *CmdData) ToSuperChatDel() SuperChatDelData {
	var superChatDel SuperChatDelData
	_ = json.Unmarshal(d.Data, &superChatDel)
	return superChatDel
}

func (d *CmdData) ToNewGuard() NewGuardData {
	var newGuard NewGuardData
	_ = json.Unmarshal(d.Data, &newGuard)
	return newGuard
}

func (d *CmdData) ToLike() LikeData {
	var like LikeData
	_ = json.Unmarshal(d.Data, &like)
	return like
}

type UserInfo struct {
	UID    int    `json:"uid"`     // 用户uid
	OpenID string `json:"open_id"` // 用户唯一标识
	UName  string `json:"uname"`   // 用户昵称
	UFace  string `json:"uface"`   // 用户头像
}

type MedalInfo struct {
	FansMedalLevel         int    `json:"fans_medal_level"`          // 粉丝勋章等级
	FansMedalName          string `json:"fans_medal_name"`           // 粉丝勋章名
	FansMedalWearingStatus bool   `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况
}

type DanmakuData struct {
	RoomID    int    `json:"room_id"` // 弹幕接收的直播间
	UID       int    `json:"uid"`     // 用户UID
	OpenID    string `json:"open_id"` // 用户唯一标识
	UName     string `json:"uname"`   // 用户昵称
	Msg       string `json:"msg"`     // 弹幕内容
	MsgID     string `json:"msg_id"`  // 消息唯一id
	MedalInfo        // 对应房间勋章信息
	//FansMedalLevel         int    `json:"fans_medal_level"`          // 对应房间勋章信息
	//FansMedalName          string `json:"fans_medal_name"`           // 粉丝勋章名
	//FansMedalWearingStatus bool   `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况
	GuardLevel  int    `json:"guard_level"`   // 对应房间大航海 1总督 2提督 3舰长
	Timestamp   int64  `json:"timestamp"`     // 弹幕发送时间秒级时间戳
	UFace       string `json:"uface"`         // 用户头像
	EmojiImgURL string `json:"emoji_img_url"` // 表情包图片地址
	DanMuType   int    `json:"dm_type"`       // 弹幕类型 0：普通弹幕 1：表情包弹幕
}

type GiftData struct {
	RoomID   int    `json:"room_id"`   // 直播间 (In auditorium mode, it represents the auditorium live room; otherwise, it represents the gifting live room)
	UID      int    `json:"uid"`       // 送礼用户 UID
	OpenID   string `json:"open_id"`   // 用户唯一标识
	UName    string `json:"uname"`     // 送礼用户昵称
	UFace    string `json:"uface"`     // 送礼用户头像
	GiftID   int    `json:"gift_id"`   // 道具 ID (For blind boxes: the ID of the item obtained)
	GiftName string `json:"gift_name"` // 道具名 (For blind boxes: the name of the item obtained)
	GiftNum  int    `json:"gift_num"`  // 赠送道具数量
	Price    int    `json:"price"`     // 礼物单价 (1000 = 1元 = 10电池); for blind boxes: the value of the obtained item
	Paid     bool   `json:"paid"`      // 是否是付费道具
	MedalInfo
	//FansMedalLevel         int      `json:"fans_medal_level"`          // 实际收礼人的勋章信息
	//FansMedalName          string   `json:"fans_medal_name"`           // 粉丝勋章名
	//FansMedalWearingStatus bool     `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况
	GuardLevel int      `json:"guard_level"` // room_id 对应的大航海等级
	Timestamp  int64    `json:"timestamp"`   // 收礼时间秒级时间戳
	MsgID      string   `json:"msg_id"`      // 消息唯一 ID
	AnchorInfo UserInfo `json:"anchor_info"` // 收礼主播
	GiftIcon   string   `json:"gift_icon"`   // 道具 icon (新增)
	ComboGift  bool     `json:"combo_gift"`  // 是否是 combo 道具
	ComboInfo  struct {
		ComboBaseNum int    `json:"combo_base_num"` // 每次连击赠送的道具数量
		ComboCount   int    `json:"combo_count"`    // 连击次数
		ComboID      string `json:"combo_id"`       // 连击 ID
		ComboTimeout int    `json:"combo_timeout"`  // 连击有效期秒
	} `json:"combo_info"`
}

type SuperChatData struct {
	RoomID     int    `json:"room_id"`     // 直播间 ID
	UID        int    `json:"uid"`         // 购买用户 UID
	OpenID     string `json:"open_id"`     // 用户唯一标识
	UName      string `json:"uname"`       // 购买的用户昵称
	UFace      string `json:"uface"`       // 购买用户头像
	MessageID  int    `json:"message_id"`  // 留言 ID (In case of risk control, this may be used to recall the message)
	Message    string `json:"message"`     // 留言内容
	MsgID      string `json:"msg_id"`      // 消息唯一 ID
	RMB        int    `json:"rmb"`         // 支付金额 (元)
	Timestamp  int64  `json:"timestamp"`   // 赠送时间秒级
	StartTime  int64  `json:"start_time"`  // 生效开始时间
	EndTime    int64  `json:"end_time"`    // 生效结束时间
	GuardLevel int    `json:"guard_level"` // 对应房间大航海等级 (新增)
	MedalInfo
	//FansMedalLevel         int    `json:"fans_medal_level"`          // 对应房间勋章信息 (新增)
	//FansMedalName          string `json:"fans_medal_name"`           // 对应房间勋章名字 (新增)
	//FansMedalWearingStatus bool   `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况 (新增)
}

type SuperChatDelData struct {
	RoomID     int    `json:"room_id"`     // 直播间 ID
	MessageIDs []int  `json:"message_ids"` // 留言 ID 列表
	MsgID      string `json:"msg_id"`      // 消息唯一 ID
}

type NewGuardData struct {
	UserInfo   UserInfo `json:"user_info"`   // 用户信息
	GuardLevel int      `json:"guard_level"` // 对应的大航海等级 1 总督 2 提督 3 舰长
	GuardNum   int      `json:"guard_num"`   // 舰长数量
	GuardUnit  string   `json:"guard_unit"`  // 舰长数量单位 (个月)
	MedalInfo           // 该房间粉丝勋章
	//FansMedalLevel         int      `json:"fans_medal_level"`          // 粉丝勋章等级
	//FansMedalName          string   `json:"fans_medal_name"`           // 粉丝勋章名
	//FansMedalWearingStatus bool     `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况
	Timestamp int64  `json:"timestamp"` // 时间戳
	RoomID    int    `json:"room_id"`   // 直播间 ID
	MsgID     string `json:"msg_id"`    // 消息唯一 ID
}

type LikeData struct {
	UserInfo
	Timestamp int64  `json:"timestamp"` // 时间戳
	LikeText  string `json:"like_text"` // 点赞文本
	MedalInfo
	//FansMedalWearingStatus bool   `json:"fans_medal_wearing_status"` // 该房间粉丝勋章佩戴情况
	//FansMedalName          string `json:"fans_medal_name"`           // 粉丝勋章名
	//FansMedalLevel         int    `json:"fans_medal_level"`          // 粉丝勋章等级
	MsgID  string `json:"msg_id"`  // 消息唯一 ID
	RoomID int    `json:"room_id"` // 直播间 ID
}
