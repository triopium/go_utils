package helper

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"testing"

	"github.com/triopium/go_utils/pkg/logging"
)

// TesterConfig
// TODO: maybe add t.Error() to recover
type TesterConfig struct {
	// Config
	TestDataSource string
	TempDirName    string
	TempDir        string
	TempDataSource string
	TempDataOutput string

	// Internals
	currentDir      string
	testType        string
	initializedTemp bool
	initializedMain bool
	failed          bool
	sigChan         chan os.Signal
	WaitCount       int
	WaitGroup       *sync.WaitGroup
}

// TesterMain
func (tc *TesterConfig) TesterMain(m *testing.M) {
	tc.InitMain()
	exitCode := m.Run()
	slog.Debug("exit code", "code", exitCode)
	tc.WaitGroup.Wait()
	tc.CleanuUP()
}

// WaitAdd
func (tc *TesterConfig) WaitAdd() {
	tc.WaitCount++
	tc.WaitGroup.Add(1)
	slog.Debug("wait count", "count", tc.WaitCount)
}

// WaitDone
func (tc *TesterConfig) WaitDone() {
	tc.WaitGroup.Done()
	tc.WaitCount--
	slog.Debug("wait count", "count", tc.WaitCount)
}

// InitMain
func (tc *TesterConfig) InitMain() {
	if !tc.initializedMain {
		tc.initializedMain = true
		level := os.Getenv("GOLOGLEVEL")
		curDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		tc.currentDir = curDir
		logging.SetLogLevel(level, "json")
		tc.testType = os.Getenv("GO_TEST_TYPE")
		flag.Parse()
		slog.Debug("test config initialized")
		tc.sigChan = make(chan os.Signal, 1)
		tc.WaitGroup = new(sync.WaitGroup)
		signal.Notify(
			tc.sigChan,
			syscall.SIGILL,
			syscall.SIGINT,
			syscall.SIGHUP,
		)
		if tc.testType == "manual" {
			tc.WaitAdd()
		}
		go tc.WaitForSignal()
	}
}

// WaitForSignal
func (tc *TesterConfig) WaitForSignal() {
	slog.Debug("waiting for signal")
	sig := <-tc.sigChan
	slog.Debug("interrupting", "signal", sig.String())
	switch sig {
	case syscall.SIGINT:
		<-tc.sigChan
		if tc.testType != "manual" {
			tc.CleanuUP()
			os.Exit(-1)
		}
	case syscall.SIGILL:
		slog.Error("bad instruction")
		if tc.testType == "manual" {
			slog.Error("bad instruction, waiting", "count", tc.WaitCount)
			<-tc.sigChan
		}
	case syscall.SIGHUP:
		slog.Info("test ends")
	}
	tc.WaitDone()
}

// InitTempSrc
func (tc *TesterConfig) InitTempSrc(
	testSubdir ...string) {
	if len(testSubdir) == 0 {
		return
	}
	if !tc.initializedTemp {
		tc.TempDir = DirectoryCreateTemporaryOrPanic(tc.TempDirName)
		packageName := filepath.Base(tc.currentDir)
		tc.TempDataSource = filepath.Join(tc.TempDir, packageName, "SRC")
		tc.TempDataOutput = filepath.Join(tc.TempDir, packageName, "DST")
		tc.initializedTemp = true
	}
	createDir := func(folder string) {
		err := os.MkdirAll(folder, 0700)
		if err != nil {
			err1 := fmt.Errorf(
				"cannot crete output directory: %s, err: %w",
				tc.TempDataOutput, err)
			panic(err1)
		}
	}
	// Crete DST temp data output
	createDir(tc.TempDataOutput)
	for _, s := range testSubdir {
		if s == "" {
			panic(
				"empty string passed as test subdirectory is not allowed as safety measure")
		}
		srcDir := filepath.Join(tc.TestDataSource, s)
		dstDir := filepath.Join(tc.TempDataSource, s)
		ok, err1 := DirectoryExists(srcDir)

		if err1 != nil || !ok {
			err := fmt.Errorf(
				"err: test data source dir path: %s, exists: %t, err1: %w", srcDir, ok, err1)
			panic(err)
		}
		ok, err2 := DirectoryExists(dstDir)
		if err2 != nil {
			err := fmt.Errorf(
				"test data source dir path: %s exists %t, err: %w", srcDir, ok, err1)
			panic(err)
		}
		if ok {
			// temp test subdir already exists
			continue
		}
		err_copy := DirectoryCopy(
			srcDir, dstDir,
			true, false, "", false,
		)
		if err_copy != nil {
			err := fmt.Errorf("cannot copy directory from %s to %s", srcDir, dstDir)
			panic(err)
		}
		outputSubDir := filepath.Join(tc.TempDataOutput, s)
		createDir(outputSubDir)
	}
}

// CleanuUP
func (tc *TesterConfig) CleanuUP() {
	if tc.initializedTemp {
		DirectoryDeleteOrPanic(tc.TempDir)
	}
}

// TempSourcePathGeter
func (tc *TesterConfig) TempSourcePathGeter(tempSubdir string) func(string) string {
	return func(relPath string) string {
		// var rel string
		// if relPath == "" || relPath == "." {
		// rel = string(os.PathSeparator)
		// }
		return filepath.Join(
			// tc.TempDataSource, tempSubdir, rel)
			tc.TempDataSource, tempSubdir, relPath)
	}
}

// TempDestinationPathGeter
func (tc *TesterConfig) TempDestinationPathGeter(tempSubdir string) func(string) string {
	return func(relPath string) string {
		return filepath.Join(
			tc.TempDataOutput, tempSubdir, relPath)
	}
}

// PrintResult
func (tc *TesterConfig) PrintResult(a ...any) {
	if tc.testType == "manual" {
		fmt.Println(a...)
	}
}

// InitTest
func (tc *TesterConfig) InitTest(
	t *testing.T, testSubdir ...string) {
	if tc.failed {
		t.SkipNow()
		return
	}
	if testing.Short() && len(testSubdir) == 0 {
		t.SkipNow()
		return
	}
	tc.WaitAdd()
	tc.InitTempSrc(testSubdir...)
	slog.Debug("test initialized", "name", t.Name())
}

// RecoverPanic
func (tc *TesterConfig) RecoverPanic(t *testing.T) {
	// TODO: maybe add t.Error(err)
	if t.Skipped() {
		return
	}
	if r := recover(); r != nil {
		tc.failed = true
		slog.Error("test panics", "reason", r)
		t.Fail()
		tc.WaitDone()
		if tc.testType == "manual" {
			tc.sigChan <- syscall.SIGILL
		}
		return
	}
	if !tc.failed {
		tc.WaitDone()
	}
}

func (tc *TesterConfig) RecoverPanicNoFail(t *testing.T) {
	if t.Skipped() {
		return
	}
	if r := recover(); r != nil {
		slog.Info("test recovered panic", "reason", r)
		tc.WaitDone()
		if tc.testType == "manual" {
			tc.sigChan <- syscall.SIGILL
		}
		return
	}
	if !tc.failed {
		tc.WaitDone()
	}
}
