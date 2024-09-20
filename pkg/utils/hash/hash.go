package hash

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"encoding/json"
)

func Hash(in any) string {
	if in == nil {
		return ""
	}
	bytes, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	hash := md5.Sum(bytes) //nolint:gosec
	return hex.EncodeToString(hash[:])
}
