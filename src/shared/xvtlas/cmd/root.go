package cmd

import (
	"fmt"
	"os"
	"xvtlas/ebpf"
	"xvtlas/logger"
	"xvtlas/utils"

	"github.com/spf13/cobra"
)

var (
	rootPath        string
	interactive     bool
	verbose         bool
	exportPath      string
	prettyFilePath  string
	kernelVersion   string
	patchPath       string
	baseFile        string
	saveLogs        bool
	runSingle       string
	destroy         bool
	keepPatched	bool
)

var rootCmd = &cobra.Command{
	Use:   "xvtlas",
	Short: "XVTlas - XDP Verifier Launch Automation Suite",
	Long:  "XVTlas is a tool to compile, patch, load, and verify eBPF programs using bpftool and a pretty verifier.",
	SilenceUsage:  true,
	SilenceErrors: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(os.Args) == 1 {
			cmd.Help()
			return fmt.Errorf("no arguments provided")
		}
		if runSingle != "" && baseFile == "" {
			return fmt.Errorf("--base-file must be set when using --run-single")
		}
		if runSingle != "" && patchPath != "" {
			return fmt.Errorf("--patch-path cannot be used with --run-single")
		}
		if patchPath != "" && baseFile == "" {
			return fmt.Errorf("--base-file must be specified when using --patch-path")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(verbose, exportPath)

		switch {
		// --run-single mode
		case runSingle != "":
			if baseFile == "" {
				fmt.Println("Error: --base-file must be set when using --run-single")
				os.Exit(1)
			}

			// Disallow any other options
			if rootPath != "." || patchPath != "" || prettyFilePath != "" || kernelVersion != "" ||
				interactive || saveLogs || exportPath != "./output" || verbose {
				fmt.Println("Error: When using --run-single, only --base-file is allowed as an additional flag.")
				os.Exit(1)
			}

			ebpf.RunSingle(runSingle, baseFile)

		// --destroy mode
		case destroy:
			// Disallow all other flags
			if runSingle != "" || baseFile != "" || patchPath != "" || rootPath != "." || prettyFilePath != "" ||
				kernelVersion != "" || interactive || saveLogs || exportPath != "./output" || verbose {
				fmt.Println("Error: --destroy must be used alone with no other flags.")
				os.Exit(1)
			}
			utils.DestroyPreviousState()

		// Default: normal multi-patch or rootPath mode
		default:
			if patchPath != "" {
				if rootPath != "." {
					fmt.Println("Error: --patch-path and --path are mutually exclusive")
					os.Exit(1)
				}
				if baseFile == "" {
					fmt.Println("Error: --base-file is required with --patch-path")
					os.Exit(1)
				}
				ebpf.RunPipelineNew(patchPath, baseFile, prettyFilePath, kernelVersion, exportPath, interactive, saveLogs, keepPatched)
			} else {
				ebpf.RunPipeline(rootPath, prettyFilePath, kernelVersion, exportPath, interactive)
			}
		}
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&rootPath, "path", "p", ".", "Root path to eBPF program folders")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Stop on errors and prompt user to continue")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&exportPath, "export", "e", "./output", "Export location for logs and CSV")
	rootCmd.Flags().StringVar(&prettyFilePath, "pretty", "", "Pretty verifier file path")
	rootCmd.Flags().StringVar(&kernelVersion, "kernel", "", "Target kernel version for accounting")
	rootCmd.Flags().StringVar(&patchPath, "patch-path", "", "Path to folders containing patches and configs")
	rootCmd.Flags().StringVar(&baseFile, "base-file", "", "Absolute path to the master file to apply each patch to")
	rootCmd.Flags().BoolVar(&saveLogs, "save-logs", false, "Save logs for each patch (default: false)")
	rootCmd.Flags().StringVar(&runSingle, "run-single", "", "Run a single patch file against a base file")
	rootCmd.Flags().BoolVar(&destroy, "destroy", false, "Restore Git HEAD and clean build based on last state")
	rootCmd.Flags().BoolVar(&keepPatched, "keep-patched", false, "Keeps the patch applied base file for every patch in the output folder")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

