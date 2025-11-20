/**
 * Convert blessed-style markup tags to ANSI escape codes
 */
export function parseMarkup(text: string): string {
  // Map of tag names to ANSI codes
  const colorMap: { [key: string]: string } = {
    // Foreground colors
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m',
    cyan: '\x1b[36m',
    white: '\x1b[37m',
    black: '\x1b[30m',
    gray: '\x1b[90m',

    // Background colors
    'bg-red': '\x1b[41m',
    'bg-green': '\x1b[42m',
    'bg-yellow': '\x1b[43m',
    'bg-blue': '\x1b[44m',
    'bg-magenta': '\x1b[45m',
    'bg-cyan': '\x1b[46m',
    'bg-white': '\x1b[47m',
    'bg-black': '\x1b[40m',

    // Styles
    bold: '\x1b[1m',
    dim: '\x1b[2m',
    italic: '\x1b[3m',
    underline: '\x1b[4m',
    blink: '\x1b[5m',
    inverse: '\x1b[7m',
  };

  const reset = '\x1b[0m';

  // Replace tags with ANSI codes
  let result = text;

  // Replace opening tags
  for (const [tag, code] of Object.entries(colorMap)) {
    const pattern = new RegExp(`\\{${tag}\\}`, 'g');
    result = result.replace(pattern, code);
  }

  // Replace closing tags
  result = result.replace(/\{\/[^}]+\}/g, reset);

  return result;
}

/**
 * Strip all markup tags from text (for length calculations)
 */
export function stripMarkup(text: string): string {
  return text.replace(/\{[^}]+\}/g, '');
}
