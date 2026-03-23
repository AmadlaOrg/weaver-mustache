package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "weaver-mustache")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run(), "failed to build binary")
	return bin
}

func TestInfoCommand(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "info")
	out, err := cmd.Output()
	require.NoError(t, err)

	var metadata map[string]any
	require.NoError(t, json.Unmarshal(out, &metadata))
	assert.Equal(t, "weaver-mustache", metadata["name"])
	assert.Equal(t, "mustache", metadata["engine"])
	assert.Equal(t, "1.0.0", metadata["version"])
}

func TestRenderCommand_JSONStdinToStdout(t *testing.T) {
	bin := buildBinary(t)

	tmplDir := t.TempDir()
	tmplPath := filepath.Join(tmplDir, "test.mustache")
	require.NoError(t, os.WriteFile(tmplPath, []byte("Hello {{name}}"), 0644))

	cmd := exec.Command(bin, "render", "-t", tmplPath, "-f", "-")
	cmd.Stdin = bytes.NewBufferString(`{"name": "world"}`)
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "Hello world", string(out))
}

func TestRenderCommand_YAMLFileToStdout(t *testing.T) {
	bin := buildBinary(t)

	tmplDir := t.TempDir()
	tmplPath := filepath.Join(tmplDir, "test.mustache")
	require.NoError(t, os.WriteFile(tmplPath, []byte("port={{port}}"), 0644))

	dataPath := filepath.Join(tmplDir, "data.yaml")
	require.NoError(t, os.WriteFile(dataPath, []byte("port: 8080"), 0644))

	cmd := exec.Command(bin, "render", "-t", tmplPath, "-f", dataPath)
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "port=8080", string(out))
}

func TestRenderCommand_ToOutputFile(t *testing.T) {
	bin := buildBinary(t)

	tmplDir := t.TempDir()
	tmplPath := filepath.Join(tmplDir, "test.mustache")
	require.NoError(t, os.WriteFile(tmplPath, []byte("result={{value}}"), 0644))

	outPath := filepath.Join(tmplDir, "out.txt")

	cmd := exec.Command(bin, "render", "-t", tmplPath, "-f", "-", "-o", outPath)
	cmd.Stdin = bytes.NewBufferString(`{"value": 42}`)
	require.NoError(t, cmd.Run())

	content, err := os.ReadFile(outPath)
	require.NoError(t, err)
	assert.Equal(t, "result=42", string(content))
}

func TestRenderCommand_Sections(t *testing.T) {
	bin := buildBinary(t)

	tmplDir := t.TempDir()
	tmplPath := filepath.Join(tmplDir, "test.mustache")
	require.NoError(t, os.WriteFile(tmplPath, []byte("{{#items}}{{name}} {{/items}}"), 0644))

	cmd := exec.Command(bin, "render", "-t", tmplPath, "-f", "-")
	cmd.Stdin = bytes.NewBufferString(`{"items": [{"name": "a"}, {"name": "b"}]}`)
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Equal(t, "a b ", string(out))
}

func TestRenderCommand_MissingTemplate(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "render", "-t", "/nonexistent.mustache", "-f", "-")
	cmd.Stdin = bytes.NewBufferString(`{}`)
	err := cmd.Run()
	assert.Error(t, err)
}

func TestVersionFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "--version")
	out, err := cmd.Output()
	require.NoError(t, err)
	assert.Contains(t, string(out), "1.0.0")
}
