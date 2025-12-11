package ebpf

import (
	"fmt"
	"os"
	"os/exec"
	"io"
	"path/filepath"
	"xvtlas/config"
	"xvtlas/logger"
	"xvtlas/report"
	"xvtlas/utils"
	"strings"
)
func RunPipeline(rootPath, prettyPath, kernelVersion, exportPath string, interactive bool) {
	var rows []report.CSVRow

	// Global compilation
	compilationLog, err := utils.RunMakeAll(rootPath)
	fmt.Println("Stuff is compiled :", compilationLog)
	if err != nil {
		logger.LogError("Makefile", compilationLog)
		if interactive && !utils.ConfirmPrompt("Compilation failed. Continue?") {
			os.Exit(0)
		}
	}

	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		ebpfFile := filepath.Join(path, "main.c")
		yamlFile := filepath.Join(path, "config.yaml")

		if !utils.FileExists(ebpfFile) || !utils.FileExists(yamlFile) {
			return nil
		}

		cfg, err := config.ParseConfig(yamlFile)
		if err != nil {
			logger.LogError(ebpfFile, err.Error())
			if interactive && !utils.ConfirmPrompt("Continue to next program?") {
				os.Exit(0)
			}
			return nil
		}

		row := report.CSVRow{
			Filename:        ebpfFile,
			LoadParameters:  fmt.Sprintf("%v", cfg.EBPFProgram),
			KernelVersion:   kernelVersion,
			Compiled:        true, // assume success if object exists ??
		}

		oFile := utils.FindObjectFile(path)
		if oFile == "" {
			row.Compiled = false
			row.LoadOutput += "No .o file found after make.\n"
		} else {
			utils.RunVerifier(oFile, ebpfFile, prettyPath, cfg)
			//utils.LoadEBPF(oFile, cfg, &row)
		}

		logger.SaveLog(ebpfFile, row.LoadOutput)
		rows = append(rows, row)
		return nil
	})

	if err != nil {
		fmt.Println("Error scanning root path:", err)
	}

	report.ExportCSV(rows, exportPath)
}


func RunPipelinePatch(patchRoot, baseFile, prettyPath, kernelVersion, exportPath string, interactive, removePatched bool) {
	var rows []report.CSVRow

	submoduleRoot := filepath.Dir(baseFile)
	absSubmoduleRoot,err1 := filepath.Abs(submoduleRoot)
	if err1 != nil {
		logger.LogError(absSubmoduleRoot, fmt.Sprintf("Failed to get absolute path: %s", err1.Error()))
		return
	}


	fmt.Println("Running pipeline")

	err := filepath.Walk(patchRoot, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		patchFile := filepath.Join(path, "patch.diff")
		yamlFile := filepath.Join(path, "config.yaml")
		if !utils.FileExists(patchFile) || !utils.FileExists(yamlFile) {
			return nil
		}


		fmt.Println(absSubmoduleRoot)
		// Ensure submodule base file is reset to clean state
		resetCmd := exec.Command("git", "-C", absSubmoduleRoot, "reset", "--hard")
		if out, err := resetCmd.CombinedOutput(); err != nil {
			logger.LogError(baseFile, fmt.Sprintf("Failed to reset repo state: %s\nOutput: %s", err.Error(), out))
			return nil
		}

		// Convert patch file to absolute path to avoid path resolution issues
		absPatchFile, err := filepath.Abs(patchFile)
		if err != nil {
			logger.LogError(patchFile, fmt.Sprintf("Failed to get absolute path: %s", err.Error()))
			return nil
		}
		
		
		fmt.Println("pathc file ", absPatchFile)
		fmt.Println("Root", absSubmoduleRoot)
		// Apply patch using git apply with absolute path from submodule root
		applyCmd := exec.Command("git", "-C", absSubmoduleRoot, "apply", absPatchFile)
		if output, err := applyCmd.CombinedOutput(); err != nil {
			logger.LogError(patchFile, fmt.Sprintf("Failed to apply patch: %s\nOutput: %s", err.Error(), output))
			if interactive && !utils.ConfirmPrompt("Continue after failed patch?") {
				os.Exit(1)
			}
			return nil
		}

		// Copy patched source to test folder as main.c
		patchedDest := filepath.Join(path, "main.c")
		srcFile, err := os.Open(baseFile)
		if err != nil {
			logger.LogError(baseFile, "Failed to open patched base file")
			return nil
		}
		defer srcFile.Close()

		dstFile, err := os.Create(patchedDest)
		if err != nil {
			logger.LogError(patchedDest, "Failed to create target patched file")
			return nil
		}
		defer dstFile.Close()

		if _, err = io.Copy(dstFile, srcFile); err != nil {
			logger.LogError(patchedDest, "Failed to copy patched source")
			return nil
		}

		// Compile with make in target dir
		compilationLog, err := utils.RunMake(path)
		if err != nil {
			logger.LogError("Makefile", compilationLog)
			if interactive && !utils.ConfirmPrompt("Compilation failed. Continue?") {
				os.Exit(1)
			}
		}

		cfg, err := config.ParseConfig(yamlFile)
		if err != nil {
			logger.LogError(yamlFile, err.Error())
			return nil
		}

		row := report.CSVRow{
			Filename:        patchedDest,
			LoadParameters:  fmt.Sprintf("%v", cfg.EBPFProgram),
			KernelVersion:   kernelVersion,
			Compiled:        true,
		}

		oFile := utils.FindObjectFile(path)
		if oFile == "" {
			row.Compiled = false
			row.LoadOutput += "No .o file found after make.\n"
		} else {
			utils.RunVerifier(oFile, patchedDest, prettyPath, cfg)
			row.Loaded = row.Verified
		}

		logger.SaveLog(patchedDest, row.LoadOutput)
		rows = append(rows, row)

		// Optionally remove main.c copy
		if removePatched {
			_ = os.Remove(patchedDest)
		}

		// Unpin loaded program if necessary
		progName := cfg.EBPFProgram.Name
		if progName == "" && oFile != "" {
			progName = filepath.Base(oFile)
		}
		_ = exec.Command("sudo", "rm", "-f", filepath.Join("/sys/fs/bpf", progName)).Run()

		return nil
	})

	if err != nil {
		fmt.Println("Error scanning patch path:", err)
	}

	report.ExportCSV(rows, exportPath)
}

