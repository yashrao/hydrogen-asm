package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	// Debug
	//"reflect"
)

// OPCODES CATEGORY
const (
	OP_MOV     = 1
	OP_ADD     = 2
	OP_RET     = 3
	OP_INVALID = 4
)

// Token Types
const (
	TK_INSTRUCTION = 10
	TK_CONSTANT    = 11
	TK_REGISTER    = 12
	TK_NEWLINE     = 13
	//TODO
	TK_PUSH        = 14
	TK_POP         = 15

	TK_INVALID     = -1
)

// Colors for error reporting
const (
	RED = "\033[0;31m" 
	GREEN = "\033[0;32m" 
	YELLOW = "\033[0;33m"
    BLUE = "\033[0;34m"
	RESET = "\033[0m" // Used to reset
)

var stack []int 

type Token struct {
	tk_type int
	str_repr string
}

type Registers struct {
	// 64 bit registers
	rax int64
	rcx int64
	rdx int64
	rsi int64
	rdi int64
	r8  int64
	r9  int64
	r10 int64
	r11 int64
	/*preserved*/
	rbx int64
	rsp int64
	rbp int64
	r12 int64
	r13 int64
	r14 int64
	r15 int64
	/**********/

	// 32 bit registers
    eax int32
    ecx int32
    edx int32
    esi int32
    edi int32
	/*preserved*/
	esp int32
	ebp int32
    ebx int32
	/**********/

	// 16 bit registers
	ax int16
	cx int16
	dx int16
	si int16
	di int16
	/*preserved*/
	bx int16
	sp int16
	bp int16
	/**********/
}

func init_regs(registers *Registers) {
	// 64 bit registers
	registers.rax = 0
	registers.rcx = 0
	registers.rdx = 0
	registers.rbx = 0
	registers.rsi = 0
	registers.rdi = 0
	registers.r8  = 0
	registers.r9  = 0
	registers.r10 = 0
	registers.r11 = 0
	registers.r12 = 0
	registers.r13 = 0
	registers.r14 = 0
	registers.r15 = 0

	// Gotta check about these two
	registers.rsp = 0 
	registers.rbp = 0
}

func print_error(msg string) {
	fmt.Printf("%s[ERROR]:%s %s\n", RED, RESET, msg)
}

func get_op_code(word string) int {
	if word == "mov" || word == "add" || word == "sub" || word == "mul" {
		return TK_INSTRUCTION 
	} else if word == "div" {
		return TK_INSTRUCTION
	} else if is_register(word) {
		return TK_REGISTER
	} else if is_digit(word) {
		return TK_CONSTANT
	} else {
		return TK_INVALID 
	}
}

func print_tokens(token_list []Token) {
	for _, token := range token_list {
		print_token(token)
	}
	/*for _, token := range token_list {
		fmt.Printf("TOKEN: %s\n", token)
	}*/
}

func print_token(token Token) {
	fmt.Printf("TOKEN: \"%s\"<", token.str_repr)
	if token.tk_type == TK_INSTRUCTION {
		fmt.Printf("TK_INSTRUCTION>")
	} else if token.tk_type == TK_CONSTANT {
		fmt.Printf("TK_CONSTANT>")
	} else if token.tk_type == TK_NEWLINE {
		fmt.Printf("TK_NEWLINE>")
	} else if token.tk_type == TK_REGISTER {
		fmt.Printf("TK_REGISTER>")
	} else {
		fmt.Printf("TK_INVALID>")
	}
	fmt.Printf("\n")
}

func convert_to_int64(str_repr string) int64 {
	value, err := strconv.Atoi(str_repr)
	if err != nil {
		fmt.Println("Conversion did not work")
	}
	return int64(value)
}

