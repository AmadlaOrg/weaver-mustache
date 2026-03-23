package input

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_JSONObject(t *testing.T) {
	result, err := Parse([]byte(`{"key": "value"}`))
	require.NoError(t, err)
	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "value", m["key"])
}

func TestParse_JSONArray(t *testing.T) {
	result, err := Parse([]byte(`[{"key": "value"}]`))
	require.NoError(t, err)
	arr, ok := result.([]any)
	require.True(t, ok)
	assert.Len(t, arr, 1)
}

func TestParse_YAMLObject(t *testing.T) {
	result, err := Parse([]byte("key: value\nother: data"))
	require.NoError(t, err)
	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "value", m["key"])
	assert.Equal(t, "data", m["other"])
}

func TestParse_YAMLList(t *testing.T) {
	result, err := Parse([]byte("- key: value\n- key: other"))
	require.NoError(t, err)
	arr, ok := result.([]any)
	require.True(t, ok)
	assert.Len(t, arr, 2)
}

func TestParse_EmptyJSON(t *testing.T) {
	result, err := Parse([]byte(`{}`))
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "data.json")
	err := os.WriteFile(f, []byte(`{"name": "test"}`), 0644)
	require.NoError(t, err)

	result, err := Load(f)
	require.NoError(t, err)
	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", m["name"])
}

func TestLoad_FromYAMLFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "data.yaml")
	err := os.WriteFile(f, []byte("name: test\nport: 8080"), 0644)
	require.NoError(t, err)

	result, err := Load(f)
	require.NoError(t, err)
	m, ok := result.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "test", m["name"])
	assert.Equal(t, 8080, m["port"])
}

func TestLoad_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "empty.json")
	err := os.WriteFile(f, []byte(""), 0644)
	require.NoError(t, err)

	result, err := Load(f)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path.json")
	assert.Error(t, err)
}
