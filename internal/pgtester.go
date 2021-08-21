package internal

import (
	"fmt"
	"github.com/mannemsolutions/pgtester/pkg/pg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
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
	configs, err := GetConfigs()
	if err != nil {
		log.Errorf("could not parse all configs: %s", err.Error())
		os.Exit(125)
	}
	for _, config := range configs {
		name := config.Name()
		log.Infof(strings.Repeat("=", 19+len(name)))
		log.Infof("Running tests from %s", name)
		log.Infof(strings.Repeat("=", 19+len(name)))
		if config.Debug {
			atom.SetLevel(zapcore.DebugLevel)
		} else {
			atom.SetLevel(zapcore.InfoLevel)
		}
		conn := pg.NewConn(config.DSN, config.Retries, config.Delay)
		for _, test := range config.Tests {
			err = test.Validate()
			if err != nil {
				log.Errorf("%s occurred on test '%s': %s", test.MsgOnError(), test.Name, err.Error())
				errors += test.IncreaseOnError()
				continue
			}
			result, err := conn.RunQueryGetOneField(test.Query)
			if err != nil {
				log.Errorf("%s occurred on test '%s': %s", test.MsgOnError(), test.Name, err.Error())
				errors += test.IncreaseOnError()
			} else {
				err = result.Compare(test.Results)
				if err != nil {
					log.Errorf("%s occurred on test '%s': %s", test.MsgOnError(), test.Name, err.Error())
					errors += test.IncreaseOnError()
				} else {
					log.Infof("%s on test '%s'", test.MsgOnSuccess(), test.Name)
					errors += test.IncreaseOnSuccess()
				}
			}
		}
	}
	if errors > 0 {
		// exit code should be between 0 and 125, where 0 is success, and 125 is config error
		// so change to a number between 1 and 125
		msg := fmt.Sprintf("unfortunately finished with %d unexpected results", errors)
		log.Errorf(strings.Repeat("=", len(msg)))
		log.Errorf(msg)
		log.Errorf(strings.Repeat("=", len(msg)))
		errors = ((errors - 1) % 124) + 1
		os.Exit(errors)
	}
	log.Infof(strings.Repeat("=", 47))
	log.Infof("succesfully finished without unexpected results")
	log.Infof(strings.Repeat("=", 47))
}
