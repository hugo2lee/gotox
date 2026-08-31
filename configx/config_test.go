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

	c := configx.New(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	assert.Equal(t, "test.log", c.LogDir())

	cu := custom{c}
	assert.Equal(t, "custom log dir", cu.LogDir())
	assert.Equal(t, map[string]string{"client": "client-auth"}, cu.Auths())
}
