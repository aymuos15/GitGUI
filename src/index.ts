#!/usr/bin/env node

import chalk from 'chalk';
import { getGitDiff } from './parser.js';
import { DiffTUI } from './tui.js';

interface Options {
  ref?: string;
  staged?: boolean;
}

async function main() {
  const args = process.argv.slice(2);
  const options: Options = {};

  // Parse command line arguments
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '--staged' || arg === '-s') {
      options.staged = true;
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

    const diffResult = getGitDiff(ref);

    if (diffResult.files.length === 0) {
      console.log(chalk.yellow('No changes to display'));
      process.exit(0);
    }

    // Launch interactive TUI
    const tui = new DiffTUI(diffResult);
    tui.show();
  } catch (error) {
    if (error instanceof Error) {
      console.error(chalk.red(`Error: ${error.message}`));
    } else {
      console.error(chalk.red('An unknown error occurred'));
    }
    process.exit(1);
  }
}

function printHelp() {
  console.log(`
${chalk.bold('Git Diff Viewer')}

Usage: diffview [options] [ref]

Options:
  -s, --staged   Show staged changes (git diff --staged)
  -h, --help     Show this help message

Examples:
  diffview              Show unstaged changes
  diffview --staged     Show staged changes
  diffview HEAD         Show changes against HEAD
  diffview main         Show changes against main branch

Interactive controls:
  n              Next file
  p              Previous file
  j/k or ↓/↑    Scroll down/up
  Page Down/Up   Scroll page
  ?/h            Show help
  q              Quit
`);
}

main();
