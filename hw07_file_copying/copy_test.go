package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type testData struct {
	tempDir      string
	testFileName string
	testBytes    []byte
	fileSize     int64
}

func createTempDir(t *testing.T) string {
	t.Helper()

	dir := os.TempDir() + "/copy_test"
	err := os.Mkdir(dir, os.FileMode(0o700))
	if err != nil {
		t.Error("create test dir", err)
	}

	return dir
}

func createInputFile(t *testing.T, tempDir string, bytes []byte) string {
	t.Helper()

	f, err := os.CreateTemp(tempDir, "copy_test")
	if err != nil {
		t.Error("create test file", err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		t.Error("write test file", err)
	}
	f.Close()

	return f.Name()
}

func createTestData(t *testing.T, fileSize int64) *testData {
	t.Helper()

	testBytes := make([]byte, fileSize)
	rand.Read(testBytes)

	tempDir := createTempDir(t)
	testFileName := createInputFile(t, tempDir, testBytes)

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return &testData{
		tempDir,
		testFileName,
		testBytes,
		fileSize,
	}
}

func TestCopy(t *testing.T) {
	var fileSize int64 = 256
	testData := createTestData(t, fileSize)

	var cases []struct {
		fromPath, toPath string
		offset, limit    int64
		err              error
	}

	t.Run("success casex", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			randOffset, _ := rand.Int(rand.Reader, big.NewInt(100))
			randLimit, _ := rand.Int(rand.Reader, big.NewInt(fileSize-randOffset.Int64()-100))

			cases = append(cases, struct {
				fromPath, toPath string
				offset, limit    int64
				err              error
			}{
				toPath: testData.tempDir + "/success_" + strconv.Itoa(i),
				offset: randOffset.Int64(),
				limit:  randLimit.Int64() + 100,
			})
		}

		for _, successCase := range cases {
			successCase := successCase

			t.Run(fmt.Sprintf("success case offset %d limit %d", successCase.offset, successCase.limit), func(t *testing.T) {
				err := Copy(testData.testFileName, successCase.toPath, successCase.offset, successCase.limit, io.Discard)
				require.NoError(t, err)

				bytes, err := os.ReadFile(successCase.toPath)
				require.NoError(t, err)

				lowerBound, upperBound := successCase.offset, successCase.offset+successCase.limit
				if successCase.limit == 0 || successCase.limit > fileSize {
					upperBound = testData.fileSize
				}

				require.ElementsMatch(t, bytes, testData.testBytes[lowerBound:upperBound])
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		cases = []struct {
			fromPath, toPath string
			offset, limit    int64
			err              error
		}{
			{
				fromPath: testData.testFileName,
				toPath:   testData.tempDir + "/errorCase",
				offset:   fileSize + 1,
				err:      ErrOffsetExceedsFileSize,
			},
			{
				fromPath: testData.testFileName,
				toPath:   testData.tempDir + "/errorCase",
				offset:   -1,
				err:      ErrWrongParams,
			},
			{
				fromPath: testData.testFileName,
				toPath:   testData.tempDir + "/errorCase",
				limit:    -1,
				err:      ErrWrongParams,
			},
			{
				fromPath: "/dev/urandom",
				toPath:   testData.tempDir + "/errorCase",
				err:      ErrUnsupportedFile,
			},
			{
				fromPath: "/tmp/this_file_is_not_exists",
				toPath:   testData.tempDir + "/errorCase",
				err:      ErrFileIsNotReadable,
			},
			{
				fromPath: testData.testFileName,
				toPath:   "/dev/try_to_write_in_dev",
				err:      ErrFileIsNotWritable,
			},
		}

		for _, errorCase := range cases {
			errorCase := errorCase

			t.Run(fmt.Sprintf("error case %v", errorCase.err), func(t *testing.T) {
				err := Copy(errorCase.fromPath, errorCase.toPath, errorCase.offset, errorCase.limit, io.Discard)

				require.ErrorIs(t, errorCase.err, err)
			})
		}
	})
}
