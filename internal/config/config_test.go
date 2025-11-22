package config

import (
	"fmt"
	"os"
	"path"
	"testing"

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
	const fileReadWritePermission = 0o0600
	const myconfigslength = 2
	const myconfigs0length = 3
	const myconfigs1length = 1
	tmpDir, err := os.MkdirTemp("", "config")
	if err != nil {
		panic(fmt.Errorf("unable to create temp dir: %w", err))
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	myCredFile := path.Join(tmpDir, "my-yaml")
	require.NoError(t, os.WriteFile(myCredFile, []byte(myConfigData), fileReadWritePermission))
	myConfigs, err := newConfigsFromFile(myCredFile)
	assert.NoError(t, err)
	require.Len(t, myConfigs, myconfigslength)
	require.Len(t, myConfigs[0].Tests, myconfigs0length)
	require.Len(t, myConfigs[1].Tests, myconfigs1length)
}

func TestTestValidate(t *testing.T) {
	const testcount0 = 0
	const testcount1 = 1
	const select1 = "select 1"
	const select2 = "select 2"
	for _, check := range []struct {
		test  Test
		valid bool
		count int
	}{
		{test: Test{Name: select1}, valid: true, count: testcount1},
		{test: Test{Query: select1, Name: select2}, valid: true, count: testcount1},
		{test: Test{Query: select1}, valid: true, count: testcount1},
		{test: Test{}, count: testcount1},
		{test: Test{Name: select1, Reverse: true}, valid: true, count: testcount0},
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
