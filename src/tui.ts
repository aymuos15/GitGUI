import blessed from 'blessed';
import { FileDiff, DiffResult } from './parser.js';
import { parseMarkup } from './colors.js';

export class DiffTUI {
  private screen: any;
  private currentFileIndex: number = 0;
  private files: FileDiff[] = [];
  private diffResult: DiffResult;
  private contentBox: any;
  private headerBox: any;
  private footerBox: any;
  private sidebarBox: any;

  constructor(diffResult: DiffResult) {
    this.diffResult = diffResult;
    this.files = diffResult.files;

    // Create blessed screen with better default style
    this.screen = blessed.screen({
      mouse: true,
      keyboard: true,
      title: 'Git Diff Viewer',
      smartCSR: true,
      trueColor: true,
      style: {
        bg: '#f8f8f8',
        fg: '#333333',
      },
    });

    // Header with summary
    this.headerBox = blessed.box({
      parent: this.screen,
      top: 0,
      left: 0,
      right: 0,
      height: 3,
      padding: {
        left: 2,
        right: 2,
        top: 1,
        bottom: 0,
      },
      style: {
        bg: 'white',
        fg: 'black',
        border: {
          fg: 'gray',
        },
      },
      border: 'bottom',
      tags: true,
    });

    // Sidebar with file list
    this.sidebarBox = blessed.box({
      parent: this.screen,
      top: 3,
      left: 0,
      width: '20%',
      bottom: 1,
      scrollable: true,
      mouse: true,
      keys: true,
      style: {
        bg: 'white',
        fg: 'black',
        border: {
          fg: 'gray',
        },
      },
      border: 'right',
      padding: {
        left: 1,
        right: 1,
      },
      tags: true,
    });

    // Main content area
    this.contentBox = blessed.box({
      parent: this.screen,
      top: 3,
      left: '20%',
      right: 0,
      bottom: 1,
      scrollable: true,
      mouse: true,
      keys: true,
      vi: true,
      alwaysScroll: true,
      style: {
        bg: 'white',
        fg: 'black',
      },
      padding: {
        left: 1,
        right: 1,
      },
      tags: true,
    });

    // Footer with help text
    this.footerBox = blessed.box({
      parent: this.screen,
      bottom: 0,
      left: 0,
      right: 0,
      height: 1,
      padding: {
        left: 2,
        right: 2,
      },
      style: {
        bg: 'black',
        fg: 'white',
      },
      tags: true,
    });

    this.setupKeyBindings();
    this.render();
  }

  private setupKeyBindings() {
    this.screen.key(['n'], () => this.nextFile());
    this.screen.key(['p'], () => this.prevFile());
    this.screen.key(['j'], () => this.contentBox.scroll(1));
    this.screen.key(['k'], () => this.contentBox.scroll(-1));
    this.screen.key(['left'], () => this.prevFile());
    this.screen.key(['right'], () => this.nextFile());
    this.screen.key(['h', '?'], () => this.showHelp());
    this.screen.key(['q', 'C-c'], () => process.exit(0));
    this.screen.key(['pagedown'], () => this.contentBox.scroll(10));
    this.screen.key(['pageup'], () => this.contentBox.scroll(-10));
  }

  private nextFile() {
    if (this.currentFileIndex < this.files.length - 1) {
      this.currentFileIndex++;
      this.contentBox.setScrollPerc(0);
      this.render();
    }
  }

  private prevFile() {
    if (this.currentFileIndex > 0) {
      this.currentFileIndex--;
      this.contentBox.setScrollPerc(0);
      this.render();
    }
  }

  private render() {
    this.renderHeader();
    this.renderSidebar();
    this.renderContent();
    this.renderFooter();
    this.screen.render();
  }

  private renderHeader() {
    const summary = this.diffResult.summary;
    const file = this.files[this.currentFileIndex];
    const path = file ? (file.newPath || file.oldPath) : 'N/A';

    let header = '\x1b[1m📊 Git Diff Viewer\x1b[0m\n';
    header += `${summary.filesChanged} file${summary.filesChanged !== 1 ? 's' : ''} changed  •  `;
    header += `\x1b[32m+${summary.insertions}\x1b[0m  `;
    header += `\x1b[31m-${summary.deletions}\x1b[0m\n\n`;
    header += `\x1b[1mCurrent:\x1b[0m \x1b[34m${path}\x1b[0m`;

    this.headerBox.setContent(header);
  }

  private renderSidebar() {
    if (this.files.length === 0) {
      this.sidebarBox.setContent('\x1b[1mFiles\x1b[0m\n(no changes)');
      return;
    }

    let content = '\x1b[1mFiles\x1b[0m\n'; // bold Files

    for (let i = 0; i < this.files.length; i++) {
      const file = this.files[i];
      const isActive = i === this.currentFileIndex;
      const fileName = (file.newPath || file.oldPath).split('/').pop() || file.newPath || file.oldPath;

      const addCount = file.hunks
        .flatMap((h: any) => h.lines)
        .filter((l: any) => l.type === 'add').length;
      const removeCount = file.hunks
        .flatMap((h: any) => h.lines)
        .filter((l: any) => l.type === 'remove').length;

      const statusIcon: { [key: string]: string } = {
        added: '✚',
        modified: '◆',
        deleted: '✕',
        renamed: '➜',
      };

      // Build the line - use ANSI codes directly
      let line = '';
      if (isActive) {
        line = '\x1b[7m'; // inverse
      }
      
      line += `${statusIcon[file.status]} `;
      
      // Truncate filename to fit
      const maxLen = 12;
      if (fileName.length > maxLen) {
        line += fileName.substring(0, maxLen - 1) + '…';
      } else {
        line += fileName;
      }

      if (isActive) {
        line += '\x1b[0m'; // reset
      }

      content += line + '\n';
      content += `  \x1b[32m+${addCount}\x1b[0m \x1b[31m-${removeCount}\x1b[0m\n`;
    }

    this.sidebarBox.setContent(content);
  }

