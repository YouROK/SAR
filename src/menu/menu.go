package menu

import (
	"fmt"
	"strconv"
)

func List(msgTop, msgBottom string, list []string) int {
	fmt.Println()
	fmt.Println(msgTop)
	for i, l := range list {
		fmt.Println("[", i+1, "]", l)
	}
	if msgBottom == "" {
		msgBottom = "0 для выхода"
	}
	fmt.Println()
	fmt.Print(msgBottom, " ")

	var input string
	fmt.Scanln(&input)
	ret, err := strconv.Atoi(input)
	if err != nil || (ret < 1 && ret > len(list)) {
		return 0
	}
	return ret
}

func Question(msg string) bool {
	fmt.Println()
	fmt.Println(msg)
	fmt.Println("y - да", "n - нет")
	input := Read()
	if input == "y" || input == "Y" {
		return true
	}
	return false
}

func Edit(msg string) string {
	fmt.Println()
	fmt.Println(msg)
	var input string
	fmt.Scanln(&input)
	return input
}

func Text(msg ...interface{}) {
	fmt.Println(msg...)
}

func Read() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func Clear() {
	fmt.Print("\n\n\n\n\n\n\n\n\n\n")
}
