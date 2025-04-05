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

func TestConfig(t *testing.T) {
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
---
dsn:
  host: postgres

tests:
- query: "select count(*) total from pg_database"
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
	myConfig1 := myConfigs[0]
	assert.Len(t, myConfig1.Tests, 1)

}
