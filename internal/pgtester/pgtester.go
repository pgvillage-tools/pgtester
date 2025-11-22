// Package pgtester is a package for testing postgres databases
package pgtester

import (
	"fmt"
	"os"
	"strings"

	"github.com/pgvillage-tools/pgtester/internal/config"
	"github.com/pgvillage-tools/pgtester/pkg/pg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.SugaredLogger
	atom zap.AtomicLevel
)

// Initialize is a function that prepares postgres
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

// Handle is a function that configures postgres, using all config files and checks for errors
func Handle() {
	var errors int
	const (
		eCANCELED     = 125
		lineCharacter = "="
		maxExitCode   = 125
	)
	configs, err := config.GetConfigs()
	if err != nil {
		log.Errorf("could not parse all configs: %s", err.Error())
		os.Exit(eCANCELED)
	}
	for _, cfg := range configs {
		name := cfg.Name()
		log.Infof(strings.Repeat(lineCharacter, 19+len(name)))
		log.Infof("Running tests from %s", name)
		log.Infof(strings.Repeat(lineCharacter, 19+len(name)))
		if cfg.Debug {
			atom.SetLevel(zapcore.DebugLevel)
		} else {
			atom.SetLevel(zapcore.InfoLevel)
		}
		conn := pg.NewConn(cfg.DSN, cfg.Retries, cfg.Delay)
		for i, test := range cfg.Tests {
			err = test.Validate()
			if err != nil {
				log.Errorf("Test %d (%s): %s occurred : %s", i, test.Name, test.MsgOnError(), err.Error())
				errors += test.IncreaseOnError()
				continue
			}
			result, err := conn.RunQueryGetOneField(test.Query)
			if err != nil {
				log.Errorf("Test %d (%s): occurred : %s", i, test.Name, test.MsgOnError(), err.Error())
				errors += test.IncreaseOnError()
			} else {
				err = result.Compare(test.Results)
				if err != nil {
					log.Errorf("Test %d (%s): occurred : %s", i, test.Name, test.MsgOnError(), err.Error())
					errors += test.IncreaseOnError()
				} else {
					log.Infof("Test %d (%s): '%s'", i, test.Name, test.MsgOnSuccess())
					errors += test.IncreaseOnSuccess()
				}
			}
		}
	}
	if errors > 0 {
		msg := fmt.Sprintf("unfortunately finished with %d unexpected results", errors)
		log.Errorf(strings.Repeat(lineCharacter, len(msg)))
		log.Errorf(msg)
		log.Errorf(strings.Repeat(lineCharacter, len(msg)))
		// exit code should be between 0 and 125, where 0 is success, and 125 is config error
		// so change to a number between 1 and 125
		errors = ((errors - 1) % (maxExitCode - 1)) + 1
		os.Exit(errors)
	}
	const exitString = "successfully finished without unexpected results"
	log.Infof(strings.Repeat(lineCharacter, len(exitString)))
	log.Infof(exitString)
	log.Infof(strings.Repeat(lineCharacter, len(exitString)))
}
