import blessed from 'blessed';
import { FileDiff, DiffResult } from './parser.js';

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
      height: 4,
      padding: {
        left: 2,
        right: 2,
      },
      style: {
        bg: '#ffffff',
        fg: '#333333',
        border: {
          fg: '#e0e0e0',
        },
      },
      border: 'bottom',
    });

    // Sidebar with file list
    this.sidebarBox = blessed.box({
      parent: this.screen,
      top: 4,
      left: 0,
      width: 30,
      bottom: 2,
      scrollable: true,
      mouse: true,
      keys: true,
      style: {
        bg: '#fafafa',
        fg: '#333333',
        border: {
          fg: '#e0e0e0',
        },
      },
      border: 'right',
      padding: {
        left: 1,
        right: 1,
      },
    });

    // Main content area
    this.contentBox = blessed.box({
      parent: this.screen,
      top: 4,
      left: 31,
      right: 0,
      bottom: 2,
      scrollable: true,
      mouse: true,
      keys: true,
      vi: true,
      alwaysScroll: true,
      style: {
        bg: '#ffffff',
        fg: '#333333',
      },
      padding: {
        left: 2,
        right: 2,
      },
    });

    // Footer with help text
    this.footerBox = blessed.box({
      parent: this.screen,
      bottom: 0,
      left: 0,
      right: 0,
      height: 2,
      padding: {
        left: 2,
        right: 2,
      },
      style: {
        bg: '#333333',
        fg: '#ffffff',
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
    const path = file.newPath || file.oldPath;

    let header = '{bold}📊 Git Diff Viewer{/bold}\n';
    header += `{gray}${summary.filesChanged} file${summary.filesChanged !== 1 ? 's' : ''} changed  •  `;
    header += `{green}+${summary.insertions}{/green}  `;
    header += `{red}-${summary.deletions}{/red}\n\n`;
    header += `{bold}Current:{/bold} {cyan}${path}{/cyan}`;

    this.headerBox.setContent(header);
  }

  private renderSidebar() {
    let content = '{bold}Files ({gray}' + this.files.length + '{/gray}){/bold}\n';
    content += '{gray}' + '─'.repeat(26) + '{/gray}\n\n';

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

      let line = '';
      if (isActive) {
        line = '{inverse}';
      }

      line += `${statusIcon[file.status]} ${fileName}`;
      line += ` {gray}(+${addCount},-${removeCount}){/gray}`;

      if (isActive) {
        line += '{/inverse}';
      }

      content += line + '\n';
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

    content += `{bold}${statusEmoji[file.status]} ${statusText[file.status]}{/bold}\n`;
    content += `{gray}${file.newPath || file.oldPath}{/gray}\n\n`;

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

      // Render side-by-side with better spacing
      const colWidth = 55;
      const maxLines = Math.max(oldLines.length, newLines.length);

      for (let i = 0; i < maxLines; i++) {
        const oldLine = oldLines[i] || { lineNum: undefined, content: '', type: 'empty' };
        const newLine = newLines[i] || { lineNum: undefined, content: '', type: 'empty' };

        const oldFormatted = this.formatLine(oldLine, colWidth);
        const newFormatted = this.formatLine(newLine, colWidth);

        content += `${oldFormatted} {gray}│{/gray} ${newFormatted}\n`;
      }

      content += '\n';
    }

    this.contentBox.setContent(content);
  }

  private formatLine(
    line: { lineNum?: number; content: string; type: string },
    width: number
  ): string {
    const lineNum = line.lineNum ? String(line.lineNum).padStart(4) : '    ';
    const marker =
      line.type === 'remove' ? '{red}−{/red}' : line.type === 'add' ? '{green}+{/green}' : ' ';

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
    const stripped = colored.replace(/{[^}]+}/g, '');
    if (stripped.length > width) {
      colored = colored.substring(0, width - 1) + '…';
    }

    return colored.padEnd(width);
  }

  private renderFooter() {
    const current = this.currentFileIndex + 1;
    const total = this.files.length;

    let footer = '';
    footer += '{bold}[n]{/bold}ext  ';
    footer += '{bold}[p]{/bold}rev  ';
    footer += '{bold}[j/k]{/bold} scroll  ';
    footer += '{bold}[?]{/bold} help  ';
    footer += '{bold}[q]{/bold} quit';
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
          fg: '#666',
        },
        bg: '#ffffff',
        fg: '#333333',
      },
      padding: {
        left: 2,
        right: 2,
        top: 1,
        bottom: 1,
      },
    });

    const helpText = `{bold}Keyboard Shortcuts{/bold}

{bold}Navigation{/bold}
  {cyan}n{/cyan}  or  {cyan}→{/cyan}      Next file
  {cyan}p{/cyan}  or  {cyan}←{/cyan}      Previous file
  {cyan}j{/cyan}  or  {cyan}↓{/cyan}      Scroll down
  {cyan}k{/cyan}  or  {cyan}↑{/cyan}      Scroll up

{bold}Scrolling{/bold}
  {cyan}PageDown{/cyan}       Scroll page down
  {cyan}PageUp{/cyan}         Scroll page up

{bold}Other{/bold}
  {cyan}?{/cyan}  or  {cyan}h{/cyan}    Show this help
  {cyan}q{/cyan}           Quit

{gray}Press any key to close{/gray}`;

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
