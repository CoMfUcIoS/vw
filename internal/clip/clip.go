package clip

import (
	"errors"
	"os/exec"
	"runtime"
)

func Copy(text string) error {
	cmd, err := command()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	_, _ = stdin.Write([]byte(text))
	_ = stdin.Close()
	return cmd.Wait()
}

func command() (*exec.Cmd, error) {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbcopy"), nil
	case "windows":
		return exec.Command("clip"), nil
	default:
		if _, err := exec.LookPath("wl-copy"); err == nil {
			return exec.Command("wl-copy"), nil
		}
		if _, err := exec.LookPath("xclip"); err == nil {
			return exec.Command("xclip", "-selection", "clipboard"), nil
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return exec.Command("xsel", "--clipboard", "--input"), nil
		}
		return nil, errors.New("no clipboard command found; install wl-copy, xclip, or xsel")
	}
}
