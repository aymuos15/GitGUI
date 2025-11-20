import { createSignal, For, Show, createEffect } from 'solid-js';
import { render, Text, Box, useKeyboard } from '@opentui/solid';
import { DiffResult, FileDiff } from './parser.js';

export function createApp(diffResult: DiffResult) {
  const [currentFileIndex, setCurrentFileIndex] = createSignal(0);
  const [scrollOffset, setScrollOffset] = createSignal(0);

  const currentFile = () => diffResult.files[currentFileIndex()];

  return () => {
    useKeyboard((key) => {
      // Navigate between files
      if (key === 'n' && currentFileIndex() < diffResult.files.length - 1) {
        setCurrentFileIndex(currentFileIndex() + 1);
        setScrollOffset(0);
      } else if (key === 'p' && currentFileIndex() > 0) {
        setCurrentFileIndex(currentFileIndex() - 1);
        setScrollOffset(0);
      }
      // Scroll
      else if (key === 'j' || key === 'down') {
        setScrollOffset(scrollOffset() + 1);
      } else if (key === 'k' || key === 'up') {
        setScrollOffset(Math.max(0, scrollOffset() - 1));
      }
      // Quit
      else if (key === 'q' || key === 'escape') {
        process.exit(0);
      }
    });

    return Box({
      flexDirection: 'column',
      height: '100%',
      children: [
        // Header
        Box({
          flexDirection: 'column',
          paddingLeft: 2,
          paddingRight: 2,
          paddingTop: 1,
          borderStyle: 'single',
          borderBottom: true,
          children: [
            Text({ bold: true, children: 'Git Diff Viewer' }),
            Text({
              children: [
                `${diffResult.summary.filesChanged} file${diffResult.summary.filesChanged !== 1 ? 's' : ''} changed  `,
                Text({ color: 'green', children: `+${diffResult.summary.insertions}` }),
                Text({ color: 'red', children: ` -${diffResult.summary.deletions}` }),
              ],
            }),
          ],
        }),

        // Main area
        Box({
          flexGrow: 1,
          flexDirection: 'row',
          children: [
            // Sidebar
            Box({
              width: 25,
              flexDirection: 'column',
              paddingLeft: 1,
              paddingRight: 1,
              borderStyle: 'single',
              borderRight: true,
              children: [
                Text({ bold: true, children: 'Files' }),
                For({
                  each: diffResult.files,
                  children: (file: FileDiff, index: () => number) => {
                    const isActive = () => index() === currentFileIndex();
                    const fileName = (file.newPath || file.oldPath).split('/').pop() || file.newPath || file.oldPath;
                    
                    return Text({
                      inverse: isActive(),
                      children: `${file.status[0]} ${fileName.substring(0, 12)}`,
                    });
                  },
                }),
              ],
            }),

            // Content
            Box({
              flexGrow: 1,
              flexDirection: 'column',
              paddingLeft: 1,
              children: [
                Text({ bold: true, children: currentFile() ? `${currentFile().status.toUpperCase()}` : 'No file' }),
                Text({ children: currentFile() ? (currentFile().newPath || currentFile().oldPath) : '' }),
              ],
            }),
          ],
        }),

        // Footer
        Box({
          paddingLeft: 2,
          backgroundColor: 'black',
          children: [
            Text({ bold: true, color: 'white', children: '[n]' }),
            Text({ color: 'white', children: 'ext  ' }),
            Text({ bold: true, color: 'white', children: '[p]' }),
            Text({ color: 'white', children: 'rev  ' }),
            Text({ bold: true, color: 'white', children: '[q]' }),
            Text({ color: 'white', children: ' quit' }),
          ],
        }),
      ],
    });
  };
}
