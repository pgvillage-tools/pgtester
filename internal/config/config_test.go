package config_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/pgvillage-tools/pgtester/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	const myConfigData = `
---
dsn:
  host: postgres

retries: 60
delay: 1s
debug: false

tests:
- query: "select count(*) total from pg_database"
  results:
  - total: 3
- name: "select count(*) total from pg_database"
  results:
  - total: 3
- results:
  - total: 3
---
dsn:
  host: postgres

tests:
- name: number of databases
  query: "select count(*) total from pg_database"
  results:
  - total: 3
  `
	tmpDir, err := os.MkdirTemp("", "config")
	if err != nil {
		panic(fmt.Errorf("unable to create temp dir: %w", err))
	}
	defer os.RemoveAll(tmpDir)
	myCredFile := path.Join(tmpDir, "my-config.yaml")
	require.NoError(t, os.WriteFile(myCredFile, []byte(myConfigData), 0o0600))
	myConfigs, err := config.NewConfigsFromFile(myCredFile)
	assert.NoError(t, err)
	require.Len(t, myConfigs, 2)
	require.Len(t, myConfigs[0].Tests, 3)
	require.Len(t, myConfigs[1].Tests, 1)
}

func TestTestValidate(t *testing.T) {
	for _, check := range []struct {
		test  config.Test
		valid bool
		count int
	}{
		{test: config.Test{Name: "select 1"}, valid: true, count: 1},
		{test: config.Test{Query: "select 1", Name: "select 2"}, valid: true, count: 1},
		{test: config.Test{Query: "select 1"}, valid: true, count: 1},
		{test: config.Test{}, count: 1},
		{test: config.Test{Name: "select 1", Reverse: true}, valid: true, count: 0},
	} {
		t.Logf("test: %v", check.test)
		if check.valid {
			assert.NoError(t, check.test.Validate())
		} else {
			assert.Error(t, check.test.Validate())
		}
		assert.Equal(t, check.test.IncreaseOnError(), check.count)
	}
}