func RunPipelineNew(patchRoot, baseFile, prettyPath, kernelVersion, exportPath string, interactive, saveLogsToFile bool, keepPatched bool) {
	var rows []report.CSVRow

	submoduleRoot := filepath.Dir(baseFile)
	absSubmoduleRoot, err := filepath.Abs(submoduleRoot)
	if err != nil {
		logger.LogError(submoduleRoot, fmt.Sprintf("Failed to get absolute path: %s", err.Error()))
		return
	}

	fmt.Println("Running pipeline")

	err = filepath.Walk(patchRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.LogError(path, fmt.Sprintf("Walk error: %s", err.Error()))
			return nil
		}

		if !info.IsDir() {
			return nil
		}

		patchFiles, err := filepath.Glob(filepath.Join(path, "*.patch"))
		if err != nil {
			logger.LogError(path, fmt.Sprintf("Failed to glob patch files: %s", err.Error()))
			return nil
		}

		for _, patchFile := range patchFiles {
			fmt.Printf("Applying patch: %s\n", patchFile)

			// Save original HEAD
			origHeadCmd := exec.Command("git", "-C", absSubmoduleRoot, "rev-parse", "HEAD")
			origHead, err := origHeadCmd.Output()
			if err != nil {
				logger.LogError("git rev-parse", "Failed to capture HEAD")
				continue
			}
			origHeadStr := strings.TrimSpace(string(origHead))

			absPatchFile, err := filepath.Abs(patchFile)
			if err != nil {
				logger.LogError(patchFile, fmt.Sprintf("Failed to get absolute path: %s", err.Error()))
				continue
			}

			// Determine output subdirectory for logs
			var outputDir string
			if saveLogsToFile {
				patchDir := filepath.Base(path)
				outputDir = filepath.Join(exportPath, patchDir)
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					logger.LogError(outputDir, fmt.Sprintf("Failed to create output directory: %s", err.Error()))
					continue
				}
			}

			// Apply patch
			applyCmd := exec.Command("git", "-C", absSubmoduleRoot, "am", absPatchFile)
			if output, err := applyCmd.CombinedOutput(); err != nil {
				logger.LogError(patchFile, fmt.Sprintf("Failed to apply patch: %s\nOutput: %s", err.Error(), string(output)))
				_ = exec.Command("git", "-C", absSubmoduleRoot, "am", "--abort").Run()
				continue
			}

			// Compile
			compilationLog, err := utils.RunMake(absSubmoduleRoot)

			row := report.CSVRow{
				Filename:       patchFile,
				LoadParameters: "none",
				KernelVersion:  kernelVersion,
				Compiled:       true,
			}

			if saveLogsToFile {
				_ = os.WriteFile(filepath.Join(outputDir, "make.log"), []byte(compilationLog), 0644)
			}

			if err != nil || strings.Contains(compilationLog, "error:") {
				row.Compiled = false
				//row.LoadOutput += compilationLog
				//logger.LogError("Makefile", compilationLog)
				//logger.SaveLog(patchFile, compilationLog)
				rows = append(rows, row)
				resetGit(absSubmoduleRoot, origHeadStr)
				continue
			}

			// Check for object file
			oFile := utils.FindObjectFile(absSubmoduleRoot)
			if oFile == "" {
				row.Compiled = false
				row.LoadOutput += "No .o file found after make.\n"
				logger.SaveLog(patchFile, row.LoadOutput)
				rows = append(rows, row)
				resetGit(absSubmoduleRoot, origHeadStr)
				continue
			}

			// Run verifier
			loadOutput := utils.RunVerifier(oFile, baseFile, prettyPath, nil)
			//row.LoadOutput += string(loadOutput)
			row.Verified = !strings.Contains(string(loadOutput), "BPF program load failed")
			row.Loaded = row.Verified

			if saveLogsToFile {
				_ = os.WriteFile(filepath.Join(outputDir, "verifier.log"), loadOutput, 0644)
			}

			if keepPatched {
				_ = exec.Command("cp", baseFile, outputDir).Run()
			}

			_ = exec.Command("sudo", "rm", "-f", filepath.Join("/sys/fs/bpf/", filepath.Base(oFile))).Run()
			//logger.SaveLog(patchFile, string(loadOutput))
			rows = append(rows, row)

			// Perform make clean
			cleanCmd := exec.Command("make", "-C", absSubmoduleRoot, "clean")
			if output, err := cleanCmd.CombinedOutput(); err != nil {
				logger.LogError("make clean", fmt.Sprintf("Failed to run make clean: %s\nOutput: %s", err.Error(), string(output)))
				// do not fail
			}

			resetGit(absSubmoduleRoot, origHeadStr)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error scanning patch path:", err)
	}

	report.ExportCSV(rows, exportPath)

	fmt.Println("\nFinal Report:")
	for _, r := range rows {
		status := "[ ]"
		if r.Loaded {
			status = "[✔]"
		}
		fmt.Printf("%s %s | Compiled: %v | Loaded: %v\n", status, r.Filename, r.Compiled, r.Loaded)
	}
}

