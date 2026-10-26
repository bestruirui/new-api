package operation_setting

import "github.com/QuantumNous/new-api/setting/config"

const (
	defaultDrainMaxWaitSeconds = 600
	maxDrainMaxWaitSeconds     = 3600
)

// StreamDrainSetting 控制客户端断开后继续读取上游流，以获取真实用量结算。
type StreamDrainSetting struct {
	DrainOnClientDisconnect bool `json:"drain_on_client_disconnect"` // 是否在客户端断开后继续读取上游流
	DrainMaxWaitSeconds     int  `json:"drain_max_wait_seconds"`     // 客户端断开后最多等待的秒数
}

var streamDrainSetting = StreamDrainSetting{ // 客户端断开后的排空配置
	DrainOnClientDisconnect: true,
	DrainMaxWaitSeconds:     defaultDrainMaxWaitSeconds,
}

// init 注册客户端断开后的排空配置。
func init() {
	config.GlobalConfig.Register("stream_drain_setting", &streamDrainSetting)
}

// GetStreamDrainSetting 返回客户端断开后的排空配置。
func GetStreamDrainSetting() *StreamDrainSetting {
	return &streamDrainSetting
}

// GetDrainMaxWaitSeconds 返回排空等待秒数，非正值使用默认值，上限为一小时。
func (s *StreamDrainSetting) GetDrainMaxWaitSeconds() int {
	if s.DrainMaxWaitSeconds <= 0 {
		return defaultDrainMaxWaitSeconds
	}
	return min(s.DrainMaxWaitSeconds, maxDrainMaxWaitSeconds)
}
