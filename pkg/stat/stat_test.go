package stat

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	l, _ := zap.NewDevelopment()
	Init(
		// test empty
		WithLog(nil),
		WithPrintInterval(0),

		WithLog(l),
		WithPrintInterval(time.Second),

		WithAlarm(WithCPUThreshold(0.9), WithMemoryThreshold(0.85)),
	)

	time.Sleep(time.Second * 2)
}
