package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"xvtlas/report"
	"xvtlas/config"
	"path/filepath"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ConfirmPrompt(msg string) bool {
	fmt.Print(msg + " (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	return response == "y\n"
}
//Still need checks for some edge cases
func RunMakeAll(rootPath string) (string, error) {
	cmd := exec.Command("make", "-C", rootPath)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

//the .o file should be found in each folder
func FindObjectFile(folderPath string) string { 
	files, _ := os.ReadDir(folderPath)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".o" {
			return filepath.Join(folderPath, file.Name())
		}
	}
	return ""
}

func RunVerifier(oFile, cFile, prettyPath string, cfg *config.EBPFConfig) []byte {
	var progName string

	if cfg == nil || cfg.EBPFProgram.Name == "" {
		progName = filepath.Base(oFile)
	} else {
		progName = cfg.EBPFProgram.Name
	}

	pinPath := progName

	cmdStr := fmt.Sprintf("sudo bpftool prog load %s /sys/fs/bpf/%s 2>&1 | python3 %s -c %s -o %s", oFile, pinPath, prettyPath, cFile, oFile)
	fmt.Println("Running Verifier Command:", cmdStr)

	cmd := exec.Command("bash", "-c", cmdStr)
	output, _ := cmd.CombinedOutput()

	return output
}

func LoadEBPF(oFile string, cfg *config.EBPFConfig, row *report.CSVRow) {
	//1) take as input the compile_object
	//2) use target from config_to try attach if configured ()
	//3) NO need for this now

	row.Loaded = true
}



//Function to run make on the single folder in order to compile only the actively patched file
func RunMake(rootPath string) (string, error) {
	cmd := exec.Command("make", "-C", rootPath)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// DestroyPreviousState restores Git HEAD and runs `make clean`
// using the state saved in /tmp/xvtlas.swp. It exits the program on failure.
func DestroyPreviousState() {
	const stateFile = "/tmp/xvtlas.swp"

	data, err := os.ReadFile(stateFile)
	if err != nil {
		fmt.Println("Nothing to destroy or failed to read", stateFile)
		os.Exit(1)
	}

	parts := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(parts) != 2 {
		fmt.Println("Invalid state file format")
		os.Exit(1)
	}

	head := parts[0]
	baseDir := parts[1]

	fmt.Println("Restoring Git HEAD to:", head)
	resetCmd := exec.Command("git", "-C", baseDir, "reset", "--hard", head)
	resetCmd.Stdout = os.Stdout
	resetCmd.Stderr = os.Stderr
	if err := resetCmd.Run(); err != nil {
		fmt.Println("Failed to reset HEAD:", err)
		os.Exit(1)
	}

	fmt.Println("Running destroy_session.sh in:", baseDir)
	destroyScript := exec.Command("bash", "kill_session.sh")
	destroyScript.Dir = baseDir
	destroyScript.Stdout = os.Stdout
	destroyScript.Stderr = os.Stderr
	if err := destroyScript.Run(); err != nil {
		fmt.Println("Warning: destroy_session.sh failed or missing:", err)
		// continue cleanup anyway
	}

	fmt.Println("Cleaning build directory:", baseDir)
	cleanCmd := exec.Command("make", "-C", baseDir, "clean")
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	_ = cleanCmd.Run()

	_ = os.Remove(stateFile)
	fmt.Println("✅ Restore and cleanup complete.")
}
