package openblive

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

// referer to https://open-live.bilibili.com/document/74eec767-e594-7ddd-6aba-257e8317c05d
func TestApiHeader_ToHeaderStr(t *testing.T) {
	ah := apiHeader{
		TimeStamp:        "1624594467",
		SignatureVersion: "1.0",
		SignatureNonce:   "ad184c09-095f-91c3-0849-230dd3744045",
		SignatureMethod:  "HMAC-SHA256",
		ContentMD5:       "fa6837e35b2f591865b288dfd859ce9d",
		AccessKey:        "xxxx",
	}
	authEncoder := hmac.New(sha256.New, []byte("JzOzZfSHeYYnAMZ"))
	assert.Equal(t,
		"x-bili-accesskeyid:xxxx\nx-bili-content-md5:fa6837e35b2f591865b288dfd859ce9d\nx-bili-signature-method:HMAC-SHA256\nx-bili-signature-nonce:ad184c09-095f-91c3-0849-230dd3744045\nx-bili-signature-version:1.0\nx-bili-timestamp:1624594467",
		ah.ToHeaderStr())
	authEncoder.Write([]byte(ah.ToHeaderStr()))
	assert.Equal(t,
		"a81c50234b6bbf15bc56e387ee4f19c6f871af2f70b837dc56db16517d4a341f",
		hex.EncodeToString(authEncoder.Sum(nil)))
}
