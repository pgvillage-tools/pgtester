// Package config reads tests from one or more .yml files and can be used as a
// starting point for running the tests
package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pgvillage-tools/pgtester/internal/version"
	"github.com/pgvillage-tools/pgtester/pkg/pg"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Tests is a struct that can hold a list of tests (as read from a tests file).
type Tests []Test

// Test is a struct that can hold test information regarding a test definition.
type Test struct {
	Name    string     `yaml:"name"`
	Query   string     `yaml:"query"`
	Results pg.Results `yaml:"results"`
	Reverse bool       `yaml:"reverse"`
}

const testsuccessful = 0
const testunsuccessful = 1

// Validate is is a function that checks if a Test definition is a valid Test.
func (t *Test) Validate() (err error) {
	const blank = ""
	if t.Name == blank {
		t.Name = t.Query
	} else if t.Query == blank {
		// Let's hope it is a query
		t.Query = t.Name
	}
	if t.Name == blank {
		return errors.New("a defined test is missing the query and name arguments")
	}
	return nil
}

// IncreaseOnError is a helper function that reverses a test when an error occurs.
func (t *Test) IncreaseOnError() (increase int) {
	if t.Reverse {
		return testsuccessful
	}
	return testunsuccessful
}

// IncreaseOnSuccess is a function that reverses the test result adder when there is an succes
func (t *Test) IncreaseOnSuccess() (increase int) {
	return testunsuccessful - t.IncreaseOnError()
}

// MsgOnError is a function that sends a message containing either 'expected'
// or 'unexpected' error when a error occurs.
func (t *Test) MsgOnError() (msg string) {
	if t.Reverse {
		return "expected error"
	}
	return "unexpected error"
}

// MsgOnSuccess is a function that sends a message containing either 'success as expected'
// or 'unexpected succes' when a succes occurs.
func (t *Test) MsgOnSuccess() (msg string) {
	if t.Reverse {
		return "unexpected success"
	}
	return "success as expected"
}

// Configs is a struct that can hold a number of configs.
type Configs []Config

// Config is a struct thatcan hold configuration data.
type Config struct {
	path    string
	index   int
	Debug   bool          `yaml:"debug"`
	Delay   time.Duration `yaml:"delay"`
	Retries uint          `yaml:"retries"`
	Tests   Tests         `yaml:"tests"`
	DSN     pg.Dsn        `yaml:"dsn"`
}

// Name returns the name of this testconfig
func (c Config) Name() (name string) {
	return fmt.Sprintf("%s (%d)", c.path, c.index)
}

// NewConfigFromReader
func newConfigsFromReader(reader io.Reader, name string) (configs Configs, err error) {
	var i int
	decoder := yaml.NewDecoder(reader)
	for {
		// create new spec here
		config := new(Config)
		// pass a reference to spec reference
		if err := decoder.Decode(&config); errors.Is(err, io.EOF) {
			log.Debug().Msg("EOF")
			break
		} else if err != nil {
			return Configs{}, err
		} else if config == nil {
			// check it was parsed
			log.Debug().Msg("empty config")
			continue
		} else if config.Delay.Nanoseconds() == 0 {
			config.Delay = time.Second
		}

		config.path = name
		config.index = i
		i++
		configs = append(configs, *config)
		log.Debug().Msgf("Number of configs: %d", len(configs))
	}
	return configs, nil
}

func newConfigsFromFile(path string) (c Configs, err error) {
	// This only parsed as yaml, nothing else
	// #nosec
	reader, err := os.Open(path)
	if err != nil {
		return c, err
	}
	return newConfigsFromReader(reader, path)
}

func newConfigsFromStdin() (configs Configs, err error) {
	log.Info().Msg("Parsing yaml from stdin")
	reader := bufio.NewReader(os.Stdin)
	return newConfigsFromReader(reader, "(stdin)")
}

// ReadFromFileOrDir returns an array of Configs parsed from all yaml files, found while recursively walking
// through a directory, while following symlinks as needed.
func ReadFromFileOrDir(path string) (configs Configs, err error) {
	const alldirectories = 0
	log.Info().Msgf("Parsing yaml from %s", path)
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return Configs{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Configs{}, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return Configs{}, err
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if fileInfo.IsDir() {
		// file is a directory
		entries, err := file.ReadDir(alldirectories)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
		// I want the entries sorted, so adding them to a list of strings
		var entryNames []string
		for _, entry := range entries {
			entryNames = append(entryNames, entry.Name())
		}
		sort.Strings(entryNames)
		for _, entryName := range entryNames {
			newConfigs, err := ReadFromFileOrDir(filepath.Join(path, entryName))
			if err != nil {
				_ = file.Close()
				return Configs{}, err
			}
			configs = append(configs, newConfigs...)
		}
	} else {
		// file is not a directory
		configs, err = newConfigsFromFile(path)
		if err != nil {
			_ = file.Close()
			return Configs{}, err
		}
	}
	return configs, file.Close()
}

// GetConfigs is a function that fetches configs from either paths or stdin and returns a list of testconfigs.
func GetConfigs() (configs Configs, err error) {
	var debug bool
	var myVersion bool
	const pathsIsEmpty = 0
	flag.BoolVar(&debug, "d", false, "Add debugging output")
	flag.BoolVar(&myVersion, "v", false, "Show version information")

	flag.Parse()
	if myVersion {
		fmt.Println(version.GetAppVersion())
		os.Exit(testsuccessful)
	}
	paths := flag.Args()
	if len(paths) == pathsIsEmpty {
		return newConfigsFromStdin()
	}
	for _, path := range paths {
		newConfigs, err := ReadFromFileOrDir(path)
		if err != nil {
			return Configs{}, err
		}
		configs = append(configs, newConfigs...)
	}

	for i := range configs {
		configs[i].Debug = configs[i].Debug || debug
	}

	return configs, err
}