func get_register(registers *Registers, str_repr string) *int64 {
	switch str_repr {
	case "rax":
		return &registers.rax
	case "rcx":
		return &registers.rcx
    case "rdx":
		return &registers.rdx
    case "rbx:":
		return &registers.rbx
    case "rsi":
		return &registers.rsi
    case "rdi":
		return &registers.rdi
    case "r8":
		return &registers.r8
    case "r9":
		return &registers.r9
    case "r10":
		return &registers.r10
    case "r11":
		return &registers.r11
    case "r12":
		return &registers.r12
    case "r13":
		return &registers.r13
    case "r14":
		return &registers.r14
    case "r15":
		return &registers.r15
    case "rsp":
		return &registers.rsp
    case "rbp":
		return &registers.rbp
	}
	return nil
}

func print_registers(registers *Registers) {
	fmt.Printf("%s%s:%s %d\n", BLUE, "rax", RESET, registers.rax)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rcx", RESET, registers.rcx)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rdx", RESET, registers.rdx)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rbx", RESET, registers.rbx)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rsi", RESET, registers.rsi)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rdi", RESET, registers.rdi)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r8 ", RESET, registers.r8)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r9 ", RESET, registers.r9)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r10", RESET, registers.r10)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r11", RESET, registers.r11)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r12", RESET, registers.r12)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r13", RESET, registers.r13)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r14", RESET, registers.r14)
	fmt.Printf("%s%s:%s %d\n", BLUE, "r15", RESET, registers.r15)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rsp", RESET, registers.rsp)
	fmt.Printf("%s%s:%s %d\n", BLUE, "rbp", RESET, registers.rbp)
}

func execute_instruction(registers *Registers, instruction Token, lhs Token, rhs Token) {
	if rhs.tk_type == TK_CONSTANT {
		if instruction.str_repr == "mov" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register = convert_to_int64(rhs.str_repr)
		} else if instruction.str_repr == "add" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register += convert_to_int64(rhs.str_repr)
		} else if instruction.str_repr == "sub" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register -= convert_to_int64(rhs.str_repr)
		} else if instruction.str_repr == "mul" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register *= convert_to_int64(rhs.str_repr)
		} else if instruction.str_repr == "div" {
			var register *int64 = get_register(registers, lhs.str_repr)
			// Do we want to handle this? Go will handle and crash by itself already lol
			*register /= convert_to_int64(rhs.str_repr)
		}
	} else if rhs.tk_type == TK_REGISTER {
		if instruction.str_repr == "mov" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register = *get_register(registers, rhs.str_repr)
		} else if instruction.str_repr == "add" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register += *get_register(registers, rhs.str_repr)
		} else if instruction.str_repr == "sub" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register -= *get_register(registers, rhs.str_repr)
		} else if instruction.str_repr == "mul" {
			var register *int64 = get_register(registers, lhs.str_repr)
			*register *= *get_register(registers, rhs.str_repr)
		} else if instruction.str_repr == "div" {
			var register *int64 = get_register(registers, lhs.str_repr)
			// Do we want to handle this? Go will handle and crash by itself already lol
			*register /= *get_register(registers, rhs.str_repr)
		}

	}
}

func execute_tokens(tokens []Token, registers *Registers) bool {
	success := false
	for i := 0; i < len(tokens); i++ {
		if tokens[i].tk_type == TK_INSTRUCTION {
			instruction := tokens[i]
			// Expecting two operands
			i++ // register
			lhs := tokens[i]
			if lhs.tk_type != TK_REGISTER {
				print_error(fmt.Sprintf("Expected register, got: \"%s\"", tokens[i].str_repr))
				return false
			}
			i++
			rhs := tokens[i]
			if rhs.tk_type == TK_REGISTER {
				//TODO:
				execute_instruction(registers, instruction, lhs, rhs)
				success = true
			} else if rhs.tk_type == TK_REGISTER || rhs.tk_type == TK_CONSTANT {
				execute_instruction(registers, instruction, lhs, rhs)
				success = true
			} else {
				print_error(fmt.Sprintf("Expected Constant or Register, got: \"%s\"", tokens[i].str_repr))
			}
		}
	}
	return success
}

