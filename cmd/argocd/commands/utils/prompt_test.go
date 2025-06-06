package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrompt_PromptsEnabled_True(t *testing.T) {
	prompt := NewPrompt(true)
	assert.True(t, prompt.enabled)
}

func TestNewPrompt_PromptsEnabled_False(t *testing.T) {
	prompt := NewPrompt(false)
	assert.False(t, prompt.enabled)
}

func TestConfirm_PromptsEnabled_False(t *testing.T) {
	prompt := NewPrompt(false)
	assert.True(t, prompt.Confirm("Are you sure you want to run this command? (y/n) "))
}

// Returns true, true when prompt is disabled
func TestConfirmAllPromptDisabled(t *testing.T) {
	p := &Prompt{enabled: false}
	result1, result2 := p.ConfirmAll("Proceed?")
	assert.True(t, result1)
	assert.True(t, result2)
}

func TestConfirmBaseOnCountPromptDisabled(t *testing.T) {
	p := &Prompt{enabled: false}
	result1, result2 := p.ConfirmBaseOnCount("Proceed?", "Process all?", 2)
	assert.True(t, result1)
	assert.True(t, result2)
}

func TestConfirmBaseOnCount(t *testing.T) {
	tests := []struct {
		input  string
		output bool
		count  int
	}{
		{
			input:  "y\n",
			output: true,
			count:  0,
		},
		{
			input:  "y\n",
			output: true,
			count:  1,
		},
		{
			input:  "n\n",
			output: false,
			count:  1,
		},
	}

	origStdin := os.Stdin

	for _, tt := range tests {
		tmpFile, err := writeToStdin(tt.input)
		require.NoError(t, err)
		p := &Prompt{enabled: true}
		result1, result2 := p.ConfirmBaseOnCount("Proceed?", "Proceed all?", tt.count)
		assert.Equal(t, tt.output, result1)
		if tt.count == 1 {
			assert.False(t, result2)
		} else {
			assert.Equal(t, tt.output, result2)
		}
		_ = tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	os.Stdin = origStdin
}

func TestConfirmPrompt(t *testing.T) {
	cases := []struct {
		input  string
		output bool
	}{
		{"y\n", true},
		{"n\n", false},
	}

	origStdin := os.Stdin

	for _, c := range cases {
		tmpFile, err := writeToStdin(c.input)
		require.NoError(t, err)
		p := &Prompt{enabled: true}
		result := p.Confirm("Are you sure you want to run this command? (y/n) \n")
		assert.Equal(t, c.output, result)
		_ = tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	os.Stdin = origStdin
}

func TestConfirmAllPrompt(t *testing.T) {
	cases := []struct {
		input      string
		confirm    bool
		confirmAll bool
	}{
		{"y\n", true, false},
		{"n\n", false, false},
		{"a\n", true, true},
	}

	origStdin := os.Stdin

	for _, c := range cases {
		tmpFile, err := writeToStdin(c.input)
		require.NoError(t, err)
		p := &Prompt{enabled: true}
		confirm, confirmAll := p.ConfirmAll("Are you sure you want to run this command? (y/n) \n")
		assert.Equal(t, c.confirm, confirm)
		assert.Equal(t, c.confirmAll, confirmAll)
		_ = tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	os.Stdin = origStdin
}

func writeToStdin(msg string) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "test-input")
	if err != nil {
		return nil, err
	}

	// Write the input to the temporary file
	if _, err := tmpFile.WriteString(msg); err != nil {
		return nil, err
	}

	// Seek to the beginning of the file so it can be read from the start
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, err
	}

	os.Stdin = tmpFile

	return tmpFile, nil
}
