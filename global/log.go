package global

import (
	"go.uber.org/zap"
)

var (
	Log *zap.SugaredLogger
	logger *zap.Logger
)

func LoadLogger(testing bool) error {

	var err error

	if testing {
		logger = zap.NewNop()
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return err
		}
	}

	Log = logger.Sugar()
	return nil
}