func is_digit(word string) bool {
	_, err := strconv.Atoi(word)
	if err != nil {
		return false
	}
	return true
}

func is_register(word string) bool {
	switch str := word; str {
	case "rax":
		return true
	case "rcx":
		return true
    case "rdx":
		return true
    case "rbx:":
		return true
    case "rsi":
		return true
    case "rdi":
		return true
    case "r8":
		return true
    case "r9":
		return true
    case "r10":
		return true
    case "r11":
		return true
    case "r12":
		return true
    case "r13":
		return true
    case "r14":
		return true
    case "r15":
		return true
    case "rsp":
		return true
    case "rbp":
		return true
	}
	return false
}

func parse_line(line string) []Token {
	//words := strings.Split(line, " ")
	ret := []Token{}
	words := strings.Fields(line)
	if len(words) > 3 {
		print_error(fmt.Sprintf("Expected instruction with 3 words or fewer, got: %d words ", len(words)))
	} else if len(words) == 1 {
		fmt.Println("ONE WORD - NOT IMPLEMENTED YET")
	} else if len(words) == 2 {
		fmt.Println("TWO WORDS - NOT IMPLEMENTED YET")
	} else {
		// VALID INSTRUCTION
		for i, word := range words {
			// Getting instruction
			op := get_op_code(word)
			if i == 0 {
				// Expecting instruction
				if op == TK_INVALID {
					print_error(fmt.Sprintf("Expected an instruction, got: \"%s\"", word))
					return ret
				}
				word = strings.ToLower(word)
				token := Token{op, word}
				ret = append(ret, token)
			} else {
				// Expecting register or constant
				if is_register(word) {
					ret = append(ret, Token{TK_REGISTER, word})
				} else {
					ret = append(ret, Token{get_op_code(word), word})
				}
			}
		}
		ret = append(ret, Token{TK_NEWLINE, "NEWLINE"})
	}
	print_tokens(ret)
	return ret
}

func main_interactive_loop() {
    fmt.Println("Hydrogen interpreter v0 on Linux")
	i := 1
	var line string
	var registers Registers
	var history []string
    _ = history // TODO: Store history of input
	init_regs(&registers)
	stack := []int{} // Initialize the stack - TODO: make it accept other things as well
	_ = stack
	for {
		fmt.Printf("%s[%d]:%s ", GREEN, i, RESET)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			line = scanner.Text()
		}
		tokens := parse_line(line)
		if len(tokens) != 0 {
			success := execute_tokens(tokens, &registers)
			if success {
				print_registers(&registers)
			}
		}

		i += 1
		fmt.Println("")
	}
}

func main_loop_backend() {
	//TODO: for integration with helium
}

func main_loop(filename string) {
	
}

func print_help(prog_name string) {
	//TODO: when other options available use this
	//fmt.Printf("Usage: %s [-i:-c] [OPTIONS] INPUT\n")
	fmt.Printf("OVERVIEW: Hydrogen NASM interpreter\n\n")
	fmt.Printf("Usage: %s [-i] [OPTIONS] [FILENAME]\n", prog_name)

	fmt.Printf("OPTIONS:\n\t-h, --help\t\tDisplay this Message\n")
	fmt.Printf("        -i, --interactive\tRun in interactive mode\n")
	fmt.Printf("\n")
}

func str_in_arr(str string, array []string) bool {
	for _, arg := range array {
		if str == arg {
			return true
		}
	}
	return false
}

func str_not_in_arr(str string, array []string) bool {
	return !str_in_arr(str, array)
}

func main() {
	args := os.Args
	fmt.Println(args)
	if str_in_arr("-h", args) || str_in_arr("--help", args) {
		print_help(args[0])
	}
	if str_in_arr("-i", args) || str_in_arr("--interactive", args) {
		main_interactive_loop()
	} else if (len(args) > 1 && str_not_in_arr("-", args)) {
		main_loop(args[1])
	} else {
		main_loop_backend()
	}
}
