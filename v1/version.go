package v1

// Version 代表policy的版本，不同版本对策略解析的方式会有不同
// 目前v1版本
type Version uint8

const (
	V1 Version = 1
)
