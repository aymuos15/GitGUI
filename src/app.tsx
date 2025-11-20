import { Component, createSignal, For, Show } from 'solid-js';
import { useKeyboard } from '@opentui/solid';
import { TextAttributes } from '@opentui/core';
import { DiffResult, FileDiff } from './parser.js';

interface AppProps {
  diffResult: DiffResult;
  onExit: () => void;
}

export const App: Component<AppProps> = (props) => {
  const [currentFileIndex, setCurrentFileIndex] = createSignal(0);
  
  const currentFile = () => props.diffResult.files[currentFileIndex()];

  useKeyboard((key) => {
    if (key === 'n' && currentFileIndex() < props.diffResult.files.length - 1) {
      setCurrentFileIndex(currentFileIndex() + 1);
    } else if (key === 'p' && currentFileIndex() > 0) {
      setCurrentFileIndex(currentFileIndex() - 1);
    } else if (key === 'q') {
      props.onExit();
    }
  });

  const summary = () => props.diffResult.summary;
  const path = () => currentFile() ? (currentFile().newPath || currentFile().oldPath) : 'N/A';

  return (
    <box flex-direction="column" width="100%" height="100%">
      {/* Header */}
      <box flex-direction="column" padding-left={2} padding-right={2} padding-top={1} border-bottom>
        <text attributes={TextAttributes.BOLD}>Git Diff Viewer</text>
        <text>
          {summary().filesChanged} file{summary().filesChanged !== 1 ? 's' : ''} changed  
          <text attributes={TextAttributes.GREEN}> +{summary().insertions}</text>
          <text attributes={TextAttributes.RED}> -{summary().deletions}</text>
        </text>
        <text>
          <text attributes={TextAttributes.BOLD}>Current:</text> {path()}
        </text>
      </box>

      {/* Main content area */}
      <box flex-direction="row" flex-grow={1}>
        {/* Sidebar */}
        <box width={25} flex-direction="column" padding-left={1} border-right>
          <text attributes={TextAttributes.BOLD}>Files</text>
          <For each={props.diffResult.files}>
            {(file, index) => {
              const isActive = () => index() === currentFileIndex();
              const fileName = (file.newPath || file.oldPath).split('/').pop() || file.newPath || file.oldPath;
              const displayName = fileName.length > 15 ? fileName.substring(0, 14) + '…' : fileName;
              
              const statusIcon = {
                added: '✚',
                modified: '◆',
                deleted: '✕',
                renamed: '➜',
              }[file.status] || '◆';

              return (
                <text attributes={isActive() ? TextAttributes.INVERSE : 0}>
                  {statusIcon} {displayName}
                </text>
              );
            }}
          </For>
        </box>

        {/* Content */}
        <box flex-grow={1} flex-direction="column" padding-left={2} padding-right={2}>
          <Show when={currentFile()} fallback={<text>No file selected</text>}>
            {(file) => (
              <>
                <text attributes={TextAttributes.BOLD}>
                  {
                    {
                      added: '✚ ADDED',
                      modified: '◆ MODIFIED',
                      deleted: '✕ DELETED',
                      renamed: '➜ RENAMED',
                    }[file().status] || '◆ CHANGED'
                  }
                </text>
                <text>{file().newPath || file().oldPath}</text>
                <text> </text>

                <For each={file().hunks}>
                  {(hunk) => (
                    <>
                      <text attributes={TextAttributes.BLUE}>
                        @@ -{hunk.oldStart},{hunk.oldCount} +{hunk.newStart},{hunk.newCount} @@
                      </text>
                      <For each={hunk.lines}>
                        {(line) => {
                          const lineNum = line.oldLineNum || line.newLineNum || 0;
                          const marker = line.type === 'remove' ? '−' : line.type === 'add' ? '+' : ' ';
                          const attrs = line.type === 'remove' ? TextAttributes.RED : 
                                       line.type === 'add' ? TextAttributes.GREEN : 0;
                          
                          return (
                            <text attributes={attrs}>
                              {String(lineNum).padStart(4)} {marker} {line.content}
                            </text>
                          );
                        }}
                      </For>
                    </>
                  )}
                </For>
              </>
            )}
          </Show>
        </box>
      </box>

      {/* Footer */}
      <box padding-left={2} padding-right={2} background-color="black">
        <text attributes={TextAttributes.BOLD | TextAttributes.WHITE}>[n]</text>
        <text attributes={TextAttributes.WHITE}>ext  </text>
        <text attributes={TextAttributes.BOLD | TextAttributes.WHITE}>[p]</text>
        <text attributes={TextAttributes.WHITE}>rev  </text>
        <text attributes={TextAttributes.BOLD | TextAttributes.WHITE}>[q]</text>
        <text attributes={TextAttributes.WHITE}> quit</text>
        <text attributes={TextAttributes.WHITE}> | File {currentFileIndex() + 1}/{props.diffResult.files.length}</text>
      </box>
    </box>
  );
};
