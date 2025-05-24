package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitAnalyzer handles Git repository analysis using the vulnerable go-git library
type GitAnalyzer struct {
	tempDir string
}

// RepositoryInfo contains information about the analyzed repository
type RepositoryInfo struct {
	URL               string    `json:"url"`
	LastCommitHash    string    `json:"last_commit_hash"`
	LastCommitDate    time.Time `json:"last_commit_date"`
	LastCommitAuthor  string    `json:"last_commit_author"`
	LastCommitMsg     string    `json:"last_commit_message"`
	BranchCount       int       `json:"branch_count"`
	CommitCount       int       `json:"commit_count"`
	Contributors      []string  `json:"contributors"`
	Languages         []string  `json:"languages"`
	VulnerabilityInfo VulnInfo  `json:"vulnerability_info"`
}

// VulnInfo contains information about the vulnerability being demonstrated
type VulnInfo struct {
	CVE           string `json:"cve"`
	Severity      string `json:"severity"`
	AffectedLib   string `json:"affected_library"`
	CurrentVer    string `json:"current_version"`
	FixedInVer    string `json:"fixed_in_version"`
	Description   string `json:"description"`
}

// NewGitAnalyzer creates a new GitAnalyzer instance
func NewGitAnalyzer() *GitAnalyzer {
	return &GitAnalyzer{
		tempDir: "/tmp/git-analyzer",
	}
}

// AnalyzeRepository clones and analyzes a Git repository
// This method uses the VULNERABLE go-git library version 5.4.2
// which is susceptible to CVE-2023-49568 (path traversal vulnerability)
func (ga *GitAnalyzer) AnalyzeRepository(repoURL string) (*Report, error) {
	// Create temporary directory for cloning
	cloneDir := filepath.Join(ga.tempDir, fmt.Sprintf("repo-%d", time.Now().Unix()))

	// Clean up clone directory when done
	defer func() {
		os.RemoveAll(cloneDir)
	}()

	fmt.Printf("🔄 Cloning repository (using VULNERABLE go-git v5.4.2)...\n")

	// Clone repository using vulnerable go-git library
	// CVE-2023-49568: This version is vulnerable to path traversal attacks
	repo, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: nil, // Suppress progress for cleaner output
		Depth:    50,  // Shallow clone for faster analysis
	})
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Printf("✅ Repository cloned successfully\n")

	// Analyze repository structure and commits
	repoInfo, err := ga.analyzeRepoStructure(repo, repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository structure: %w", err)
	}

	// Analyze files for language detection
	languages, err := ga.detectLanguages(cloneDir)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not detect languages: %v\n", err)
	}
	repoInfo.Languages = languages

	// Create vulnerability information
	repoInfo.VulnerabilityInfo = VulnInfo{
		CVE:         "CVE-2023-49568",
		Severity:    "HIGH",
		AffectedLib: "github.com/go-git/go-git/v5",
		CurrentVer:  "5.4.2",
		FixedInVer:  "5.11.0",
		Description: "Path traversal vulnerability allowing unauthorized file system access during Git operations",
	}

	return NewReport(repoInfo), nil
}

// analyzeRepoStructure extracts information from the Git repository
func (ga *GitAnalyzer) analyzeRepoStructure(repo *git.Repository, repoURL string) (*RepositoryInfo, error) {
	info := &RepositoryInfo{
		URL: repoURL,
	}

	// Get HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	info.LastCommitHash = ref.Hash().String()

	// Get last commit information
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit object: %w", err)
	}

	info.LastCommitDate = commit.Author.When
	info.LastCommitAuthor = commit.Author.Name
	info.LastCommitMsg = strings.Split(commit.Message, "\n")[0] // First line only

	// Count commits (limited for performance)
	commitCount, contributors := ga.countCommitsAndContributors(repo)
	info.CommitCount = commitCount
	info.Contributors = contributors

	// Count branches
	branches, err := repo.Branches()
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not count branches: %v\n", err)
		info.BranchCount = 1 // At least main/master branch
	} else {
		branchCount := 0
		branches.ForEach(func(ref *plumbing.Reference) error {
			branchCount++
			return nil
		})
		info.BranchCount = branchCount
	}

	return info, nil
}

// countCommitsAndContributors counts commits and extracts unique contributors
func (ga *GitAnalyzer) countCommitsAndContributors(repo *git.Repository) (int, []string) {
	ref, err := repo.Head()
	if err != nil {
		return 0, []string{}
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return 0, []string{}
	}
	defer commitIter.Close()

	commitCount := 0
	contributorMap := make(map[string]bool)
	maxCommits := 100 // Limit for performance

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commitCount >= maxCommits {
			return fmt.Errorf("max commits reached") // Break iteration
		}

		commitCount++
		contributorMap[commit.Author.Name] = true
		return nil
	})

	// Convert map to slice
	contributors := make([]string, 0, len(contributorMap))
	for contributor := range contributorMap {
		contributors = append(contributors, contributor)
	}

	return commitCount, contributors
}

// detectLanguages analyzes files to detect programming languages
func (ga *GitAnalyzer) detectLanguages(repoPath string) ([]string, error) {
	languageMap := make(map[string]bool)

	// Define file extension to language mapping
	extToLang := map[string]string{
		".go":    "Go",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".py":    "Python",
		".java":  "Java",
		".cpp":   "C++",
		".c":     "C",
		".cs":    "C#",
		".php":   "PHP",
		".rb":    "Ruby",
		".rs":    "Rust",
		".kt":    "Kotlin",
		".swift": "Swift",
		".scala": "Scala",
		".sh":    "Shell",
		".yaml":  "YAML",
		".yml":   "YAML",
		".json":  "JSON",
		".xml":   "XML",
		".html":  "HTML",
		".css":   "CSS",
		".md":    "Markdown",
	}

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on errors
		}

		if info.IsDir() {
			// Skip hidden directories and common non-source directories
			dirName := info.Name()
			if strings.HasPrefix(dirName, ".") ||
				dirName == "node_modules" ||
				dirName == "vendor" ||
				dirName == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if lang, exists := extToLang[ext]; exists {
			languageMap[lang] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to slice
	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}

	return languages, nil
}
