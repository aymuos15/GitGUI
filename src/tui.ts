import blessed from 'blessed';
import { FileDiff, DiffResult } from './parser.js';

export class DiffTUI {
  private screen: any;
  private currentFileIndex: number = 0;
  private files: FileDiff[] = [];
  private diffResult: DiffResult;
  private contentBox: any;
  private tabBar: any;
  private statusBar: any;

  constructor(diffResult: DiffResult) {
    this.diffResult = diffResult;
    this.files = diffResult.files;

    // Create blessed screen
    this.screen = blessed.screen({
      mouse: true,
      keyboard: true,
      title: 'Git Diff Viewer',
    });

    // Create layout boxes
    this.tabBar = blessed.box({
      parent: this.screen,
      top: 0,
      left: 0,
      right: 0,
      height: 2,
      style: {
        bg: 'blue',
        fg: 'white',
      },
    });

    this.contentBox = blessed.box({
      parent: this.screen,
      top: 2,
      left: 0,
      right: 0,
      bottom: 2,
      scrollable: true,
      mouse: true,
      keys: true,
      vi: true,
      alwaysScroll: true,
      style: {
        bg: 'black',
        fg: 'white',
      },
    });

    this.statusBar = blessed.box({
      parent: this.screen,
      bottom: 0,
      left: 0,
      right: 0,
      height: 2,
      style: {
        bg: 'black',
        fg: 'cyan',
      },
    });

    this.setupKeyBindings();
    this.render();
  }

  private setupKeyBindings() {
    this.screen.key(['n'], () => this.nextFile());
    this.screen.key(['p'], () => this.prevFile());
    this.screen.key(['j'], () => this.contentBox.scroll(1));
    this.screen.key(['k'], () => this.contentBox.scroll(-1));
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
    this.renderTabBar();
    this.renderContent();
    this.renderStatusBar();
    this.screen.render();
  }

  private renderTabBar() {
    let tabContent = '';
    const summary = this.diffResult.summary;

    // Header line with summary
    const headerLine = ` 📊 Git Diff: ${summary.filesChanged} file${
      summary.filesChanged !== 1 ? 's' : ''
    } changed, +${summary.insertions} -${summary.deletions}`;
    tabContent += headerLine + '\n';

    // Tab line
    const tabs: string[] = [];
    for (let i = 0; i < this.files.length; i++) {
      const file = this.files[i];
      const fileName = (file.newPath || file.oldPath).split('/').pop() || file.newPath || file.oldPath;
      const isActive = i === this.currentFileIndex;

      if (isActive) {
        tabs.push(`{bold}[ ${fileName} ]{/bold}`);
      } else {
        tabs.push(` ${fileName} `);
      }
    }

    tabContent += tabs.join('  ');
    this.tabBar.setContent(tabContent);
  }

  private renderContent() {
    if (this.files.length === 0) return;

    const file = this.files[this.currentFileIndex];
    let content = '';

    // File header
    const statusEmoji: { [key: string]: string } = {
      added: '📄',
      modified: '✏️ ',
      deleted: '🗑️ ',
      renamed: '➜ ',
    };
    const statusText: { [key: string]: string } = {
      added: 'ADDED',
      modified: 'MODIFIED',
      deleted: 'DELETED',
      renamed: 'RENAMED',
    };

    content += `{bold}${statusEmoji[file.status]} ${statusText[file.status]} ${file.newPath || file.oldPath}{/bold}\n\n`;

    // Render hunks
    for (const hunk of file.hunks) {
      content += `{cyan}@@ -${hunk.oldStart},${hunk.oldCount} +${hunk.newStart},${hunk.newCount} @@{/cyan}\n`;

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

      // Render side-by-side
      const colWidth = 50;
      const maxLines = Math.max(oldLines.length, newLines.length);

      for (let i = 0; i < maxLines; i++) {
        const oldLine = oldLines[i] || { lineNum: undefined, content: '', type: 'empty' };
        const newLine = newLines[i] || { lineNum: undefined, content: '', type: 'empty' };

        const oldFormatted = this.formatLine(oldLine, colWidth);
        const newFormatted = this.formatLine(newLine, colWidth);

        content += `${oldFormatted} │ ${newFormatted}\n`;
      }

      content += '\n';
    }

    this.contentBox.setContent(content);
  }

  private formatLine(
    line: { lineNum?: number; content: string; type: string },
    width: number
  ): string {
    const lineNum = line.lineNum ? String(line.lineNum).padStart(4) : '   ';
    const marker = line.type === 'remove' ? '−' : line.type === 'add' ? '+' : ' ';

    let colored = `${lineNum} ${marker} ${line.content}`;

    if (line.type === 'remove') {
      colored = `{red}${colored}{/red}`;
    } else if (line.type === 'add') {
      colored = `{green}${colored}{/green}`;
    } else if (line.type === 'context') {
      colored = `{gray}${colored}{/gray}`;
    } else {
      colored = ' '.repeat(width);
    }

    // Truncate to width
    if (colored.length > width) {
      colored = colored.substring(0, width - 1) + '…';
    }
    return colored.padEnd(width);
  }

  private renderStatusBar() {
    const current = this.currentFileIndex + 1;
    const total = this.files.length;
    const file = this.files[this.currentFileIndex];
    const addCount = file.hunks
      .flatMap((h) => h.lines)
      .filter((l) => l.type === 'add').length;
    const removeCount = file.hunks
      .flatMap((h) => h.lines)
      .filter((l) => l.type === 'remove').length;

    let status = `File ${current}/${total} | `;
    status += `+${addCount} -${removeCount} | `;
    status += `[n]ext [p]rev [?]help [q]uit [j/k]scroll`;

    this.statusBar.setContent(status);
  }

  private showHelp() {
    const helpBox = blessed.box({
      parent: this.screen,
      top: 'center',
      left: 'center',
      width: 60,
      height: 20,
      border: 'line',
      style: {
        border: {
          fg: 'cyan',
        },
        bg: 'black',
        fg: 'white',
      },
    });

    const helpText = `
Navigation:
  n/→        Next file
  p/←        Previous file
  j/↓        Scroll down
  k/↑        Scroll up
  PgDn       Page down
  PgUp       Page up

Controls:
  ?/h        Show this help
  q/Ctrl-C   Quit
`;

    helpBox.setContent(helpText);
    this.screen.render();

    this.screen.once('key', () => {
      helpBox.destroy();
      this.screen.render();
    });
  }

  public show() {
    this.screen.render();
  }
}
