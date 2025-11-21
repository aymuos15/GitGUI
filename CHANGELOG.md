# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.0] - 2025-11-21

### Added
- Untracked file detection and display in diff and stats views
- `ReadUntrackedFiles()` function to fetch untracked files via git ls-files
- Curl-based installation script for easy setup
- Installation instructions in README
- Color-coded "Untracked" status (magenta) in stats view
- Placeholder message when viewing untracked files in diff view

### Changed
- Updated README with simple curl installation command
- Stats view now includes untracked files with 0 additions/deletions

### Notes
- Currently supports Linux binaries (x86_64, ARM64)
- Installation via: `curl -fsSL https://raw.githubusercontent.com/aymuos15/diffview/main/install.sh | bash`
