package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/manifoldco/promptui"
)

const (
	ShellZsh  = "zsh"
	ShellBash = "bash"
	ShellPwsh = "windows"
)

var shellChoices = []string{
	ShellZsh + " (Oh My Zsh / Zsh)",
	ShellBash + " (Bash)",
	ShellPwsh + " (PowerShell)",
}

func InstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		return filepath.Join(home, ".local", "bin"), nil
	}
	localBin := filepath.Join(home, ".local", "bin")
	if err := os.MkdirAll(localBin, 0755); err == nil {
		return localBin, nil
	}
	goBin := filepath.Join(home, "go", "bin")
	if err := os.MkdirAll(goBin, 0755); err == nil {
		return goBin, nil
	}
	return localBin, nil
}

func FindRepoRoot() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

const binaryName = "aws-profile"

func removeExisting(installDir string) {
	bin := filepath.Join(installDir, binaryName)
	_ = os.Remove(bin)
	if runtime.GOOS == "windows" {
		_ = os.Remove(bin + ".exe")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configDir := filepath.Join(home, ".config", "awsp")
	_ = os.Remove(filepath.Join(configDir, "completion.zsh"))
	_ = os.Remove(filepath.Join(configDir, "completion.bash"))
}

func BuildBinary(repoRoot string) (string, error) {
	out := filepath.Join(repoRoot, binaryName)
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", out, ".")
	cmd.Dir = repoRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("go build: %w", err)
	}
	return out, nil
}

func CopyBinary(src, destDir string) (string, error) {
	name := binaryName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	dest := filepath.Join(destDir, name)
	data, err := os.ReadFile(src)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(dest, data, 0755); err != nil {
		return "", err
	}
	return dest, nil
}

func SelectShell() (string, error) {
	sel := promptui.Select{
		Label: "Select your shell",
		Items: shellChoices,
		Size:  4,
	}
	idx, chosen, err := sel.Run()
	if err != nil {
		return "", err
	}
	_ = idx
	if strings.HasPrefix(chosen, ShellZsh) {
		return ShellZsh, nil
	}
	if strings.HasPrefix(chosen, ShellBash) {
		return ShellBash, nil
	}
	return ShellPwsh, nil
}

func ConfigPath(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch shell {
	case ShellZsh:
		return filepath.Join(home, ".zshrc"), nil
	case ShellBash:
		rc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(rc); os.IsNotExist(err) {
			rc = filepath.Join(home, ".bash_profile")
		}
		return rc, nil
	case ShellPwsh:
		doc, _ := os.UserConfigDir()
		return filepath.Join(doc, "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
	default:
		return filepath.Join(home, ".zshrc"), nil
	}
}

func Snippet(shell, installDir string) string {
	pathLine := ""
	if runtime.GOOS != "windows" {
		pathLine = fmt.Sprintf("export PATH=%q:$PATH\n", installDir)
	} else {
		pathLine = fmt.Sprintf("$env:Path = %q + \";\" + $env:Path\n", installDir)
	}
	compLine := completionSourceLine(shell)
	switch shell {
	case ShellZsh, ShellBash:
		return "\n# awsp - AWS Profile Switcher\n" +
			pathLine +
			"awsp() { if [ $# -gt 0 ]; then command aws-profile \"$@\"; else local f=$(mktemp); AWSP_EXPORT_FILE=$f command aws-profile; eval \"$(cat \"$f\")\"; rm -f \"$f\"; fi; }\n" +
			compLine
	case ShellPwsh:
		return "\n# awsp - AWS Profile Switcher\n" +
			pathLine +
			"function awsp { if ($args.Count -gt 0) { aws-profile @args } else { aws-profile | ForEach-Object { $l = $_.Trim(); if ($l -match '^export ([^=]+)=(.+)$') { Set-Item -Path \"env:$($matches[1])\" -Value $matches[2].Trim('\"') } } } }\n"
	default:
		return "\n# awsp\n" + pathLine +
			"awsp() { if [ $# -gt 0 ]; then command aws-profile \"$@\"; else local f=$(mktemp); AWSP_EXPORT_FILE=$f command aws-profile; eval \"$(cat \"$f\")\"; rm -f \"$f\"; fi; }\n" +
			compLine
	}
}

func writeCompletionScript(awspPath, shell string) {
	ext := "zsh"
	if shell == ShellBash {
		ext = "bash"
	}
	if shell != ShellZsh && shell != ShellBash {
		return
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := filepath.Join(home, ".config", "awsp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	outPath := filepath.Join(dir, "completion."+ext)
	cmd := exec.Command(awspPath, "completion", shell)
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return
	}
	_ = os.WriteFile(outPath, out, 0644)
	fmt.Printf("Completion script: %s\n", outPath)
}

func completionSourceLine(shell string) string {
	if shell == ShellZsh {
		return "[ -f ~/.config/awsp/completion.zsh ] && source ~/.config/awsp/completion.zsh\n"
	}
	if shell == ShellBash {
		return "[ -f ~/.config/awsp/completion.bash ] && source ~/.config/awsp/completion.bash\n"
	}
	return ""
}

func AppendConfig(path, snippet string) (appended bool, err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil && !os.IsNotExist(err) {
		return false, err
	}
	content, _ := os.ReadFile(path)
	if strings.Contains(string(content), "# awsp") {
		return false, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.WriteString(snippet)
	return err == nil, err
}

func appendCompletionIfMissing(path, shell string) {
	line := completionSourceLine(shell)
	if line == "" {
		return
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	s := string(content)
	if strings.Contains(s, "completion.zsh") || strings.Contains(s, "completion.bash") {
		return
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(line)
}

func Run() error {
	installDir, err := InstallDir()
	if err != nil {
		return err
	}

	fmt.Println("Removing existing aws-profile binary and completion scripts...")
	removeExisting(installDir)

	var binaryPath string
	repoRoot, inRepo := FindRepoRoot()
	if inRepo {
		fmt.Println("Building binary (go build)...")
		binaryPath, err = BuildBinary(repoRoot)
		if err != nil {
			return err
		}
		fmt.Println("Build done. Copying to install dir...")
	} else {
		fmt.Println("Not in project directory: copying current binary (no rebuild).")
		fmt.Println("To install a new build, run 'go run . install' from the project.")
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot get executable path: %w", err)
		}
		binaryPath = execPath
	}

	fmt.Printf("Installing to %s ...\n", installDir)
	dest, err := CopyBinary(binaryPath, installDir)
	if err != nil {
		return fmt.Errorf("copy binary: %w", err)
	}
	fmt.Printf("Installed: %s\n", dest)

	shell, err := SelectShell()
	if err != nil {
		fmt.Println("Canceled")
		return nil
	}

	configPath, err := ConfigPath(shell)
	if err != nil {
		return err
	}

	snippet := Snippet(shell, installDir)
	appended, err := AppendConfig(configPath, snippet)
	if err != nil {
		fmt.Printf("Could not write to %s: %v\n", configPath, err)
		fmt.Println("Add this to your shell config:")
		fmt.Print(snippet)
		return nil
	}
	if !appended {
		appendCompletionIfMissing(configPath, shell)
	}
	writeCompletionScript(dest, shell)
	fmt.Printf("Config written to %s\n", configPath)
	fmt.Println("Restart your terminal or run: source", configPath)
	return nil
}
