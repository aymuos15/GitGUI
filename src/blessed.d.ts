declare module 'blessed' {
  interface Screen {
    key(keys: string[], callback: () => void): void;
    once(event: string, callback: () => void): void;
    render(): void;
    focus(): void;
  }

  interface BoxElement {
    setContent(content: string): void;
    setScrollPerc(perc: number): void;
    scroll(offset: number): void;
    destroy(): void;
  }

  interface ScreenOptions {
    mouse?: boolean;
    keyboard?: boolean;
    title?: string;
    style?: any;
  }

  interface BoxOptions {
    parent?: any;
    top?: number | string;
    left?: number | string;
    right?: number;
    bottom?: number;
    height?: number | string;
    width?: number | string;
    scrollable?: boolean;
    mouse?: boolean;
    keys?: boolean;
    vi?: boolean;
    alwaysScroll?: boolean;
    border?: string;
    style?: any;
    padding?: any;
  }

  function screen(options: ScreenOptions): Screen;
  function box(options: BoxOptions): BoxElement;

  export default {
    screen,
    box,
  };

  export namespace Widgets {
    type Screen = import('./blessed.d.ts').Screen;
    type BoxElement = import('./blessed.d.ts').BoxElement;
  }
}
