package emulator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func usage(name string, c map[string]Control, t map[string]Control) {
	pUsage, pKey := usagePadding(c, t)
	fmt.Printf("Usage: %s [FILE]\n\n", "chip8")
	fmt.Printf("Emulator Controls:\n")
	fmt.Printf("|" + line(pKey) + "|" + line(pUsage) + "|\n")
	fmt.Printf("| %-*s| %-*s|\n", pKey, "key", pUsage, "function")
	fmt.Printf("|" + line(pKey) + "|" + line(pUsage) + "|\n")
	for key := range c {
		fmt.Printf("| %-*s| %-*s|\n", pKey, key, pUsage, c[key].usage)
	}
	fmt.Printf("|" + line(pKey) + "|" + line(pUsage) + "|\n")
	for key := range t {
		fmt.Printf("| %-*s| %-*s|\n", pKey, key, pUsage, t[key].usage)
	}
	fmt.Printf("|" + line(pKey) + "|" + line(pUsage) + "|\n")
}

func line(amount int) string {
	return strings.Repeat("-", amount+1)
}
func usagePadding(c map[string]Control, t map[string]Control) (int, int) {
	maxLenUsage := 0
	maxLenKey := 0
	for key := range c {
		if maxLenUsage < len(c[key].usage) {
			maxLenUsage = len(c[key].usage)
		}
		if maxLenKey < len(key) {
			maxLenKey = len(key)
		}
	}
	for key := range t {
		if maxLenUsage < len(t[key].usage) {
			maxLenUsage = len(t[key].usage)
		}
		if maxLenKey < len(key) {
			maxLenKey = len(key)
		}
	}
	return maxLenUsage + 1, maxLenKey + 1
}
func clearTerminal() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		log.Println("operating system not supported for clearing terminal")
	}
}
