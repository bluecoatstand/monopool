package algorithm

import (
	"encoding/hex"
	"math/big"
	"strings"

	"monopool/utils"

	logging "github.com/ipfs/go-log/v2"
	"github.com/samli88/go-x11-hash"
	"golang.org/x/crypto/scrypt"
)

var log = logging.Logger("algorithm")

// difficulty = MAX_TARGET / current_target.
var (
	MaxTargetTruncated, _ = new(big.Int).SetString("00000000FFFF0000000000000000000000000000000000000000000000000000", 16)
	MaxTarget, _          = new(big.Int).SetString("00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 16)
)

func GetHashFunc(hashName string) func([]byte) []byte {
	switch strings.ToLower(hashName) {
	case "scrypt":
		return ScryptHash
	case "x11":
		return X11Hash
	case "sha256d":
		return DoubleSha256Hash
	case "sha256dt":
		return TaggedDoubleSha256
	default:
		log.Panic(hashName, " is not officially supported yet, but you can easily add it with cgo binding by yourself")
		return nil
	}
}

// ScryptHash is the algorithm which litecoin uses as its PoW mining algorithm
func ScryptHash(data []byte) []byte {
	b, _ := scrypt.Key(data, data, 1024, 1, 1, 32)

	return b
}

// X11Hash is the algorithm which dash uses as its PoW mining algorithm
func X11Hash(data []byte) []byte {
	dst := make([]byte, 32)
	x11.New().Hash(data, dst)
	return dst
}

// DoubleSha256Hash is the algorithm which litecoin uses as its PoW mining algorithm
func DoubleSha256Hash(b []byte) []byte {
	return utils.Sha256d(b)
}

var tag, _ = hex.DecodeString("ce04921a0d6d9badccfec07629a154a6637383d6979ba5eb6298c1667a2a5a11")

func TaggedDoubleSha256(b []byte) []byte {
	var payload []byte

	payload = append(payload, tag...)
	payload = append(payload, tag...)

	payload = append(payload, b...)

	hash := utils.Sha256(payload)

	var payload2 []byte
	payload2 = append(payload2, tag...)
	payload2 = append(payload2, tag...)
	payload2 = append(payload2, hash...)

	hash2 := utils.Sha256(payload2)
	hash2 = utils.ReverseBytes(hash2)

	return hash2
}
