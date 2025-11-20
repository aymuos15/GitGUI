#!/usr/bin/env python3
"""A side-by-side diff viewer for git diff output."""

import sys
from typing import List, Tuple
from rich.console import Console
from rich.table import Table
from rich.text import Text
from rich.panel import Panel


class DiffLine:
    """Represents a single line in a diff."""

    def __init__(self, content: str, line_type: str, line_num: int = 0):
        self.content = content.rstrip("\n")
        self.line_type = line_type  # 'add', 'remove', 'context', 'empty'
        self.line_num = line_num

    def __repr__(self):
        return f"DiffLine({self.line_type}, {self.line_num})"


class GitDiffParser:
    """Parse unified diff format from git diff."""

    @staticmethod
    def parse(diff_text: str) -> Tuple[List[DiffLine], List[DiffLine]]:
        """Parse unified diff into side-by-side format."""
        left_lines: List[DiffLine] = []
        right_lines: List[DiffLine] = []

        left_num = 1
        right_num = 1

        for line in diff_text.split("\n"):
            # Skip file headers
            if line.startswith("---") or line.startswith("+++"):
                continue
            if line.startswith("diff ") or line.startswith("index "):
                continue

            # Skip hunk headers
            if line.startswith("@@"):
                continue

            # Process diff lines
            if line.startswith("-") and not line.startswith("---"):
                # Removed line
                left_lines.append(DiffLine(line[1:], "remove", left_num))
                right_lines.append(DiffLine("", "empty", 0))
                left_num += 1

            elif line.startswith("+") and not line.startswith("+++"):
                # Added line
                left_lines.append(DiffLine("", "empty", 0))
                right_lines.append(DiffLine(line[1:], "add", right_num))
                right_num += 1

            elif line.startswith(" "):
                # Context line
                left_lines.append(DiffLine(line[1:], "context", left_num))
                right_lines.append(DiffLine(line[1:], "context", right_num))
                left_num += 1
                right_num += 1

            elif line == "":
                # Empty line - skip
                pass
            elif line.startswith("\\"):
                # "\ No newline at end of file" - skip
                pass

        return left_lines, right_lines


class DiffViewer:
    """Display side-by-side diff in terminal."""

    def __init__(self, diff_text: str):
        self.console = Console()
        self.diff_text = diff_text
        self.left_lines, self.right_lines = GitDiffParser.parse(diff_text)

    def display(self) -> None:
        """Display the diff in a formatted table."""
        # Create table with two columns
        table = Table.grid(padding=(0, 1))
        table.add_column(style="dim")  # Left file
        table.add_column(style="dim")  # Right file

        # Add lines to table
        max_lines = max(len(self.left_lines), len(self.right_lines))

        for i in range(max_lines):
            left = self.left_lines[i] if i < len(self.left_lines) else None
            right = self.right_lines[i] if i < len(self.right_lines) else None

            left_text = self._format_line(left)
            right_text = self._format_line(right)

            table.add_row(left_text, right_text)

        # Display the table
        self.console.print(table)

    def _format_line(self, line: DiffLine | None) -> Text:
        """Format a single line with color and line number."""
        if line is None or line.line_type == "empty":
            return Text("")

        # Format line number
        if line.line_num > 0:
            line_num_str = f"{line.line_num:5d} "
        else:
            line_num_str = "      "

        # Format content
        content = line.content
        if not content:
            content = "∅"

        # Apply color based on type
        if line.line_type == "add":
            return Text(line_num_str + content, style="green")
        elif line.line_type == "remove":
            return Text(line_num_str + content, style="red")
        else:  # context
            return Text(line_num_str + content, style="dim")


def main():
    """Main entry point."""
    # Read from stdin
    if not sys.stdin.isatty():
        diff_text = sys.stdin.read()

        if not diff_text.strip():
            print("No diff provided. Usage: git diff | diffview.py")
            sys.exit(1)

        viewer = DiffViewer(diff_text)
        viewer.display()
    else:
        print("Usage: git diff | python3 diffview.py")
        print()
        print("Examples:")
        print("  git diff | python3 diffview.py")
        print("  git diff --staged | python3 diffview.py")
        print("  git diff HEAD~1 | python3 diffview.py")
        print("  diff file1 file2 | python3 diffview.py")
        sys.exit(1)


if __name__ == "__main__":
    main()
