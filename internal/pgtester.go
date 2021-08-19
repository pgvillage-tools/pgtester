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
		log.Errorf("could not parse config: %s", err.Error())
		os.Exit(125)
	}
	if config.Debug {
		atom.SetLevel(zapcore.DebugLevel)
	}
	conn := pg.NewConn(config.DSN)
	for _, test := range config.Tests {
		err = test.Validate()
		if err != nil {
			log.Errorf("err occurred on test: %s", err.Error())
			errors += 1
			continue
		}
		result, err := conn.RunQueryGetOneField(test.Query)
		if err != nil {
			log.Errorf("err occurred on test '%s': %s", test.Name, err.Error())
		} else {
			err = result.Compare(test.Results)
			if err != nil {
				log.Errorf("err occurred on test '%s': %s", test.Name, err.Error())
			} else {
				log.Infof("succes on test '%s'", test.Name)
				continue
			}
		}
		errors += 1
	}
	if errors > 0 {
		// exit code should be between 0 and 125, where 0 is success, and 125 is config error
		// so change to a number between 1 and 125
		errors = ((errors -1 ) % 124) + 1
		log.Errorf("finished with %d errors", errors)
		os.Exit(errors)
	}
	log.Infof("finished without errors")
}