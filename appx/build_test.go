package appx_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hugo2lee/gotox/appx"
	"github.com/hugo2lee/gotox/configx"
	"github.com/hugo2lee/gotox/ormx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeAppConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	logFile := filepath.Join(dir, "logs", "app.log")
	content := fmt.Sprintf("[log]\ndir = %q\n", filepath.ToSlash(logFile))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "dev.toml"), []byte(content), 0o600))
	return dir
}

func TestBuildReturnsConfiguredApp(t *testing.T) {
	dir := writeAppConfig(t)

	app, err := appx.Build(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)
	require.NotNil(t, app)
	assert.NotNil(t, app.Configx)
	assert.NotNil(t, app.Logger)
}

func TestBuildReturnsConfigError(t *testing.T) {
	_, err := appx.Build(configx.WithPath(t.TempDir()), configx.WithMode(configx.RUNDEV))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "load config")
}

func TestTryEnableDBReturnsOptionError(t *testing.T) {
	dir := writeAppConfig(t)
	app, err := appx.Build(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)

	expected := errors.New("db option failed")
	err = app.TryEnableDB(func(*ormx.Ormx) error {
		return expected
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expected)
	assert.Nil(t, app.DBs)
}

func TestMigrateTablesReturnsCallbackError(t *testing.T) {
	dir := writeAppConfig(t)
	app, err := appx.Build(configx.WithPath(dir), configx.WithMode(configx.RUNDEV))
	require.NoError(t, err)

	expected := errors.New("migration failed")
	err = app.MigrateTables(
		func() error { return nil },
		func() error { return expected },
	)

	require.Error(t, err)
	assert.ErrorIs(t, err, expected)
	assert.Contains(t, err.Error(), "callback 1")
}