  private renderContent() {
    if (this.files.length === 0) return;

    const file = this.files[this.currentFileIndex];
    let content = '';

    const statusEmoji: { [key: string]: string } = {
      added: '✚',
      modified: '◆',
      deleted: '✕',
      renamed: '➜',
    };

    const statusText: { [key: string]: string } = {
      added: 'ADDED',
      modified: 'MODIFIED',
      deleted: 'DELETED',
      renamed: 'RENAMED',
    };

    content += `\x1b[1m${statusEmoji[file.status]} ${statusText[file.status]}\x1b[0m\n`;
    content += `${file.newPath || file.oldPath}\n`;

    // Calculate column width early
    const terminalWidth = this.screen.width || 120;
    const sidebarWidth = Math.floor(terminalWidth * 0.2); // 20% for sidebar
    const padding = 2; // left + right padding on content box
    const separatorWidth = 3; // " │ "
    const availableWidth = Math.max(50, terminalWidth - sidebarWidth - padding - separatorWidth);
    const colWidth = Math.max(20, Math.floor(availableWidth / 2));

    // Render hunks
    for (const hunk of file.hunks) {
      content += `\n\x1b[34m@@ -${hunk.oldStart},${hunk.oldCount} +${hunk.newStart},${hunk.newCount} @@\x1b[0m\n`;

      // Separate old and new lines for side-by-side display
      let oldLines: Array<{ lineNum?: number; content: string; type: string }> = [];
      let newLines: Array<{ lineNum?: number; content: string; type: string }> = [];

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

      const maxLines = Math.max(oldLines.length, newLines.length);

      for (let i = 0; i < maxLines; i++) {
        const oldLine = oldLines[i] || { lineNum: undefined, content: '', type: 'empty' };
        const newLine = newLines[i] || { lineNum: undefined, content: '', type: 'empty' };

        const oldFormatted = this.formatLine(oldLine, colWidth);
        const newFormatted = this.formatLine(newLine, colWidth);

        content += `${oldFormatted} │ ${newFormatted}\n`;
      }
    }

    this.contentBox.setContent(content);
  }

  private formatLine(
    line: { lineNum?: number; content: string; type: string },
    width: number
  ): string {
    // Handle empty lines
    if (line.type === 'empty') {
      return ' '.repeat(width);
    }

    const lineNum = line.lineNum ? String(line.lineNum).padStart(4) : '    ';
    const marker = line.type === 'remove' ? '−' : line.type === 'add' ? '+' : ' ';
    const text = `${lineNum} ${marker} ${line.content}`;

    // Truncate to width if needed (accounting for actual character count, not markup)
    let displayText = text;
    if (text.length > width) {
      displayText = text.substring(0, width - 1) + '…';
    }

    // Pad to exact width
    displayText = displayText.padEnd(width, ' ');

    // Apply color based on type using ANSI codes
    if (line.type === 'remove') {
      return `\x1b[31m${displayText}\x1b[0m`;
    } else if (line.type === 'add') {
      return `\x1b[32m${displayText}\x1b[0m`;
    } else {
      return displayText;
    }
  }

  private renderFooter() {
    const current = this.currentFileIndex + 1;
    const total = this.files.length;

    let footer = '';
    footer += '\x1b[1m[n]\x1b[0mext  ';
    footer += '\x1b[1m[p]\x1b[0mrev  ';
    footer += '\x1b[1m[j/k]\x1b[0m scroll  ';
    footer += '\x1b[1m[?]\x1b[0m help  ';
    footer += '\x1b[1m[q]\x1b[0m quit';
    footer += '                                          ';
    footer += `File ${current}/${total}`;

    this.footerBox.setContent(footer);
  }

  private showHelp() {
    const helpBox = blessed.box({
      parent: this.screen,
      top: 'center',
      left: 'center',
      width: 70,
      height: 24,
      border: 'line',
      style: {
        border: {
          fg: 'blue',
        },
        bg: 'white',
        fg: 'black',
      },
      padding: {
        left: 2,
        right: 2,
        top: 1,
        bottom: 1,
      },
    });

    const helpText = `\x1b[1mKeyboard Shortcuts\x1b[0m

\x1b[1mNavigation\x1b[0m
  \x1b[34mn\x1b[0m  or  \x1b[34m→\x1b[0m      Next file
  \x1b[34mp\x1b[0m  or  \x1b[34m←\x1b[0m      Previous file
  \x1b[34mj\x1b[0m  or  \x1b[34m↓\x1b[0m      Scroll down
  \x1b[34mk\x1b[0m  or  \x1b[34m↑\x1b[0m      Scroll up

\x1b[1mScrolling\x1b[0m
  \x1b[34mPageDown\x1b[0m       Scroll page down
  \x1b[34mPageUp\x1b[0m         Scroll page up

\x1b[1mOther\x1b[0m
  \x1b[34m?\x1b[0m  or  \x1b[34mh\x1b[0m    Show this help
  \x1b[34mq\x1b[0m           Quit

Press any key to close`;

    helpBox.setContent(helpText);
    this.screen.render();

    this.screen.once('key', () => {
      helpBox.destroy();
      this.render();
    });
  }

  public show() {
    this.screen.render();
  }
}
