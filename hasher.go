package vcache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type Hasher interface {
	ToString(interface{}) string
}

type HashFunc func(interface{}) string

func (h HashFunc) ToString(x interface{}) string {
	return h(x)
}

func MustMD5(x interface{}) string {
	var (
		bs  []byte
		err error
	)

	switch x.(type) {
	case string:
		bs = []byte(x.(string))
	default:
		bs, err = json.Marshal(x)
		if err != nil {
			panic(err)
		}
	}

	digest := md5.Sum(bs)
	return hex.EncodeToString(digest[:])
}
