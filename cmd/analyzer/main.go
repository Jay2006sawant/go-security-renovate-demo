package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/Jay2006sawant/go-security-renovate-demo/internal/analyzer"
)

var (
	version = "1.0.0"
	
	// Color functions for enhanced output
	red    = color.New(color.FgRed, color.Bold).SprintFunc()
	green  = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow = color.New(color.FgYellow, color.Bold).SprintFunc()
	blue   = color.New(color.FgBlue, color.Bold).SprintFunc()
)

func main() {
	rootCmd := &cobra.Command{ 
		Use:   "analyzer",
		Short: "Git Repository Security Analyzer",
		Long: `A demonstration tool that analyzes Git repositories for security insights.
This tool uses go-git library (vulnerable version) to show how Renovate
can help maintain secure dependencies in Go projects.`,
		Version: version,
	}

	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze a Git repository",
		Long:  "Analyze a Git repository and generate a security report",
		RunE:  runAnalyze,
	}

	analyzeCmd.Flags().StringP("repo", "r", "", "Repository URL to analyze (required)")
	analyzeCmd.Flags().StringP("output", "o", "console", "Output format: console, json")
	analyzeCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	analyzeCmd.MarkFlagRequired("repo")

	demoCmd := &cobra.Command{
		Use:   "demo",
		Short: "Run demo analysis with sample repositories",
		Long:  "Analyze sample repositories to demonstrate the tool functionality",
		RunE:  runDemo,
	}

	vulnerabilityCmd := &cobra.Command{
		Use:   "vulnerability",
		Short: "Show information about the vulnerable dependency",
		Long:  "Display details about CVE-2023-49568 affecting go-git library",
		RunE:  showVulnerability,
	}

	rootCmd.AddCommand(analyzeCmd, demoCmd, vulnerabilityCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %v\n", red("✗"), err)
		os.Exit(1)
	}
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	repoURL, _ := cmd.Flags().GetString("repo")
	outputFormat, _ := cmd.Flags().GetString("output")
	verbose, _ := cmd.Flags().GetBool("verbose")

	fmt.Printf("%s Starting analysis of repository: %s\n", blue("ℹ"), repoURL)
	
	if verbose {
		fmt.Printf("%s Using vulnerable go-git version for demonstration\n", yellow("⚠"))
	}

	gitAnalyzer := analyzer.NewGitAnalyzer()
	report, err := gitAnalyzer.AnalyzeRepository(repoURL)
	if err != nil {
		return fmt.Errorf("failed to analyze repository: %w", err)
	}

	if outputFormat == "json" {
		return report.OutputJSON()
	}
	
	return report.OutputConsole()
}

func runDemo(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s Running demo analysis with sample repositories\n", green("✓"))
	
	sampleRepos := []string{
		"https://github.com/go-git/go-git",
		"https://github.com/spf13/cobra",
		"https://github.com/fatih/color",
	}

	gitAnalyzer := analyzer.NewGitAnalyzer()
	
	for i, repo := range sampleRepos {
		fmt.Printf("\n%s [%d/%d] Analyzing: %s\n", blue("→"), i+1, len(sampleRepos), repo)
		
		report, err := gitAnalyzer.AnalyzeRepository(repo)
		if err != nil {
			fmt.Printf("%s Failed to analyze %s: %v\n", red("✗"), repo, err)
			continue
		}
		
		report.OutputConsole()
	}
	
	return nil
}

func showVulnerability(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s CVE-2023-49568 - Path Traversal Vulnerability in go-git\n", red("🔒"))
	fmt.Println()
	fmt.Printf("%s Severity: HIGH\n", red("•"))
	fmt.Printf("%s CVSS Score: 7.5\n", red("•"))
	fmt.Printf("%s Affected Versions: < 5.11.0\n", red("•"))
	fmt.Printf("%s Current Version: 5.4.2 (VULNERABLE)\n", red("•"))
	fmt.Println()
	fmt.Printf("%s Description:\n", blue("📋"))
	fmt.Println("  The go-git library is vulnerable to path traversal attacks when")
	fmt.Println("  processing Git repositories. An attacker could potentially access")
	fmt.Println("  files outside the intended directory structure during Git operations.")
	fmt.Println()
	fmt.Printf("%s Impact:\n", yellow("⚠"))
	fmt.Println("  • Unauthorized file system access")
	fmt.Println("  • Potential data exfiltration")
	fmt.Println("  • Directory traversal attacks")
	fmt.Println()
	fmt.Printf("%s Mitigation:\n", green("🛡"))
	fmt.Println("  Update go-git to version 5.11.0 or later")
	fmt.Println("  This vulnerability demonstrates why automated dependency")
	fmt.Println("  updates with tools like Renovate are crucial for security.")
	fmt.Println()
	fmt.Printf("%s Reference: https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-49568\n", blue("🔗"))
	
	return nil
}