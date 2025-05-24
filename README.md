


## 🔐 Go Security Renovate Demo

A demonstration tool built in Go that analyzes Git repositories and highlights the importance of dependency security management using **Renovate**.  
This project intentionally uses a **vulnerable version** of the [`go-git`](https://github.com/go-git/go-git) library to demonstrate real-world CVE detection and resolution.

---

## 📌 Purpose

This project is designed to:
- Demonstrate how vulnerable dependencies can impact Go projects.
- Showcase [Renovate](https://github.com/renovatebot/renovate) as a powerful tool to automatically detect and fix such vulnerabilities.
- Provide a real-world example of using Renovate in a Go module that depends on a package with a known security vulnerability.

---

## ⚠️ Vulnerable Dependency

- **Package:** `github.com/go-git/go-git/v5`
- **Current Version Used:** `v5.4.2`
- **Known Vulnerability:** [CVE-2023-49568](https://nvd.nist.gov/vuln/detail/CVE-2023-49568)
- **Severity:** High
- **Issue:** Path traversal vulnerability that may allow unauthorized file system access during Git operations.
- **Fixed In:** `v5.11.0`

---

## 🛠 Features

- Clone and analyze any Git repository.
- Extract repository insights: commits, contributors, branches, languages.
- Simulate the impact of using a vulnerable dependency.
- Generate structured analysis reports with vulnerability metadata.

---

## 🚀 Quick Start

### Prerequisites
- Go 1.18 or higher installed on your machine.

### Clone the Repository
```bash
git clone https://github.com/Jay2006sawant/go-security-renovate-demo.git
cd go-security-renovate-demo
````

### Run the Analyzer

```bash
go run ./cmd/analyzer demo
```

This command runs a demo that:

* Analyzes a few real-world GitHub repositories.
* Demonstrates how the vulnerable `go-git` library is used.
* Outputs a detailed security and dependency report.

---

## 📋 Example Output

![image](https://github.com/user-attachments/assets/f25899b9-b329-4887-a5a6-32c575836c8e)




![image](https://github.com/user-attachments/assets/5a5d1b02-4a60-4b7d-b5e1-8dd1bdf03b75)



![image](https://github.com/user-attachments/assets/477b0176-b34e-41b1-a0d6-932b4bbeb64f)


---

## 🤖 Renovate Integration

By integrating Renovate with this project, you can:

* Automatically detect outdated and insecure dependencies.
* Get pull requests to update to secure versions.
* Maintain a strong and proactive security posture.


---

## 📁 Project Structure

```
go-security-renovate-demo/
│
├── cmd/
│   └── analyzer/         # CLI logic for the analyzer tool
│
├── internal/
│   ├── analyzer/         # Repository analysis logic
│   └── report/           # Report formatting and display
│
├── go.mod                # Go module file with vulnerable dependency
├── renovate.json         # Renovate configuration for automated updates
└── README.md             # Project documentation
```

---



```
