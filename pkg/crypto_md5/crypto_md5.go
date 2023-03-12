package crypto_md5

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "夏天夏天"

// EncyptPassword 密码加密函数
func EncyptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
