import { execSync } from 'child_process';

export interface DiffHunk {
  oldStart: number;
  oldCount: number;
  newStart: number;
  newCount: number;
  lines: DiffLine[];
}

export interface DiffLine {
  type: 'add' | 'remove' | 'context';
  content: string;
  oldLineNum?: number;
  newLineNum?: number;
}

export interface FileDiff {
  oldPath: string;
  newPath: string;
  status: 'added' | 'modified' | 'deleted' | 'renamed';
  hunks: DiffHunk[];
}

export interface DiffResult {
  files: FileDiff[];
  summary: {
    filesChanged: number;
    insertions: number;
    deletions: number;
  };
}

export function getGitDiff(ref?: string): DiffResult {
  let diffCommand = 'git diff';
  
  if (ref) {
    diffCommand += ` ${ref}`;
  }

  let output: string;
  try {
    output = execSync(diffCommand, { encoding: 'utf-8' });
  } catch (error) {
    throw new Error(`Failed to execute git diff: ${error}`);
  }

  return parseDiff(output);
}

function parseDiff(rawDiff: string): DiffResult {
  const files: FileDiff[] = [];
  const lines = rawDiff.split('\n');
  
  let currentFile: Partial<FileDiff> | null = null;
  let currentHunk: Partial<DiffHunk> | null = null;
  let oldLineNum = 0;
  let newLineNum = 0;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];

    // File header
    if (line.startsWith('diff --git')) {
      if (currentFile && currentHunk) {
        (currentFile as FileDiff).hunks.push(currentHunk as DiffHunk);
      }
      if (currentFile) {
        files.push(currentFile as FileDiff);
      }

      const match = line.match(/diff --git a\/(.*) b\/(.*)/);
      if (match) {
        currentFile = {
          oldPath: match[1],
          newPath: match[2],
          hunks: [],
        };
        currentHunk = null;
      }
      continue;
    }

    // Status line
    if (line.startsWith('index ')) {
      continue;
    }

    if (line.startsWith('---')) {
      if (!currentFile) currentFile = { hunks: [], oldPath: '', newPath: '', status: 'modified' };
      const path = line.slice(4);
      if (path !== '/dev/null') {
        currentFile.oldPath = path;
      }
      continue;
    }

    if (line.startsWith('+++')) {
      if (!currentFile) currentFile = { hunks: [], oldPath: '', newPath: '', status: 'modified' };
      const path = line.slice(4);
      if (path !== '/dev/null') {
        currentFile.newPath = path;
      }

      // Determine status
      if (currentFile.oldPath === '/dev/null') {
        currentFile.status = 'added';
      } else if (currentFile.newPath === '/dev/null') {
        currentFile.status = 'deleted';
      } else {
        currentFile.status = 'modified';
      }
      continue;
    }

    // Hunk header
    if (line.startsWith('@@')) {
      if (currentHunk && currentFile && currentFile.hunks) {
        currentFile.hunks.push(currentHunk as DiffHunk);
      }

      const match = line.match(/@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@/);
      if (match) {
        oldLineNum = parseInt(match[1], 10);
        newLineNum = parseInt(match[3], 10);
        
        currentHunk = {
          oldStart: oldLineNum,
          oldCount: match[2] ? parseInt(match[2], 10) : 1,
          newStart: newLineNum,
          newCount: match[4] ? parseInt(match[4], 10) : 1,
          lines: [],
        };
      }
      continue;
    }

    // Diff content
    if (currentHunk) {
      if (line.startsWith('-')) {
        currentHunk.lines!.push({
          type: 'remove',
          content: line.slice(1),
          oldLineNum: oldLineNum++,
        });
      } else if (line.startsWith('+')) {
        currentHunk.lines!.push({
          type: 'add',
          content: line.slice(1),
          newLineNum: newLineNum++,
        });
      } else if (line.startsWith(' ')) {
        currentHunk.lines!.push({
          type: 'context',
          content: line.slice(1),
          oldLineNum: oldLineNum++,
          newLineNum: newLineNum++,
        });
      } else if (line.startsWith('\\')) {
        // "No newline at end of file" marker
        continue;
      }
    }
  }

  // Push the last hunk and file
  if (currentHunk && currentFile && currentFile.hunks) {
    currentFile.hunks.push(currentHunk as DiffHunk);
  }
  if (currentFile && currentFile.hunks) {
    files.push(currentFile as FileDiff);
  }

  // Calculate summary
  let insertions = 0;
  let deletions = 0;
  
  files.forEach((file) => {
    file.hunks.forEach((hunk) => {
      hunk.lines.forEach((line) => {
        if (line.type === 'add') insertions++;
        if (line.type === 'remove') deletions++;
      });
    });
  });

  return {
    files,
    summary: {
      filesChanged: files.length,
      insertions,
      deletions,
    },
  };
}
