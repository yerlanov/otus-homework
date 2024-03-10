package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	srcContent := "Hello, Otus!\nThis is a test file.\nIt is used for testing purposes.\n"
	temp, err := os.CreateTemp("", "src")
	require.NoError(t, err)
	defer os.Remove(temp.Name())

	_, err = temp.WriteString(srcContent)
	require.NoError(t, err)

	t.Run("TestCopy", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy(temp.Name(), dst.Name(), 0, 0)
		require.NoError(t, err)

		dstContent, err := os.ReadFile(dst.Name())
		require.NoError(t, err)
		require.Equal(t, srcContent, string(dstContent))
	})

	t.Run("TestCopyWithOffset", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy(temp.Name(), dst.Name(), 5, 0)
		require.NoError(t, err)

		dstContent, err := os.ReadFile(dst.Name())
		require.NoError(t, err)
		require.Equal(t, srcContent[5:], string(dstContent))
	})

	t.Run("TestCopyWithLimit", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy(temp.Name(), dst.Name(), 0, 5)
		require.NoError(t, err)

		dstContent, err := os.ReadFile(dst.Name())
		require.NoError(t, err)
		require.Equal(t, srcContent[:5], string(dstContent))
	})

	t.Run("TestCopyWithOffsetAndLimit", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy(temp.Name(), dst.Name(), 5, 5)
		require.NoError(t, err)

		dstContent, err := os.ReadFile(dst.Name())
		require.NoError(t, err)
		require.Equal(t, srcContent[5:10], string(dstContent))
	})

	t.Run("TestCopyWithOffsetExceedsFileSize", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy(temp.Name(), dst.Name(), int64(len(srcContent)+1), 0)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})

	t.Run("TestCopyWithUnsupportedFile", func(t *testing.T) {
		dst, err := os.CreateTemp("", "dst")
		require.NoError(t, err)
		defer os.Remove(dst.Name())

		err = Copy("/dev/urandom", dst.Name(), 0, 0)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})
}
