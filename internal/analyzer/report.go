package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Report represents the analysis report
type Report struct {
	RepoInfo  *RepositoryInfo `json:"repository_info"`
	Timestamp time.Time       `json:"timestamp"`
	ToolInfo  ToolInfo        `json:"tool_info"`
}

// ToolInfo contains information about the analysis tool
type ToolInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// NewReport creates a new analysis report
func NewReport(repoInfo *RepositoryInfo) *Report {
	return &Report{
		RepoInfo:  repoInfo,
		Timestamp: time.Now(),
		ToolInfo: ToolInfo{
			Name:        "Git Repository Security Analyzer",
			Version:     "1.0.0",
			Description: "Demonstrates CVE-2023-49568 vulnerability in go-git library",
		},
	}
}

// OutputConsole prints the report to console with colored output
func (r *Report) OutputConsole() error {
	// Color functions
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

	fmt.Println()
	fmt.Printf("%s Git Repository Analysis Report\n", blue("🔍"))
	fmt.Printf("%s %s\n", blue("═"), strings.Repeat("═", 50))
	fmt.Println()

	// Repository Information
	fmt.Printf("%s Repository Information\n", cyan("📁"))
	fmt.Printf("   URL: %s\n", r.RepoInfo.URL)
	fmt.Printf("   Branches: %s\n", green(fmt.Sprintf("%d", r.RepoInfo.BranchCount)))
	fmt.Printf("   Commits Analyzed: %s\n", green(fmt.Sprintf("%d", r.RepoInfo.CommitCount)))
	fmt.Printf("   Contributors: %s\n", green(fmt.Sprintf("%d", len(r.RepoInfo.Contributors))))
	fmt.Println()

	// Last Commit Information
	fmt.Printf("%s Latest Commit\n", magenta("📝"))
	fmt.Printf("   Hash: %s\n", r.RepoInfo.LastCommitHash[:12]+"...")
	fmt.Printf("   Author: %s\n", r.RepoInfo.LastCommitAuthor)
	fmt.Printf("   Date: %s\n", r.RepoInfo.LastCommitDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Message: %s\n", r.RepoInfo.LastCommitMsg)
	fmt.Println()

	// Programming Languages
	if len(r.RepoInfo.Languages) > 0 {
		fmt.Printf("%s Programming Languages Detected\n", blue("💻"))
		for _, lang := range r.RepoInfo.Languages {
			fmt.Printf("   • %s\n", lang)
		}
		fmt.Println()
	}

	// Top Contributors
	if len(r.RepoInfo.Contributors) > 0 {
		fmt.Printf("%s Contributors\n", green("👥"))
		maxShow := 5
		if len(r.RepoInfo.Contributors) < maxShow {
			maxShow = len(r.RepoInfo.Contributors)
		}
		for i := 0; i < maxShow; i++ {
			fmt.Printf("   • %s\n", r.RepoInfo.Contributors[i])
		}
		if len(r.RepoInfo.Contributors) > maxShow {
			fmt.Printf("   ... and %d more\n", len(r.RepoInfo.Contributors)-maxShow)
		}
		fmt.Println()
	}

	// Vulnerability Information (The key part of this demo)
	fmt.Printf("%s SECURITY VULNERABILITY DEMONSTRATION\n", red("🚨"))
	fmt.Printf("%s %s\n", red("═"), strings.Repeat("═", 50))
	fmt.Println()
	
	vuln := r.RepoInfo.VulnerabilityInfo
	fmt.Printf("%s Vulnerability: %s\n", red("🔒"), red(vuln.CVE))
	fmt.Printf("%s Severity: %s\n", red("⚠"), red(vuln.Severity))
	fmt.Printf("%s Affected Library: %s\n", yellow("📦"), vuln.AffectedLib)
	fmt.Printf("%s Current Version: %s %s\n", red("🔴"), vuln.CurrentVer, red("(VULNERABLE)"))
	fmt.Printf("%s Fixed in Version: %s %s\n", green("🟢"), vuln.FixedInVer, green("(SECURE)"))
	fmt.Println()
	
	fmt.Printf("%s Description:\n", blue("📋"))
	fmt.Printf("   %s\n", vuln.Description)
	fmt.Println()

	// Impact and Remediation
	fmt.Printf("%s Potential Impact:\n", red("⚠"))
	fmt.Printf("   • %s\n", "Unauthorized file system access during Git operations")
	fmt.Printf("   • %s\n", "Potential data exfiltration through path traversal")
	fmt.Printf("   • %s\n", "Compromise of application security boundaries")
	fmt.Println()

	fmt.Printf("%s Remediation:\n", green("🛡"))
	fmt.Printf("   • %s\n", "Update go-git library to version 5.11.0 or later")
	fmt.Printf("   • %s\n", "Enable Renovate to automatically detect and fix such vulnerabilities")
	fmt.Printf("   • %s\n", "Implement regular security audits of dependencies")
	fmt.Println()

	// Renovate Information
	fmt.Printf("%s Renovate Integration\n", cyan("🤖"))
	fmt.Printf("%s %s\n", cyan("═"), strings.Repeat("═", 50))
	fmt.Println()
	fmt.Printf("%s This project demonstrates how Renovate can help:\n", green("✅"))
	fmt.Printf("   • %s\n", "Automatically detect vulnerable dependencies")
	fmt.Printf("   • %s\n", "Create pull requests to update to secure versions")
	fmt.Printf("   • %s\n", "Maintain up-to-date security posture")
	fmt.Printf("   • %s\n", "Reduce manual overhead of dependency management")
	fmt.Println()

	// Analysis Metadata
	fmt.Printf("%s Analysis Metadata\n", blue("ℹ"))
	fmt.Printf("   Tool: %s v%s\n", r.ToolInfo.Name, r.ToolInfo.Version)
	fmt.Printf("   Timestamp: %s\n", r.Timestamp.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("   Purpose: %s\n", r.ToolInfo.Description)
	fmt.Println()

	return nil
}

// OutputJSON prints the report as JSON
func (r *Report) OutputJSON() error {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// SaveToFile saves the report to a file
func (r *Report) SaveToFile(filename string, format string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	switch format {
	case "json":
		jsonData, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal report to JSON: %w", err)
		}
		_, err = file.Write(jsonData)
		return err
	case "text":
		// Redirect console output to file
		oldStdout := os.Stdout
		os.Stdout = file
		err := r.OutputConsole()
		os.Stdout = oldStdout
		return err
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}