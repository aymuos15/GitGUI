#!/usr/bin/env node

import chalk from 'chalk';
import * as readline from 'readline';
import { getGitDiff } from './parser.js';
import { DiffViewer } from './viewer.js';

interface Options {
  ref?: string;
  staged?: boolean;
  color?: boolean;
}

async function main() {
  const args = process.argv.slice(2);
  const options: Options = {
    color: true,
  };

  // Parse command line arguments
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '--staged' || arg === '-s') {
      options.staged = true;
    } else if (arg === '--no-color') {
      options.color = false;
      chalk.level = 0;
    } else if (arg === '--help' || arg === '-h') {
      printHelp();
      process.exit(0);
    } else if (!arg.startsWith('-')) {
      options.ref = arg;
    }
  }

  try {
    // Get the diff
    let ref = options.ref;
    if (options.staged) {
      ref = '--staged';
    }

    console.log(chalk.blue.bold('\n📊 Git Diff Viewer\n'));

    const diffResult = getGitDiff(ref);

    if (diffResult.files.length === 0) {
      console.log(chalk.yellow('No changes to display'));
      process.exit(0);
    }

    // Print summary
    printSummary(diffResult);
    console.log('');

    // Render diffs
    const viewer = new DiffViewer(diffResult.files);
    const output = viewer.renderAllDiffs(diffResult.files);
    console.log(output);

    // Show navigation hints if multiple files
    if (diffResult.files.length > 1) {
      console.log('');
      printNavigationHints(viewer);
      await interactiveMode(viewer, diffResult.files);
    }
  } catch (error) {
    if (error instanceof Error) {
      console.error(chalk.red(`Error: ${error.message}`));
    } else {
      console.error(chalk.red('An unknown error occurred'));
    }
    process.exit(1);
  }
}

async function interactiveMode(viewer: DiffViewer, files: any[]) {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  const askUser = (): Promise<void> => {
    return new Promise((resolve) => {
      rl.question(chalk.cyan('> '), (answer) => {
        const cmd = answer.toLowerCase().trim();

        if (cmd === 'n' || cmd === 'next') {
          if (viewer.nextFile()) {
            clearScreen();
            const output = viewer.renderAllDiffs(files);
            console.log(output);
            printNavigationHints(viewer);
            askUser().then(resolve);
          } else {
            console.log(chalk.yellow('Already on the last file'));
            askUser().then(resolve);
          }
        } else if (cmd === 'p' || cmd === 'prev') {
          if (viewer.prevFile()) {
            clearScreen();
            const output = viewer.renderAllDiffs(files);
            console.log(output);
            printNavigationHints(viewer);
            askUser().then(resolve);
          } else {
            console.log(chalk.yellow('Already on the first file'));
            askUser().then(resolve);
          }
        } else if (cmd === 'q' || cmd === 'quit') {
          rl.close();
          process.exit(0);
        } else if (cmd === 'h' || cmd === 'help') {
          printNavigationHints(viewer);
          askUser().then(resolve);
        } else {
          console.log(chalk.red('Unknown command. Type "h" for help'));
          askUser().then(resolve);
        }
      });
    });
  };

  await askUser();
}

function clearScreen() {
  console.clear();
  console.log(chalk.blue.bold('📊 Git Diff Viewer\n'));
}

function printNavigationHints(viewer: DiffViewer) {
  const current = viewer.getCurrentFileIndex() + 1;
  const total = viewer.getTotalFiles();
  console.log(chalk.dim(`File ${current}/${total} — Press 'n' for next, 'p' for prev, 'q' to quit, 'h' for help`));
}

function printSummary(diffResult: any) {
  const { files, summary } = diffResult;
  const insertionColor = chalk.green;
  const deletionColor = chalk.red;

  console.log(chalk.bold('Summary:'));
  console.log(
    `  ${chalk.cyan(summary.filesChanged)} file${summary.filesChanged !== 1 ? 's' : ''} changed, ` +
    `${insertionColor(`+${summary.insertions}`)} insertions, ` +
    `${deletionColor(`-${summary.deletions}`)} deletions`
  );

  console.log(chalk.bold('\nFiles:'));
  for (const file of files) {
    const addCount = file.hunks
      .flatMap((h: any) => h.lines)
      .filter((l: any) => l.type === 'add').length;
    const removeCount = file.hunks
      .flatMap((h: any) => h.lines)
      .filter((l: any) => l.type === 'remove').length;

    const statusEmoji: { [key: string]: string } = {
      added: '📄',
      modified: '✏️ ',
      deleted: '🗑️ ',
      renamed: '➜ ',
    };

    const stats = `${insertionColor(`+${addCount}`)} ${deletionColor(`-${removeCount}`)}`;
    console.log(`  ${statusEmoji[file.status]} ${file.newPath || file.oldPath} (${stats})`);
  }
}

function printHelp() {
  console.log(`
${chalk.bold('Git Diff Viewer')}

Usage: diffview [options] [ref]

Options:
  -s, --staged          Show staged changes (git diff --staged)
  --no-color            Disable colored output
  -h, --help            Show this help message

Examples:
  diffview              Show unstaged changes
  diffview --staged     Show staged changes
  diffview HEAD         Show changes against HEAD
  diffview main         Show changes against main branch
`);
}

main();
