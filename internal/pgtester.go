package internal

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mannemsolutions/pgtester/pkg/pg"
	"os"
)

var (
	log  *zap.SugaredLogger
	atom zap.AtomicLevel
)

func Initialize() {
	atom = zap.NewAtomicLevel()
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	log = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).Sugar()

	pg.Initialize(log)
}

func Handle() {
	var errors int
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("could not parse config: %e", err)
	}
	if config.Debug {
		atom.SetLevel(zapcore.DebugLevel)
	}
	conn := pg.NewConn(config.DSN)
	for _, test := range config.Tests {
		result, err := conn.RunQueryGetOneField(test.Query)
		if err != nil {
			log.Errorf("err occurred on query %s: %e", test.Query, err)
		} else {
			err = result.Compare(test.Results)
			if err != nil {
				log.Errorf("err occurred on query %s: %e", test.Query, err)
			} else {
				log.Infof("succes on query %s", test.Query)
				continue
			}
		}
		errors += 1
	}
	if errors > 0 {
		log.Fatalf("finished with %d errors", errors)
	}
	log.Infof("finished with %d errors", errors)
}