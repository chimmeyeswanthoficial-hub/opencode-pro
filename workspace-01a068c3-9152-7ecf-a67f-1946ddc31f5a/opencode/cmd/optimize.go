package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/spf13/cobra"
)

var optimizeJSON bool

var optimizeCmd = &cobra.Command{
	Use:   "optimize [prompt]",
	Short: "Optimize and enhance a prompt using project introspection and context triggers",
	Long: `Analyze a raw or rough prompt, parse @, /, and # mentions, introspect the workspace,
and generate a high-precision, structured Golden Prompt tailored for AI coding agents.`,
	Example: `
  # Optimize a raw prompt
  opencode optimize "add jwt auth middleware @auth.go /test #rules"

  # Output as JSON
  opencode optimize "fix db connection leak #rules:strict-types" --json
  `,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rawPrompt := strings.Join(args, " ")
		opt := optimizer.NewPromptOptimizer(".")
		result := opt.Optimize(context.Background(), rawPrompt)

		if optimizeJSON {
			data, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		fmt.Println("================================================================================")
		fmt.Println("⚡ OPENCODE SMART PROMPT OPTIMIZER")
		fmt.Println("================================================================================")
		fmt.Printf("🔍 Detected Intent: %s\n", result.DetectedIntent)
		if len(result.TargetFiles) > 0 {
			fmt.Printf("📁 Target Files:    %s\n", strings.Join(result.TargetFiles, ", "))
		}
		if len(result.AppliedSkills) > 0 {
			fmt.Printf("🛠️  Skills Applied:  %s\n", strings.Join(result.AppliedSkills, ", "))
		}
		if len(result.AppliedRules) > 0 {
			fmt.Printf("📜 Rules Attached:  %s\n", strings.Join(result.AppliedRules, ", "))
		}
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("🧠 Reasoning Chain:")
		fmt.Println(result.Reasoning)
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("✨ Synthesized Golden Prompt:")
		fmt.Println(result.OptimizedPrompt)
		fmt.Println("================================================================================")
		return nil
	},
}

func init() {
	optimizeCmd.Flags().BoolVar(&optimizeJSON, "json", false, "Output optimization result as JSON")
	rootCmd.AddCommand(optimizeCmd)
}
