package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type TerminalOption struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func GetAvailableTerminals() []TerminalOption {
	platform := runtime.GOOS
	var terminals []TerminalOption

	switch platform {
	case "linux":
		terminals = []TerminalOption{
			{"Kitty", "kitty"},
			{"GNOME Terminal", "gnome-terminal"},
			{"Konsole", "konsole"},
			{"XTerm", "xterm"},
			{"Alacritty", "alacritty"},
			{"Terminator", "terminator"},
		}
	case "darwin":
		terminals = []TerminalOption{
			{"iTerm2", "open -a iTerm"},
			{"Terminal", "open -a Terminal"},
			{"Kitty", "kitty"},
			{"Alacritty", "alacritty"},
		}
	default:
		return []TerminalOption{}
	}

	var available []TerminalOption
	for _, term := range terminals {
		if isCommandAvailable(strings.Split(term.Command, " ")[0]) {
			available = append(available, term)
		}
	}
	return available
}

func OpenTerminal(terminalCmd, path string) error {
	platform := runtime.GOOS
	var cmd *exec.Cmd

	switch platform {
	case "linux":
		switch {
		case strings.Contains(terminalCmd, "kitty"):
			cmd = exec.Command("kitty", "--directory", path, "nvim")
		case strings.Contains(terminalCmd, "gnome-terminal"):
			cmd = exec.Command("gnome-terminal", "--working-directory", path, "--", "nvim")
		case strings.Contains(terminalCmd, "konsole"):
			cmd = exec.Command("konsole", "--workdir", path, "-e", "nvim")
		case strings.Contains(terminalCmd, "alacritty"):
			cmd = exec.Command("alacritty", "--working-directory", path, "-e", "nvim")
		case strings.Contains(terminalCmd, "terminator"):
			cmd = exec.Command("terminator", "--working-directory", path, "-x", "nvim")
		default:
			cmd = exec.Command("xterm", "-e", fmt.Sprintf("cd %s && nvim", path))
		}
	case "darwin":
		script := fmt.Sprintf(`tell application "Terminal" to do script "cd %s && nvim"`, path)
		cmd = exec.Command("osascript", "-e", script)
	default:
		return fmt.Errorf("platform tidak didukung: %s", platform)
	}

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}

func GetOS() string {
	return runtime.GOOS
}