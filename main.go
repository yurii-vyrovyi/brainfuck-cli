package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/yurii-vyrovyi/brainfuck"
	bfreader "github.com/yurii-vyrovyi/brainfuck/reader"
	bfwriter "github.com/yurii-vyrovyi/brainfuck/writer"

	"golang.org/x/exp/constraints"
)

type Config struct {
	cmdFileName string
	dataSize    int
	input       string
	output      string
}

const (
	InputStdIn   = "stdin"
	OutputStdOut = "stdout"

	HelpString = `brainfuck-cli tool is a brainfuck language interpreter
	
  Usage:
    brainfuck-cli -f=CMD_FILENAME [-s=MEMORY_SIZE] [-i=INPUT] [-o=OUTPUT]

    CMD_FILENAME    File with brainfuck commands

    MEMORY_SIZE     [optional] Size of brainfuck interpreter memory. If not set brainfuck package uses default value.

    INPUT           [optional] Specifies the way to get input. Possible options are:
                      - 'stdin' – interpreter reads input from StdIn. This is default value.
                      - filename – File that contains input that will be read

    OUTPUT          [optional] Specifies the way to write output. Possible options are:
                      - 'stdout' – interpreter writes output to StdOut. This is default value.
                      - filename – File that will get output values
`
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

	config, err := parseParams()
	if err != nil {

		fmt.Println("failed to parse parameters: ", err)
		fmt.Println()
		fmt.Println(HelpString)

		os.Exit(1)
	}

	if err := run(config); err != nil {
		os.Exit(1)
	}

	if err := run(config); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func run(config *Config) error {

	cmdFile, err := os.Open(config.cmdFileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = cmdFile.Close() }()

	inReader, err := createInputReader[int16](config)
	if err != nil {
		return err
	}
	defer func() { _ = inReader.Close() }()

	outWriter, err := createOutputWriter[int16](config)
	if err != nil {
		return err
	}
	defer func() { _ = outWriter.Close() }()

	bfInterpreter := brainfuck.New[int16](config.dataSize, inReader, outWriter).
		WithCmd('^', func(bf *brainfuck.BfInterpreter[int16]) error {
			bf.Data[bf.DataPtr] = bf.Data[bf.DataPtr] * bf.Data[bf.DataPtr]
			return nil
		})

	res, err := bfInterpreter.Run(cmdFile)
	if err != nil {
		return fmt.Errorf("interpreter run failed: %w", err)
	}

	fmt.Println()
	fmt.Println("----- MEMORY DUMP -----\r")
	for i, v := range res {
		fmt.Printf("%d", v)
		if i < len(res)-1 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	fmt.Println()

	return nil
}

func createInputReader[DataType constraints.Signed](config *Config) (brainfuck.InputReader[DataType], error) {
	var inReader brainfuck.InputReader[DataType]
	var err error

	if config.input == InputStdIn {
		inReader, err = bfreader.BuildStdInReader[DataType]()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdin input reader: %w", err)
		}
	} else {
		inReader, err = bfreader.BuildFileReader[DataType](config.input)
		if err != nil {
			return nil, fmt.Errorf("failed to create file input reader: %w", err)
		}
	}
	return inReader, nil
}

func createOutputWriter[DataType constraints.Signed](config *Config) (brainfuck.OutputWriter[DataType], error) {
	var outWriter brainfuck.OutputWriter[DataType]
	var err error

	if config.output == OutputStdOut {
		outWriter = bfwriter.BuildStdOutWriter[DataType]()
		if err != nil {
			return nil, fmt.Errorf("failed to create stout output writer: %w", err)
		}
	} else {
		outWriter, err = bfwriter.BuildFileWriter[DataType](config.output)
		if err != nil {
			return nil, fmt.Errorf("failed to create file output writer: %w", err)
		}
	}
	return outWriter, nil
}

func parseParams() (*Config, error) {

	if len(os.Args) < 2 {
		return nil, errors.New("not enough arguments")
	}

	config := Config{}

	const (
		delim     = "="
		argPrefix = "-"
	)

	for iArg := 1; iArg < len(os.Args); iArg++ {

		delimPos := strings.Index(os.Args[iArg], delim)
		if delimPos == -1 || delimPos == 0 || delimPos == len(os.Args[iArg])-1 {
			continue
		}

		k := os.Args[iArg][:delimPos]
		v := os.Args[iArg][delimPos+1:]

		if !strings.HasPrefix(k, argPrefix) {
			continue
		}

		switch k {

		case "-f":
			config.cmdFileName = v

		case "-s":
			s := v
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad data size parameter [%s]", s)
			}

			config.dataSize = int(n)

		case "-i":
			config.input = v

		case "-o":
			config.output = v
		}
	}

	// checking required parameters
	if len(config.cmdFileName) == 0 {
		return nil, errors.New("filename is empty")
	}

	if len(config.input) == 0 || strings.ToLower(config.input) == InputStdIn {
		config.input = InputStdIn
	}

	if len(config.output) == 0 || strings.ToLower(config.output) == OutputStdOut {
		config.output = OutputStdOut
	}

	return &config, nil
}
