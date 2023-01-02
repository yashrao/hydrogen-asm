# Hydrogen - An Assembly interpreter

## What is this for?
This is my attempt at trying to build a virtual machine that uses intel flavoured GCC assembly for use with another project. The entire idea is you can compile into assembly and Hydrogen can be used to run the program directly.

## Features
I just started so it barely does anything
- Supports MOV, ADD, SUB, MUL, DIV
- Has an iPython like interactive mode with the -i flag - not entirely sure if this is useful but we'll see how it goes

## Build
```bash
git clone https://github.com/yashrao/hydrogen-asm
cd hydrogen-asm
go build -o hydrogen main.go
```

## Usage
```txt
$ hydrogen --help
OVERVIEW: Hydrogen NASM interpreter

Usage: ./hydrogen [-i] [OPTIONS] [FILENAME]
OPTIONS:
	-h, --help		Display this Message
        -i, --interactive	Run in interactive mode
```
