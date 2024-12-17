package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	io "github.com/faideww/aoc-2024/lib"
)

type CPU struct {
	a, b, c int
	program []int
	pc      int
}

func main() {
	input := io.ReadInputFile(os.Args[1])

	cpu := parseCPU(input)

	output := runCPUToHalt(&cpu)

	outputStrs := make([]string, len(output))
	for i, o := range output {
		outputStrs[i] = fmt.Sprintf("%d", o)
	}
	fmt.Printf("output: %s\n", strings.Join(outputStrs, ","))

	digits := []int{}

	value := io.PowInt(8, len(cpu.program)-1)

	// The cpu output follows a pattern whereby each value is individually
	// solvable by testing values in base 8. The general idea is to test the CPU
	// program with a value 8^i + 8^(i-1) * j, where i is the index of the
	// current program value we're looking at, and j is some fraction from 0 to 8
	// (ideally, in some cases with carry values j can be more than 8. In my
	// output I have one j value = 11). We then inspect the ith value in the
	// output for a match, and repeat this until we find one.
	for i := len(cpu.program); i > 0; i-- {
		found := false
		for j := 0; j < 64; j++ {
			initial := value + io.PowInt(8, i-1)*j
			cpu.a = initial
			cpu.b = 0
			cpu.c = 0
			cpu.pc = 0

			output := runCPUToHalt(&cpu)

			if len(output) != len(cpu.program) {
				fmt.Printf("wrong output length (expected %d, got %d)\n", len(cpu.program), len(output))
				panic("output is not the correct length")
			}
			if output[i-1] == cpu.program[i-1] {
				value += io.PowInt(8, i-1) * j
				digits = append(digits, j)
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("no digit found for place %d\n", i)
		}
	}

	fmt.Printf("value: %v\n", value)

	cpu.a = value
	cpu.b = 0
	cpu.c = 0
	cpu.pc = 0

	output = runCPUToHalt(&cpu)
	// fmt.Printf("program:%v\n", cpu.program)
	// fmt.Printf("output:%v\n", output)
	isEqual := len(output) == len(cpu.program)
	for j, val := range output {
		if cpu.program[j] != val {
			isEqual = false
			break
		}
	}
	if isEqual {
		fmt.Printf("program copy found with A=%d\n", value)
	} else {
		fmt.Printf("program and output do not match :(\n")
	}
}

func parseCPU(input string) CPU {
	sections := io.TrimAndSplitBy(input, "\n\n")

	registerRegex := regexp.MustCompile("Register A: (.+)\nRegister B: (.+)\nRegister C: (.+)")

	match := registerRegex.FindStringSubmatch(sections[0])

	aStr, bStr, cStr := match[1], match[2], match[3]

	a, _ := strconv.Atoi(aStr)
	b, _ := strconv.Atoi(bStr)
	c, _ := strconv.Atoi(cStr)

	programRegex := regexp.MustCompile("Program: (.+)")
	progMatch := programRegex.FindStringSubmatch(sections[1])

	programStr := strings.Split(progMatch[1], ",")

	program := make([]int, len(programStr))
	for i, codeStr := range programStr {
		code, _ := strconv.Atoi(codeStr)
		program[i] = code
	}

	return CPU{
		a:       a,
		b:       b,
		c:       c,
		program: program,
		pc:      0,
	}
}

func runCPUToHalt(cpu *CPU) []int {
	output := make([]int, 0)

	for cpu.pc < len(cpu.program) {
		opcode, operand := cpu.program[cpu.pc], cpu.program[cpu.pc+1]

		switch opcode {
		case 0: // adv, combo operand
			dv(cpu, operand, &cpu.a)
		case 1: // bxl, literal operand
			cpu.b = cpu.b ^ operand
		case 2: // bst
			cpu.b = comboOperand(cpu, operand) % 8
		case 3: // jnz
			if cpu.a != 0 {
				cpu.pc = operand - 2
			}
		case 4: // bxc
			cpu.b = cpu.b ^ cpu.c
		case 5: // out
			value := comboOperand(cpu, operand) % 8
			output = append(output, value)
		case 6: // bdv
			dv(cpu, operand, &cpu.b)
		case 7: // cdv
			dv(cpu, operand, &cpu.c)
		}

		cpu.pc += 2
	}

	return output
}

func dv(cpu *CPU, operand int, dest *int) {
	denom := io.PowInt(2, comboOperand(cpu, operand))
	*dest = cpu.a / denom
}

func comboOperand(cpu *CPU, opval int) int {
	operand := opval
	switch operand {
	case 4:
		operand = cpu.a
	case 5:
		operand = cpu.b
	case 6:
		operand = cpu.c
	}

	return operand
}
