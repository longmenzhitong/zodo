package stdin

import (
	"bufio"
	"os"
	"strings"
)

func ReadString() string {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(input)
}