func resetGit(repoPath, commit string) {
	_ = exec.Command("git", "-C", repoPath, "reset", "--hard", commit).Run()
}

func RunSingle(singlePatchPath string, baseFile string) {
	swapFilePath := "/tmp/xvtlas.swp"

	// Check if swap file already exists
	if _, err := os.Stat(swapFilePath); err == nil {
		fmt.Println("File already exists first run --destroy")
		os.Exit(1)
	}

	info, err := os.Stat(singlePatchPath)
	if info.IsDir() || err != nil {
		fmt.Printf("The path to patch is not a file (got: %s)", singlePatchPath)
		os.Exit(1)
	}

	absSubmoduleRoot, err := filepath.Abs(filepath.Dir(baseFile))
	if err != nil {
		logger.LogError(baseFile, fmt.Sprintf("Failed to get absolute base path: %s", err))
		return
	}

	// Save HEAD and baseFile path to /tmp/xvtlas.swp
	origHeadCmd := exec.Command("git", "-C", absSubmoduleRoot, "rev-parse", "HEAD")
	origHead, err := origHeadCmd.Output()
	if err != nil {
		logger.LogError("git rev-parse", "Failed to get HEAD")
		os.Exit(1)
	}
	head := strings.TrimSpace(string(origHead))
	stateFileContent := fmt.Sprintf("%s\n%s\n", head, absSubmoduleRoot)
	err = os.WriteFile("/tmp/xvtlas.swp", []byte(stateFileContent), 0644)
	if err != nil {
		logger.LogError("state-file", "Failed to write /tmp/xvtlas.swp")
		os.Exit(1)
	}

	// Apply patch
	absPatchFile, err := filepath.Abs(singlePatchPath)
	if err != nil {
		logger.LogError(singlePatchPath, fmt.Sprintf("Failed to resolve absolute patch path: %s", err))
		os.Exit(1)
	}

	fmt.Println("Applying patch:", absPatchFile)
	applyCmd := exec.Command("git", "-C", absSubmoduleRoot, "am", absPatchFile)
	if output, err := applyCmd.CombinedOutput(); err != nil {
		logger.LogError("git am", fmt.Sprintf("Failed to apply patch:\n%s", string(output)))
		_ = exec.Command("git", "-C", absSubmoduleRoot, "am", "--abort").Run()
		os.Exit(1)
	}

	// Compile
	compilationLog, err := utils.RunMake(absSubmoduleRoot)
	fmt.Println(compilationLog)

	if err != nil || strings.Contains(compilationLog, "error:") {
		logger.LogError("compile", "Compilation failed")
		os.Exit(1)
	}

	// Launch start_session.sh 
	cmd := exec.Command("bash", "./start_session.sh")
	cmd.Dir = absSubmoduleRoot
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting interactive session (exit tmux to continue)...")
	if err := cmd.Run(); err != nil {
		logger.LogError("start_session", fmt.Sprintf("Failed to start session: %s", err))
		os.Exit(1)
	}
}
