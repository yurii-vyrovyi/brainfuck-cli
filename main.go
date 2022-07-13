package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/yurii-vyrovyi/brainfuck/reader"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	bfrunner "github.com/yurii-vyrovyi/brainfuck"
)

type Config struct {
	cmdFileName string
	dataSize    int
	input       string
}

const (
	InputStdin = "stdin"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)

			buf := make([]byte, 4096)
			runtime.Stack(buf, true)
			fmt.Printf("%s\n", buf)
		}
	}()

	if err := setup(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func setup() error {
	config, err := parseParams()
	if err != nil {
		return fmt.Errorf("failed to parse parameters: %w", err)
	}

	if err := run(config); err != nil {
		return err
	}

	return nil
}

func run(config *Config) error {

	cmdFile, err := os.Open(config.cmdFileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = cmdFile.Close() }()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var inReader bfrunner.InputReader

	if config.input == InputStdin {
		inReader, err = reader.BuildStdInReader()
		if err != nil {
			return fmt.Errorf("failed to create stdin input reader: %w", err)
		}
	} else {
		inReader, err = reader.BuildFileReader(config.input)
		if err != nil {
			return fmt.Errorf("failed to create file input reader: %w", err)
		}
	}
	defer func() { _ = inReader.Close() }()

	bfRunner := bfrunner.New(config.dataSize, os.Stdout, inReader)

	buf, err := io.ReadAll(cmdFile)
	if err != nil {
		return fmt.Errorf("faield to read file: %w", err)
	}

	if err := bfRunner.Run(ctx, string(buf)); err != nil {
		return fmt.Errorf("bfRunner run failed: %w", err)
	}

	return nil
}

func parseParams() (*Config, error) {

	if len(os.Args) < 2 {
		return nil, errors.New("not enough arguments")
	}

	config := Config{}
	iParam := 1

	for {

		if iParam >= len(os.Args)-1 {
			break
		}

		switch os.Args[iParam] {

		case "-f":
			config.cmdFileName = os.Args[iParam+1]

		case "-s":
			s := os.Args[iParam+1]
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad data size parameter [%s]", s)
			}

			config.dataSize = int(n)

		case "-i":
			config.input = os.Args[iParam+1]
		}

		iParam = iParam + 2
	}

	// checking required parameters
	if len(config.cmdFileName) == 0 {
		return nil, errors.New("filename is empty")
	}

	if len(config.input) == 0 || strings.ToLower(config.input) == InputStdin {
		config.input = InputStdin
	}

	return &config, nil
}
