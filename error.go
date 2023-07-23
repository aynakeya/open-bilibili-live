package openblive

import "fmt"

type PublicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"-"`
}

func (e *PublicError) WithDetail(detail error) *PublicError {
	return &PublicError{
		Code:    e.Code,
		Message: e.Message,
		Detail:  detail.Error(),
	}
}

func (e *PublicError) Error() string {
	if e.Detail == "" {
		return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("Error %d: %s (%s)", e.Code, e.Message, e.Detail)

}

var _errors = map[int]*PublicError{}

func NewPublicError(code int, msg string, detail string) *PublicError {
	_, ok := _errors[code]
	if ok {
		panic("error code already exists")
	}
	e := &PublicError{
		Code:    code,
		Message: msg,
		Detail:  detail,
	}
	_errors[code] = e
	return e
}

func GetErrorFromCode(errcode int) *PublicError {
	e, ok := _errors[errcode]
	if ok {
		return e
	}
	return ErrUnknown
}

var (
	ErrUnknown                 = NewPublicError(0, "unknown error", "")
	ErrInvalidParameter        = NewPublicError(4000, "参数错误", "请检查必填参数，参数大小限制")
	ErrInvalidApp              = NewPublicError(4001, "应用无效", "请检查header的x-bili-accesskeyid是否为空，或者有效")
	ErrSignature               = NewPublicError(4002, "签名异常", "请检查header的Authorization")
	ErrExpiredRequest          = NewPublicError(4003, "请求过期", "请检查header的x-bili-timestamp")
	ErrDuplicateRequest        = NewPublicError(4004, "重复请求", "请检查header的x-bili-nonce")
	ErrInvalidSignatureMethod  = NewPublicError(4005, "签名method异常", "请检查header的x-bili-signature-method")
	ErrInvalidVersion          = NewPublicError(4006, "版本异常", "请检查header的x-bili-version")
	ErrIPWhitelist             = NewPublicError(4007, "IP白名单限制", "请确认请求服务器是否在报备的白名单内")
	ErrPermission              = NewPublicError(4008, "权限异常", "请确认接口权限")
	ErrAPILimit                = NewPublicError(4009, "接口访问限制", "请确认接口权限及请求频率")
	ErrNotFound                = NewPublicError(4010, "接口不存在", "请确认请求接口url")
	ErrInvalidContentType      = NewPublicError(4011, "Content-Type不为application/json", "请检查header的Content-Type")
	ErrMD5Validation           = NewPublicError(4012, "MD5校验失败", "请检查header的x-bili-content-md5")
	ErrInvalidAcceptType       = NewPublicError(4013, "Accept不为application/json", "请检查header的Accept")
	ErrService                 = NewPublicError(5000, "服务异常", "请联系B站对接同学")
	ErrRequestTimeout          = NewPublicError(5001, "请求超时", "请求超时")
	ErrInternal                = NewPublicError(5002, "内部错误", "请联系B站对接同学")
	ErrConfiguration           = NewPublicError(5003, "配置错误", "请联系B站对接同学")
	ErrRoomWhitelist           = NewPublicError(5004, "房间白名单限制", "请联系B站对接同学")
	ErrRoomBlacklist           = NewPublicError(5005, "房间黑名单限制", "请联系B站对接同学")
	ErrInvalidVerificationCode = NewPublicError(6000, "验证码错误", "验证码校验失败")
	ErrInvalidPhoneNumber      = NewPublicError(6001, "手机号码错误", "检查手机号码")
	ErrExpiredVerificationCode = NewPublicError(6002, "验证码已过期", "验证码超过规定有效期")
	ErrVerificationRateLimit   = NewPublicError(6003, "验证码频率限制", "检查获取验证码的频率")
	ErrNotInGame               = NewPublicError(7000, "不在游戏内", "当前房间未进行互动游戏")
	ErrRequestCooldown         = NewPublicError(7001, "请求冷却期", "上个游戏正在结算中，建议10秒后进行重试")
	ErrRoomInGame              = NewPublicError(7002, "房间重复游戏", "当前房间正在进行游戏,无法开启下一局互动游戏")
	ErrExpiredHeartbeat        = NewPublicError(7003, "心跳过期", "当前game_id错误或互动游戏已关闭")
	ErrMaxHeartbeatBatchSize   = NewPublicError(7004, "批量心跳超过最大值", "批量心跳单次最大值为200")
	ErrDuplicateHeartbeatID    = NewPublicError(7005, "批量心跳ID重复", "批量心跳game_id存在重复,请检查参数")
	ErrInvalidIdentityCode     = NewPublicError(7007, "身份码错误", "请检查身份码是否正确")
	ErrNoProjectAccess         = NewPublicError(8002, "项目无权限访问", "确认项目ID是否正确")
)
