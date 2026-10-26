package helper

import (
	"time"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/gin-gonic/gin"
)

// streamDrainEligibleKey 在流 goroutine 启动前冻结排空决定，确保读流和写入守卫一致。
const streamDrainEligibleKey = "stream_drain_eligible"

// ClientGoneButDraining 判断客户端是否已断开且本次流允许继续读取上游。
func ClientGoneButDraining(c *gin.Context) bool {
	return requestContextDone(c) && c.GetBool(streamDrainEligibleKey)
}

// waitForUpstreamAfterClientGone 等待上游结束、流空闲超时或排空等待超时，以获取最终用量。
func waitForUpstreamAfterClientGone(c *gin.Context, info *relaycommon.RelayInfo, stopChan <-chan bool, ticker *time.Ticker) {
	drainTimer := time.NewTimer(time.Duration(operation_setting.GetStreamDrainSetting().GetDrainMaxWaitSeconds()) * time.Second)
	defer drainTimer.Stop()

	select {
	case <-stopChan:
	case <-ticker.C:
		info.StreamStatus.SetEndReason(relaycommon.StreamEndReasonTimeout, nil)
	case <-drainTimer.C:
		info.StreamStatus.SetEndReason(relaycommon.StreamEndReasonClientGone, c.Request.Context().Err())
	}
}
