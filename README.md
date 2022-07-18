# brainfuck
You can read more about the language itself [here](https://en.wikipedia.org/wiki/Brainfuck).

## brainfuck-cli
`brainfuck-cli` is a commandline brainfuck language interpreter that uses this brainfuck package under the hood
[https://github.com/yurii-vyrovyi/brainfuck](https://github.com/yurii-vyrovyi/brainfuck).

## parameters

```
    brainfuck-cli -f=CMD_FILENAME [-s=MEMORY_SIZE] [-i=INPUT] [-o=OUTPUT]

    CMD_FILENAME    File with brainfuck commands

    MEMORY_SIZE     [optional] Size of brainfuck interpreter memory. If not set brainfuck package uses default value.

    INPUT           [optional] Specifies the way to get input. Possible options are:
                      - 'stdin' – interpreter reads input from StdIn. This is default value.
                      - filename – File that contains input that will be read

    OUTPUT          [optional] Specifies the way to write output. Possible options are:
                      - 'stdout' – interpreter writes output to StdOut. This is default value.
                      - filename – File that will get output values
```