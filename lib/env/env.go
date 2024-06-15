package env

import (
	"os"
	"strconv"
	"time"
)

// Parse environment variable as string value.
func String(key string, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return v
}

// Parse environment variable as byte value.
func Byte(key string, defaultValue byte) byte {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseUint(v, 0, 8)
	if err != nil {
		return defaultValue
	}
	return byte(i)
}

// Parse environment variable as int32 value.
func Int32(key string, defaultValue int32) int32 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 0, 32)
	if err != nil {
		return defaultValue
	}
	return int32(i)
}

// Parse environment variable as int64 value.
func Int64(key string, defaultValue int64) int64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 0, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// Parse environment variable as uint32 value.
func Uint32(key string, defaultValue uint32) uint32 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseUint(v, 0, 32)
	if err != nil {
		return defaultValue
	}
	return uint32(i)
}

// Parse environment variable as uint64 value.
func Uint64(key string, defaultValue uint64) uint64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseUint(v, 0, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// Parse environment variable as bool value.
func Bool(key string, defaultValue bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return i
}

// Parse environment variable as float64 value.
func Float64(key string, defaultValue float64) float64 {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// Parse environment variable as time struct.
// Required format is RFC3339.
//
// Example: 2006-01-02T15:04:05Z07:00
func Time(key string, defaultValue time.Time) time.Time {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return defaultValue
	}
	return t
}

// Parse environment variable as time struct.
// Required format is RFC3339.
//
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Example: "300ms", "-1.5h" or "2h45m"
func Duration(key string, defaultValue time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultValue
	}
	return d
}
