package cache

import (
	"crypto/md5"
	"github.com/phzfi/RIC/server/ops"
	"strings"
)

type cacheKey string

// Returns a unique representation of an ops chain. This unique representation can be used as a map key unlike the original ops chain (slice cannot be a key).
func toKey(operations []ops.Operation) cacheKey {
	marshaled := make([]string, len(operations))
	for i, op := range operations {
		marshaled[i] = op.Marshal()
	}
	bytes := md5.Sum([]byte(strings.Join(marshaled, "")))
	return cacheKey(bytes[:])
}
