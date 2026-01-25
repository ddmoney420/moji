/**
 * ANSI escape code to HTML converter.
 * Handles the subset of ANSI codes that moji produces:
 * - 24-bit foreground/background colors
 * - Basic attributes (bold, italic, underline, inverse)
 * - Combined codes (e.g. \033[1;33m)
 * - Standard foreground/background colors
 * - Reset sequences
 */
const AnsiToHtml = (() => {
  function convert(text) {
    if (!text || !text.includes('\x1b')) {
      return escapeHtml(text || '');
    }

    let result = '';
    let i = 0;
    let openSpans = 0;

    while (i < text.length) {
      if (text[i] === '\x1b' && text[i + 1] === '[') {
        let j = i + 2;
        while (j < text.length && text[j] !== 'm') {
          j++;
        }
        if (j < text.length) {
          const code = text.substring(i + 2, j);
          const parsed = parseAnsiCode(code);

          if (parsed.reset) {
            for (let k = 0; k < openSpans; k++) {
              result += '</span>';
            }
            openSpans = 0;
          }

          if (parsed.styles.length > 0) {
            result += '<span style="' + parsed.styles.join(';') + '">';
            openSpans++;
          }

          i = j + 1;
        } else {
          result += escapeHtml(text[i]);
          i++;
        }
      } else {
        result += escapeHtml(text[i]);
        i++;
      }
    }

    for (let k = 0; k < openSpans; k++) {
      result += '</span>';
    }

    return result;
  }

  function parseAnsiCode(code) {
    const parts = code.split(';');
    const result = { reset: false, styles: [] };

    let i = 0;
    while (i < parts.length) {
      const n = parseInt(parts[i], 10);

      if (isNaN(n) || n === 0) {
        result.reset = true;
        i++;
      } else if (n === 1) {
        result.styles.push('font-weight:bold');
        i++;
      } else if (n === 2) {
        result.styles.push('opacity:0.6');
        i++;
      } else if (n === 3) {
        result.styles.push('font-style:italic');
        i++;
      } else if (n === 4) {
        result.styles.push('text-decoration:underline');
        i++;
      } else if (n === 5) {
        result.styles.push('animation:ansiBlink 1s step-end infinite');
        i++;
      } else if (n === 7) {
        // Inverse - swap fg/bg
        result.styles.push('background-color:#e6edf3');
        result.styles.push('color:#000');
        result.styles.push('padding:0 2px');
        i++;
      } else if (n === 9) {
        result.styles.push('text-decoration:line-through');
        i++;
      } else if (n >= 30 && n <= 37) {
        const colors = ['#000', '#c00', '#0a0', '#ca0', '#00c', '#c0c', '#0cc', '#ccc'];
        result.styles.push('color:' + colors[n - 30]);
        i++;
      } else if (n >= 40 && n <= 47) {
        const colors = ['#000', '#c00', '#0a0', '#ca0', '#00c', '#c0c', '#0cc', '#fff'];
        result.styles.push('background-color:' + colors[n - 40]);
        i++;
      } else if (n >= 90 && n <= 97) {
        const colors = ['#555', '#f55', '#5f5', '#ff5', '#55f', '#f5f', '#5ff', '#fff'];
        result.styles.push('color:' + colors[n - 90]);
        i++;
      } else if (n >= 100 && n <= 107) {
        const colors = ['#555', '#f55', '#5f5', '#ff5', '#55f', '#f5f', '#5ff', '#fff'];
        result.styles.push('background-color:' + colors[n - 100]);
        i++;
      } else if (n === 38 && parts[i + 1] === '2') {
        const r = parts[i + 2] || 0;
        const g = parts[i + 3] || 0;
        const b = parts[i + 4] || 0;
        result.styles.push('color:rgb(' + r + ',' + g + ',' + b + ')');
        i += 5;
      } else if (n === 48 && parts[i + 1] === '2') {
        const r = parts[i + 2] || 0;
        const g = parts[i + 3] || 0;
        const b = parts[i + 4] || 0;
        result.styles.push('background-color:rgb(' + r + ',' + g + ',' + b + ')');
        i += 5;
      } else {
        i++;
      }
    }

    return result;
  }

  function escapeHtml(text) {
    return text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;');
  }

  return { convert };
})();
