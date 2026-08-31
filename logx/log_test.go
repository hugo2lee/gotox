/*
 * @Author: hugo
 * @Date: 2024-03-19 19:57
 * @LastEditors: hugo
 * @LastEditTime: 2024-04-19 16:44
 * @FilePath: \gotox\logx\log_test.go
 * @Description:
 *
 * Copyright (c) 2024 by hugo, All Rights Reserved.
 */
package logx_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/logx"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "logs", "test.log")
	configContent := fmt.Sprintf("[log]\ndir = %q\n", filepath.ToSlash(logFile))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dev.toml"), []byte(configContent), 0o600))

	conf, err := configx.Load(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)
	logger, err := logx.Open(conf)
	require.NoError(t, err)
	logger.Debug("debug %v", 1)
	logger.Info("info %v", 1)

	_, err = os.Stat(logFile)
	require.NoError(t, err)
}

func TestOpenReturnsLogDirectoryError(t *testing.T) {
	dir := t.TempDir()
	blockingPath := filepath.Join(dir, "not-a-directory")
	require.NoError(t, os.WriteFile(blockingPath, []byte("file"), 0o600))

	logFile := filepath.Join(blockingPath, "test.log")
	configContent := fmt.Sprintf("[log]\ndir = %q\n", filepath.ToSlash(logFile))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dev.toml"), []byte(configContent), 0o600))

	conf, err := configx.Load(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)

	_, err = logx.Open(conf)
	require.Error(t, err)
}
