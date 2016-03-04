package cache

import (
	"crypto/md5"
	"github.com/phzfi/RIC/server/ops"
	"reflect"
	"unsafe"
)

type cacheKey string

// Returns a unique representation of an ops chain. This unique representation can be used as a map key unlike the original ops chain (slice cannot be a key).
func toKey(operations []ops.Operation) cacheKey {
	result := make([]byte, 0)

	for _, op := range operations {
		operation := reflect.ValueOf(&op).Elem()
		inner := operation.Elem()
		if inner.Kind() == reflect.Ptr {
			inner = inner.Elem()
		}
		result = append(result, []byte(inner.Type().Name())...)
		for i := 0; i < inner.NumField(); i++ {
			if inner.Type().Field(i).PkgPath == "" {
				// is an exported field
				if field := inner.Field(i); field.Kind() == reflect.String {
					result = append(result, []byte(field.String())...)
				} else {
					address := operation.InterfaceData()[1] + inner.Type().Field(i).Offset
					copyMemoryToByteSlice(address, field.Type().Size(), &result)
				}
			}
		}
	}

	hash := md5.Sum(result)
	return cacheKey(hash[:])
}

func copyMemoryToByteSlice(address, size uintptr, slice *[]byte) {
	//*slice = adress[1] - address[0]
	for pos := address; pos < address+size; pos++ {
		*slice = append(*slice, *(*byte)(unsafe.Pointer(pos)))
	}
}
