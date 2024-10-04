package util

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

func SizeInGB(size int64) float64 {
	return float64(size) / GB
}
