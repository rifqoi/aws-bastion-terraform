package logs

import (
	"go.uber.org/zap"
)

func GetLogger(opts ...zap.Option) *zap.SugaredLogger {
	if opts != nil {
		opts = append(opts, zap.AddCallerSkip(2), zap.AddStacktrace(zap.ErrorLevel))
	}
	l, _ := zap.NewProduction(opts...)
	defer l.Sync()
	sugar := l.Sugar()

	return sugar
}
