package timeout_control_helper

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/global/variable"
	"time"
)

var timeout = 3 * time.Second

func TimeoutControl(f func(chan error), actionMsg string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error)
	go func() {
		// 调用主方法逻辑
		f(errChan)
		// 调用完成，关闭错误通道
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		variable.ZapLog.Error(fmt.Sprintf("%s call timeout", actionMsg), zap.Error(ctx.Err()))
	case <-errChan:
		variable.ZapLog.Info(fmt.Sprintf("%s call finish", actionMsg), zap.Error(<-errChan))
	}
}
