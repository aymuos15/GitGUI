#!/usr/bin/env node

import chalk from 'chalk';
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
    const viewer = new DiffViewer();
    const output = viewer.renderAllDiffs(diffResult.files);
    console.log(output);
  } catch (error) {
    if (error instanceof Error) {
      console.error(chalk.red(`Error: ${error.message}`));
    } else {
      console.error(chalk.red('An unknown error occurred'));
    }
    process.exit(1);
  }
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
