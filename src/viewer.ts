import chalk from 'chalk';
import { FileDiff, DiffHunk, DiffLine } from './parser.js';

interface TerminalSize {
  width: number;
  height: number;
}

export class DiffViewer {
  private terminalWidth: number;
  private terminalHeight: number;

  constructor() {
    this.terminalWidth = process.stdout.columns || 120;
    this.terminalHeight = process.stdout.rows || 40;
  }

  /**
   * Render a side-by-side diff for a single file
   */
  renderFileDiff(file: FileDiff): string {
    const lines: string[] = [];

    // File header
    lines.push(this.renderFileHeader(file));
    lines.push('');

    // Render each hunk
    for (const hunk of file.hunks) {
      lines.push(this.renderHunkHeader(hunk));
      lines.push(this.renderSideBySideDiff(hunk));
      lines.push('');
    }

    return lines.join('\n');
  }

  /**
   * Render all file diffs
   */
  renderAllDiffs(files: FileDiff[]): string {
    const lines: string[] = [];

    for (const file of files) {
      lines.push(this.renderFileDiff(file));
    }

    return lines.join('\n');
  }

  private renderFileHeader(file: FileDiff): string {
    const statusColor = {
      added: chalk.green,
      modified: chalk.yellow,
      deleted: chalk.red,
      renamed: chalk.cyan,
    }[file.status];

    const statusText = {
      added: '📄 ADDED',
      modified: '✏️  MODIFIED',
      deleted: '🗑️  DELETED',
      renamed: '➜  RENAMED',
    }[file.status];

    const path = file.newPath || file.oldPath;
    return statusColor.bold(`${statusText} ${path}`);
  }

  private renderHunkHeader(hunk: DiffHunk): string {
    const oldRange = `${hunk.oldStart},${hunk.oldCount}`;
    const newRange = `${hunk.newStart},${hunk.newCount}`;
    return chalk.cyan(`@@ -${oldRange} +${newRange} @@`);
  }

  /**
   * Render side-by-side diff lines
   */
  private renderSideBySideDiff(hunk: DiffHunk): string {
    const lines: string[] = [];
    const colWidth = Math.floor((this.terminalWidth - 4) / 2);
    const separator = chalk.gray('│');

    // Track lines for pairing
    let oldLines: Array<{ lineNum?: number; content: string; type: string }> = [];
    let newLines: Array<{ lineNum?: number; content: string; type: string }> = [];

    // First pass: separate old and new lines
    for (const line of hunk.lines) {
      if (line.type === 'remove') {
        oldLines.push({
          lineNum: line.oldLineNum,
          content: line.content,
          type: 'remove',
        });
      } else if (line.type === 'add') {
        newLines.push({
          lineNum: line.newLineNum,
          content: line.content,
          type: 'add',
        });
      } else {
        // Context line - display on both sides
        oldLines.push({
          lineNum: line.oldLineNum,
          content: line.content,
          type: 'context',
        });
        newLines.push({
          lineNum: line.newLineNum,
          content: line.content,
          type: 'context',
        });
      }
    }

    // Render pairs of lines side-by-side
    const maxLines = Math.max(oldLines.length, newLines.length);
    for (let i = 0; i < maxLines; i++) {
      const oldLine = oldLines[i] || { lineNum: undefined, content: '', type: 'empty' };
      const newLine = newLines[i] || { lineNum: undefined, content: '', type: 'empty' };

      const oldRendered = this.renderDiffLine(oldLine, colWidth, 'old');
      const newRendered = this.renderDiffLine(newLine, colWidth, 'new');

      lines.push(`${oldRendered}${separator}${newRendered}`);
    }

    return lines.join('\n');
  }

  private renderDiffLine(
    line: { lineNum?: number; content: string; type: string },
    width: number,
    side: 'old' | 'new'
  ): string {
    const lineNum = line.lineNum ? String(line.lineNum).padStart(4) : '   ';
    const marker = this.getLineMarker(line.type);
    const content = this.truncateContent(line.content, width - 8);

    let colored = `${lineNum} ${marker} ${content}`;

    if (line.type === 'remove') {
      colored = chalk.red(colored);
    } else if (line.type === 'add') {
      colored = chalk.green(colored);
    } else if (line.type === 'context') {
      colored = chalk.gray(colored);
    } else if (line.type === 'empty') {
      colored = chalk.dim(' '.repeat(width));
    }

    return colored.padEnd(width);
  }

  private getLineMarker(type: string): string {
    switch (type) {
      case 'remove':
        return '−';
      case 'add':
        return '+';
      case 'context':
        return ' ';
      default:
        return ' ';
    }
  }

  private truncateContent(content: string, maxLen: number): string {
    if (content.length > maxLen) {
      return content.substring(0, maxLen - 1) + '…';
    }
    return content.padEnd(maxLen);
  }
}
