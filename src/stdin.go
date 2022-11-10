package zodo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readString() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(input)
}

func ReadInt(min int, max int, msg string) (int, error) {
	if min > max {
		panic(fmt.Errorf("min must not bigger than max"))
	}

	if min == max {
		return min, nil
	}

	fmt.Println(msg)
	input := readString()
	if input == "" {
		return -1, &CancelledError{}
	}
	num, err := strconv.Atoi(input)
	if err != nil || num < min || num > max {
		fmt.Printf("Number incorrect, expect [%d ~ %d], got [%s].\n", min, max, input)
		return ReadInt(min, max, msg)
	}
	return num, nil
}

func AskYes(format string, a ...interface{}) bool {
	format += " [y/n]\n"
	fmt.Printf(format, a...)
	return strings.ToLower(readString()) == Yes
}
