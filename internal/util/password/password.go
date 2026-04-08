package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// 参数参考 OWASP 推荐：https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16
)

// Hash 生成 argon2id 哈希，格式：$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
func Hash(plain string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(plain), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// Verify 验证明文密码与哈希是否匹配
func Verify(plain, encoded string) bool {
	if !strings.HasPrefix(encoded, "$argon2id$") {
		return false
	}
	parts := strings.Split(encoded, "$")
	// $argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>  → 6 parts after split by $
	if len(parts) != 6 {
		return false
	}
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false
	}
	var m, t uint32
	var p uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	wantHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	gotHash := argon2.IDKey([]byte(plain), salt, t, m, p, uint32(len(wantHash)))
	return subtle.ConstantTimeCompare(gotHash, wantHash) == 1
}

// ErrUnsupportedFormat 用于区分哈希格式不支持的情况
var ErrUnsupportedFormat = errors.New("unsupported password hash format")
