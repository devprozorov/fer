package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

// boxWidth is total box width including borders.
const boxWidth = 54

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// drawBox prints a titled ASCII box with the given menu lines.
func drawBox(title string, lines []string) {
	border := "+" + strings.Repeat("-", boxWidth-2) + "+"
	inner := boxWidth - 4
	fmt.Println(border)

	// Centered title row.
	titlePad := (inner - len(title)) / 2
	if titlePad < 0 {
		titlePad = 0
	}
	trailing := inner - titlePad - len(title)
	if trailing < 0 {
		trailing = 0
	}
	fmt.Printf("| %s%s%s |\n",
		strings.Repeat(" ", titlePad), title, strings.Repeat(" ", trailing))
	fmt.Println(border)

	for _, l := range lines {
		pad := inner - len(l)
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("| %s%s |\n", l, strings.Repeat(" ", pad))
	}
	fmt.Println(border)
}

// prompt prints label and returns the trimmed user input.
func prompt(reader *bufio.Reader, label string) string {
	fmt.Printf("%s: ", label)
	v, _ := reader.ReadString('\n')
	return strings.TrimSpace(v)
}

// pause waits for the user to press Enter.
func pause(reader *bufio.Reader) {
	fmt.Print("\n[Press Enter to continue]")
	_, _ = reader.ReadString('\n')
}

// runCapture runs a command and returns its combined output.
func runCapture(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// restartService restarts a systemd service and prints the result.
func restartService(name string) {
	out, err := runCapture("systemctl", "restart", name)
	if err != nil {
		fmt.Printf("Failed to restart %s: %s\n", name, out)
		return
	}
	fmt.Printf("Service %s restarted.\n", name)
}
