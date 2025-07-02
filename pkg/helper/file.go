package helper

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	// _ "unicode"

	enc_unicode "golang.org/x/text/encoding/unicode"
)

var ErrFilePathExists = errors.New("file path exists")
var ErrFileEncodingUnknown = errors.New("file encoding name is unknown")

// CharEncoding specifies character encoding in file.
type CharEncoding string

const (
	CharEncodingUTF8    CharEncoding = "UTF8"
	CharEncodingUTF16le CharEncoding = "UTF16le"
	CharEncodingUTF16be CharEncoding = "UTF16be"
)

type EncodingHandler struct {
	// ReaderFunc func(io.ReadCloser) io.ReadCloser
	ReaderFunc func(io.Reader) io.Reader
}

var CharEncodingMap = map[CharEncoding]EncodingHandler{
	CharEncodingUTF8:    {FileGetReaderUTF8},
	CharEncodingUTF16le: {FileGetReaderUTF16le},
}

func FileOsEncoding(file *os.File) (CharEncoding, error) {
	// Read the first few bytes to detect the BOM
	buf := make([]byte, 4)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Reset the file read pointer
	_, err = file.Seek(0, io.SeekStart) // nolint:errcheck
	if err != nil {
		return "", err
	}

	var enc CharEncoding
	switch {
	case bytes.HasPrefix(buf[:n], []byte{0xFE, 0xFF}):
		enc = CharEncodingUTF16be
		// enc = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	case bytes.HasPrefix(buf[:n], []byte{0xFF, 0xFE}):
		enc = CharEncodingUTF16le
		// enc = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	default:
		enc = CharEncodingUTF8
	}
	// enc.NewDecoder()
	return enc, nil
}

func FilePathEncoding(filePath string) (CharEncoding, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return FileOsEncoding(file)
}

func FileReaderHandleEncoding(filePath string) (
	io.Reader, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	enc, err := FileOsEncoding(file)
	if err != nil {
		file.Close()
		return nil, nil, err
	}
	reader := CharEncodingMap[enc]
	return reader.ReaderFunc(file), file, nil
}

func FileReadAllHandleEncoding(filePath string) (
	[]byte, CharEncoding, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	enc, err := FileOsEncoding(file)
	if err != nil {
		return nil, "", err
	}
	reader := CharEncodingMap[enc]
	data, err := io.ReadAll(reader.ReaderFunc(file))
	if err != nil {
		return nil, "", err
	}
	return data, enc, err
}

func FileGetReaderUTF8(reader io.Reader) io.Reader {
	// func FileGetReaderUTF8(ioReaderCloser io.ReadCloser) io.ReadCloser {
	return reader
}

func FileGetReaderUTF16le(reader io.Reader) io.Reader {
	// func FileGetReaderUTF16le(reader io.ReadCloser) io.ReadCloser {
	utf8reader := enc_unicode.UTF16(
		enc_unicode.LittleEndian,
		enc_unicode.IgnoreBOM).NewDecoder().Reader(reader)
	return utf8reader
}

// HandleFileEncoding read file according to specified file encoding code.
func HandleFileEncoding(
	enc CharEncoding, ioReaderCloser io.ReadCloser) ([]byte, error) {
	var data []byte
	var err error
	switch enc {
	case CharEncodingUTF8:
		data, err = io.ReadAll(ioReaderCloser)
	case CharEncodingUTF16le:
		utf8reader := enc_unicode.UTF16(
			enc_unicode.LittleEndian,
			enc_unicode.IgnoreBOM).NewDecoder().Reader(ioReaderCloser)
		data, err = io.ReadAll(utf8reader)
	default:
		err = fmt.Errorf("unknown encoding")
	}
	return data, err
}

// // HandleFileEncoding read file according to specified file encoding code.
// func HandleFileEncodingB(
// 	enc FileEncodingCode, ioReaderCloser io.ReadCloser) ([]byte, error) {
// 	var data []byte
// 	var err error
// 	switch enc {
// 	case UTF8:
// 		data, err = io.ReadAll(ioReaderCloser)
// 	case UTF16le:
// 		utf8reader := enc_unicode.UTF16(enc_unicode.LittleEndian, enc_unicode.IgnoreBOM).NewDecoder().Reader(ioReaderCloser)
// 		data, err = io.ReadAll(utf8reader)
// 	default:
// 		err = fmt.Errorf("unknown encoding")
// 	}
// 	return data, err
// }

// CopyFile copies file from source path to destination path.
func CopyFile(
	srcFilePath, dstFilePath string,
	overwrite bool, verbose bool,
) error {
	if verbose {
		slog.Debug(
			"copying file",
			"source_file", srcFilePath,
			"dst_file", dstFilePath,
		)
	}
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file in the destination directory for writing
	pathExists, err := PathExists(dstFilePath)
	if err != nil {
		return err
	}
	if !overwrite && pathExists {
		return fmt.Errorf(
			"err: %w, filepath: %s",
			ErrFilePathExists, dstFilePath,
		)
	}
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

// FileExists: returns true if file exists, false when the filePath doesnot exists, error when it is directory
func FileExists(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if fileInfo.IsDir() {
			return false, fmt.Errorf("specified path is a directory, not a file: %s", filePath)
		}
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// ProcessedFileRename adds prfix to original filepath.
func ProcessedFileRename(originalPath string) error {
	fileName := filepath.Base(originalPath)
	directory := filepath.Dir(originalPath)
	newPath := filepath.Join(directory, "processed_"+fileName)
	err := os.Rename(originalPath, newPath)
	if err != nil {
		return fmt.Errorf("Error renaming file: %s", err)
	}
	return nil
}

// ReadCSVfile
func ReadCSVfile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	// reads all the records from the CSV file and return [][]string
	return reader.ReadAll()
}

func ReadCSVfileSep(filePath string, sep rune) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.LazyQuotes = true
	// reads all the records from the CSV file and return [][]string
	return reader.ReadAll()
}

// GetFilenameWithoutExtension
func FilenameWithoutExtension(filePath string) string {
	// Get the base name of the file
	base := filepath.Base(filePath)
	// Get the file extension
	ext := filepath.Ext(base)
	// Remove the extension from the base name
	return strings.TrimSuffix(base, ext)
}
