/*
 * @Author: hugo
 * @Date: 2024-03-12 15:28
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-25 16:55
 * @FilePath: \gotox\configx\config_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package configx_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type custom struct {
	*configx.Configx
}

func (c custom) LogDir() string {
	return "custom log dir"
}

func TestConfigExample(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "dev.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[log]
dir = "test.log"

[auths]
client = "client-auth"
`), 0o600))

	c, err := configx.Load(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)
	assert.Equal(t, "test.log", c.LogDir())

	cu := custom{c}
	assert.Equal(t, "custom log dir", cu.LogDir())
	assert.Equal(t, map[string]string{"client": "client-auth"}, cu.Auths())
}

func TestLoadReturnsInvalidModeError(t *testing.T) {
	_, err := configx.Load(configx.WithMode("invalid"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid mode")
}

func TestLoadReturnsInvalidPathError(t *testing.T) {
	_, err := configx.Load(configx.WithPath(""))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestLoadReturnsMissingConfigError(t *testing.T) {
	_, err := configx.Load(configx.WithPath(t.TempDir()), configx.WithMode(configx.RUNDEV))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read dev config")
}
