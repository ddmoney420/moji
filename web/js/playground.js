/**
 * Playground - tab switching, input handling, live preview.
 */
const Playground = (() => {
  let currentTab = 'kaomoji';
  let debounceTimer = null;
  let lastOutput = '';
  const STORAGE_KEY = 'moji-playground';

  function init() {
    setupTabs();
    setupInputListeners();
    setupActions();
    setupToolbar();
    loadSavedState();

    MojiBridge.init().then(() => {
      hideSkeleton();
      populateSelects();
      render();
    }).catch(err => {
      hideSkeleton();
      showError('Failed to load WASM: ' + err.message);
    });
  }

  function setupTabs() {
    document.querySelectorAll('.playground-tab').forEach(tab => {
      tab.addEventListener('click', () => switchTab(tab.dataset.pg));
    });
  }

  function switchTab(pg) {
    currentTab = pg;
    document.querySelectorAll('.playground-tab').forEach(t => {
      t.classList.toggle('active', t.dataset.pg === pg);
    });
    document.querySelectorAll('.playground-panel').forEach(p => {
      p.classList.toggle('active', p.id === 'pg-' + pg);
    });
    document.getElementById('pg-cli-display').classList.remove('visible');
    render();
    saveState();
  }

  function setupInputListeners() {
    // Kaomoji
    on('kaomoji-search', 'input', () => debounceRender(100));
    on('kaomoji-category', 'change', render);
    document.getElementById('kaomoji-random').addEventListener('click', () => {
      if (!MojiBridge.isReady()) return;
      const result = JSON.parse(MojiBridge.kaomojiRandom());
      lastOutput = result.kaomoji || '';
      showOutput(lastOutput);
      updateCli();
    });

    // Banner
    on('banner-text', 'input', () => debounceRender(100));
    on('banner-font', 'change', render);

    // Effects
    on('effects-text', 'input', () => debounceRender(50));
    on('effects-type', 'change', render);

    // Filters
    on('filters-text', 'input', () => debounceRender(50));
    on('filters-type', 'change', render);

    // Gradient
    on('gradient-text', 'input', () => debounceRender(50));
    on('gradient-theme', 'change', render);
    on('gradient-mode', 'change', render);

    // Styles
    on('styles-text', 'input', () => debounceRender(50));
    on('styles-type', 'change', render);
    on('styles-border', 'change', render);
    on('styles-align', 'change', render);

    // QR
    on('qr-text', 'input', () => debounceRender(200));
    on('qr-charset', 'change', render);
    on('qr-invert', 'change', render);
    on('qr-compact-mode', 'change', render);

    // Patterns
    on('patterns-text', 'input', () => debounceRender(50));
    on('patterns-border', 'change', render);
    on('patterns-divider', 'change', render);
    on('patterns-preset', 'change', render);

    // Speech
    on('speech-text', 'input', () => debounceRender(50));
    on('speech-style', 'change', render);
    on('speech-width', 'input', () => debounceRender(100));
    on('speech-art', 'change', render);

    // Art
    on('art-search', 'input', () => debounceRender(100));
    on('art-category', 'change', render);
    document.getElementById('art-random').addEventListener('click', () => {
      if (!MojiBridge.isReady()) return;
      const result = JSON.parse(MojiBridge.artRandom());
      lastOutput = result.art || '';
      showOutput(result.name + ' (' + result.category + '):\n\n' + lastOutput);
      updateCli();
    });


    // Calendar
    on('calendar-view', 'change', render);
    on('calendar-year', 'input', () => debounceRender(200));
    on('calendar-month', 'change', render);
    on('calendar-monday', 'change', render);

    // Chain
    on('chain-text', 'input', () => debounceRender(50));
    on('chain-effect', 'change', render);
    on('chain-gradient', 'change', render);
    on('chain-border', 'change', render);
    on('chain-bubble', 'change', render);
  }

  function on(id, event, handler) {
    const el = document.getElementById(id);
    if (el) el.addEventListener(event, handler);
  }

  function debounceRender(ms) {
    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      render();
      saveState();
    }, ms);
  }

  function setupActions() {
    document.getElementById('pg-copy').addEventListener('click', () => {
      if (!lastOutput) return;
      navigator.clipboard.writeText(lastOutput).then(() => {
        const btn = document.getElementById('pg-copy');
        btn.textContent = 'Copied!';
        setTimeout(() => { btn.textContent = 'Copy Output'; }, 1500);
      });
    });

    document.getElementById('pg-save-image').addEventListener('click', saveImage);

    document.getElementById('pg-cli-toggle').addEventListener('click', () => {
      const display = document.getElementById('pg-cli-display');
      display.classList.toggle('visible');
      if (display.classList.contains('visible')) updateCli();
    });
  }

  function saveImage() {
    const outputEl = document.getElementById('playground-output');
    if (!outputEl.textContent.trim()) return;

    const fontSize = parseInt(document.getElementById('pg-fontsize').value) || 13;
    const lineHeight = Math.round(fontSize * 1.4);
    const padding = 20;
    const fontFamily = "'JetBrains Mono', monospace";

    // Extract styled lines from the output element
    const lines = extractStyledLines(outputEl);
    if (lines.length === 0) return;

    // Measure max width
    const measureCanvas = document.createElement('canvas');
    const measureCtx = measureCanvas.getContext('2d');
    measureCtx.font = fontSize + 'px ' + fontFamily;

    let maxWidth = 0;
    for (const line of lines) {
      let lineWidth = 0;
      for (const seg of line) {
        lineWidth += measureCtx.measureText(seg.text).width;
      }
      maxWidth = Math.max(maxWidth, lineWidth);
    }

    const canvasWidth = Math.ceil(maxWidth + padding * 2);
    const canvasHeight = Math.ceil(lines.length * lineHeight + padding * 2);

    const canvas = document.createElement('canvas');
    canvas.width = canvasWidth;
    canvas.height = canvasHeight;
    const ctx = canvas.getContext('2d');

    // Terminal background
    ctx.fillStyle = '#0d1117';
    ctx.fillRect(0, 0, canvasWidth, canvasHeight);

    // Render text
    ctx.font = fontSize + 'px ' + fontFamily;
    ctx.textBaseline = 'top';

    for (let i = 0; i < lines.length; i++) {
      let x = padding;
      const y = padding + i * lineHeight;
      for (const seg of lines[i]) {
        // Background color
        if (seg.bg) {
          const w = ctx.measureText(seg.text).width;
          ctx.fillStyle = seg.bg;
          ctx.fillRect(x, y, w, lineHeight);
        }
        // Text color
        ctx.fillStyle = seg.color || '#e6edf3';
        if (seg.bold) ctx.font = 'bold ' + fontSize + 'px ' + fontFamily;
        ctx.fillText(seg.text, x, y + (lineHeight - fontSize) / 2);
        if (seg.bold) ctx.font = fontSize + 'px ' + fontFamily;
        x += ctx.measureText(seg.text).width;
      }
    }

    // Download
    canvas.toBlob(blob => {
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'moji-' + currentTab + '.png';
      a.click();
      URL.revokeObjectURL(url);
    }, 'image/png');

    const btn = document.getElementById('pg-save-image');
    btn.textContent = 'Saved!';
    setTimeout(() => { btn.textContent = 'Save Image'; }, 1500);
  }

  function extractStyledLines(el) {
    // Walk DOM nodes and extract text segments with color info
    const segments = [];
    walkNodes(el, { color: null, bg: null, bold: false }, segments);

    // Split into lines
    const lines = [[]];
    for (const seg of segments) {
      const parts = seg.text.split('\n');
      for (let i = 0; i < parts.length; i++) {
        if (i > 0) lines.push([]);
        if (parts[i].length > 0) {
          lines[lines.length - 1].push({
            text: parts[i],
            color: seg.color,
            bg: seg.bg,
            bold: seg.bold
          });
        }
      }
    }
    return lines;
  }

  function walkNodes(node, style, segments) {
    if (node.nodeType === Node.TEXT_NODE) {
      if (node.textContent) {
        segments.push({
          text: node.textContent,
          color: style.color,
          bg: style.bg,
          bold: style.bold
        });
      }
      return;
    }
    if (node.nodeType !== Node.ELEMENT_NODE) return;

    const newStyle = { ...style };
    const cs = node.style;
    if (cs.color) newStyle.color = cs.color;
    if (cs.backgroundColor) newStyle.bg = cs.backgroundColor;
    if (cs.fontWeight === 'bold' || node.tagName === 'B') newStyle.bold = true;

    for (const child of node.childNodes) {
      walkNodes(child, newStyle, segments);
    }
  }

  function populateSelects() {
    // Banner fonts
    populateSelect('banner-font', MojiBridge.bannerFonts(), 'name', 'name');

    // Effects
    populateSelect('effects-type', MojiBridge.effectList(), 'name', 'name', 'desc');

    // Filters
    populateSelect('filters-type', MojiBridge.filterList(), 'name', 'name', 'desc');

    // Gradient themes
    populateSelect('gradient-theme', MojiBridge.gradientThemes(), 'name', 'name', 'desc');

    // Styles
    populateSelect('styles-type', MojiBridge.styleList(), 'name', 'name', 'desc');

    // Styles borders
    populateSelectWithEmpty('styles-border', MojiBridge.styleBorderList(), 'name', 'name', 'None');

    // Kaomoji categories
    try {
      const cats = JSON.parse(MojiBridge.kaomojiCategories());
      const select = document.getElementById('kaomoji-category');
      select.innerHTML = '<option value="">All</option>';
      cats.forEach(c => {
        const opt = document.createElement('option');
        opt.value = c;
        opt.textContent = c;
        select.appendChild(opt);
      });
    } catch (e) {}

    // Art categories
    try {
      const cats = JSON.parse(MojiBridge.artCategories());
      const select = document.getElementById('art-category');
      select.innerHTML = '<option value="">All</option>';
      cats.forEach(c => {
        const opt = document.createElement('option');
        opt.value = c;
        opt.textContent = c;
        select.appendChild(opt);
      });
    } catch (e) {}

    // Speech art options
    try {
      const arts = JSON.parse(MojiBridge.artList(''));
      const select = document.getElementById('speech-art');
      select.innerHTML = '<option value="">None</option>';
      arts.forEach(a => {
        const opt = document.createElement('option');
        opt.value = a.name;
        opt.textContent = a.name + ' (' + a.category + ')';
        select.appendChild(opt);
      });
    } catch (e) {}

    // Pattern presets
    try {
      const presets = JSON.parse(MojiBridge.patternListPatterns());
      const select = document.getElementById('patterns-preset');
      select.innerHTML = '<option value="">None</option>';
      presets.forEach(p => {
        const opt = document.createElement('option');
        opt.value = p;
        opt.textContent = p;
        select.appendChild(opt);
      });
    } catch (e) {}

    // Chain effect/gradient/border dropdowns
    populateSelectWithEmpty('chain-effect', MojiBridge.effectList(), 'name', 'name', 'None');
    populateSelectWithEmpty('chain-gradient', MojiBridge.gradientThemes(), 'name', 'name', 'None');

    try {
      const borders = JSON.parse(MojiBridge.patternListBorders());
      const select = document.getElementById('chain-border');
      select.innerHTML = '<option value="">None</option>';
      borders.forEach(b => {
        const opt = document.createElement('option');
        opt.value = b;
        opt.textContent = b;
        select.appendChild(opt);
      });
    } catch (e) {}
  }

  function populateSelect(id, jsonStr, valueKey, labelKey, descKey) {
    try {
      const items = JSON.parse(jsonStr);
      const select = document.getElementById(id);
      select.innerHTML = '';
      items.forEach(item => {
        const opt = document.createElement('option');
        opt.value = item[valueKey];
        opt.textContent = descKey ? item[labelKey] + ' - ' + item[descKey] : item[labelKey];
        select.appendChild(opt);
      });
    } catch (e) {}
  }

  function populateSelectWithEmpty(id, jsonStr, valueKey, labelKey, emptyLabel) {
    try {
      const items = JSON.parse(jsonStr);
      const select = document.getElementById(id);
      select.innerHTML = '<option value="">' + emptyLabel + '</option>';
      items.forEach(item => {
        const opt = document.createElement('option');
        opt.value = item[valueKey] || item;
        opt.textContent = item[labelKey] || item;
        select.appendChild(opt);
      });
    } catch (e) {}
  }

  function render() {
    if (!MojiBridge.isReady()) return;

    switch (currentTab) {
      case 'kaomoji': renderKaomoji(); break;
      case 'banner': renderBanner(); break;
      case 'effects': renderEffects(); break;
      case 'filters': renderFilters(); break;
      case 'gradient': renderGradient(); break;
      case 'styles': renderStyles(); break;
      case 'qr': renderQR(); break;
      case 'patterns': renderPatterns(); break;
      case 'speech': renderSpeech(); break;
      case 'art': renderArt(); break;

      case 'calendar': renderCalendar(); break;
      case 'chain': renderChain(); break;
    }
    updateCli();
  }

  function renderKaomoji() {
    const search = val('kaomoji-search');
    const category = val('kaomoji-category');
    const items = JSON.parse(MojiBridge.kaomojiList(search, category));
    const lines = items.slice(0, 50).map(i => i.name.padEnd(20) + i.kaomoji);
    lastOutput = lines.join('\n');
    showOutput(lastOutput);
  }

  function renderBanner() {
    const text = val('banner-text') || 'Hello';
    const font = val('banner-font') || 'standard';
    lastOutput = MojiBridge.banner(text, font);
    showOutput(lastOutput);
  }

  function renderEffects() {
    const text = val('effects-text') || 'Hello World';
    const effect = val('effects-type') || 'bubble';
    lastOutput = MojiBridge.effect(effect, text);
    showOutput(lastOutput);
  }

  function renderFilters() {
    const text = val('filters-text') || 'Hello World';
    const filter = val('filters-type') || 'rainbow';
    const raw = MojiBridge.filter(filter, text);
    lastOutput = raw;
    showAnsiOutput(raw);
  }

  function renderGradient() {
    const text = val('gradient-text') || 'Hello World';
    const theme = val('gradient-theme') || 'rainbow';
    const mode = val('gradient-mode') || 'horizontal';
    const raw = MojiBridge.gradientApply(text, theme, mode);
    lastOutput = raw;
    showAnsiOutput(raw);
  }

  function renderStyles() {
    const text = val('styles-text') || 'Hello World';
    const styleName = val('styles-type') || 'rainbow';
    const border = val('styles-border');
    const align = val('styles-align');

    let result = text;
    if (styleName && styleName !== 'none') {
      result = MojiBridge.style(result, styleName);
    }
    if (border) {
      result = MojiBridge.styleBorder(result, border);
    }
    if (align) {
      result = MojiBridge.styleAlign(result, align, 60);
    }
    lastOutput = result;
    showAnsiOutput(result);
  }

  function renderQR() {
    const text = val('qr-text') || 'https://github.com';
    const charset = val('qr-charset') || 'blocks';
    const invert = checked('qr-invert');
    const compact = checked('qr-compact-mode');

    let raw;
    if (compact) {
      raw = MojiBridge.qrCompact(text, invert);
    } else {
      raw = MojiBridge.qr(text, charset, invert);
    }
    lastOutput = raw;
    showAnsiOutput(raw);
  }

  function renderPatterns() {
    const text = val('patterns-text') || 'Hello World';
    const border = val('patterns-border') || 'single';
    const divider = val('patterns-divider');
    const preset = val('patterns-preset');

    let output = MojiBridge.pattern(text, border);
    if (divider) {
      output += '\n' + MojiBridge.divider(divider, 40);
    }
    if (preset) {
      output += '\n\n' + MojiBridge.patternPreset(preset, 40, 5);
    }
    lastOutput = output;
    showOutput(output);
  }

  function renderSpeech() {
    const text = val('speech-text') || 'Hello World!';
    const style = val('speech-style') || 'round';
    const width = parseInt(val('speech-width')) || 40;
    const artName = val('speech-art');

    let bubble = MojiBridge.speechWrap(text, style, width);
    if (artName) {
      const art = MojiBridge.artGet(artName);
      if (art) {
        bubble = MojiBridge.speechCombine(bubble, art);
      }
    }
    lastOutput = bubble;
    showOutput(bubble);
  }

  function renderArt() {
    const search = val('art-search');
    const category = val('art-category');

    if (search) {
      const results = JSON.parse(MojiBridge.artSearch(search));
      if (results.length > 0) {
        const art = MojiBridge.artGet(results[0].name);
        lastOutput = art;
        showOutput(results[0].name + ' (' + results[0].category + '):\n\n' + art);
      } else {
        lastOutput = '';
        showOutput('No art found for "' + search + '"');
      }
    } else {
      const items = JSON.parse(MojiBridge.artList(category));
      if (items.length > 0) {
        const art = MojiBridge.artGet(items[0].name);
        lastOutput = art;
        showOutput(items[0].name + ' (' + items[0].category + '):\n\n' + art);
      } else {
        lastOutput = '';
        showOutput('No art in this category');
      }
    }
  }


  function renderCalendar() {
    const view = val('calendar-view') || 'current';
    const year = parseInt(val('calendar-year')) || new Date().getFullYear();
    const month = parseInt(val('calendar-month')) || (new Date().getMonth() + 1);
    const monday = checked('calendar-monday');

    switch (view) {
      case 'current':
        lastOutput = MojiBridge.calendarCurrent();
        break;
      case 'month':
        lastOutput = MojiBridge.calendarMonth(year, month, monday);
        break;
      case 'year':
        lastOutput = MojiBridge.calendarYear(year);
        break;
      case 'week':
        lastOutput = MojiBridge.calendarWeek();
        break;
      case 'today':
        lastOutput = MojiBridge.calendarToday();
        break;
      case 'art':
        lastOutput = MojiBridge.calendarArt();
        break;
    }
    showAnsiOutput(lastOutput);
  }

  function renderChain() {
    const text = val('chain-text') || 'Hello World';
    const effect = val('chain-effect');
    const gradientTheme = val('chain-gradient');
    const border = val('chain-border');
    const bubble = val('chain-bubble');

    const opts = {};
    if (effect) opts.effect = effect;
    if (gradientTheme) { opts.gradient = gradientTheme; opts.gradientMode = 'horizontal'; }
    if (border) { opts.border = border; opts.borderPad = 1; }
    if (bubble) { opts.bubble = bubble; opts.bubbleWidth = 50; }

    const raw = MojiBridge.chainApply(text, JSON.stringify(opts));
    lastOutput = raw;
    // Chain may produce ANSI output if gradient is applied
    if (gradientTheme || raw.includes('\x1b[')) {
      showAnsiOutput(raw);
    } else {
      showOutput(raw);
    }
  }

  function showOutput(text) {
    document.getElementById('playground-output').innerHTML = escapeHtml(text);
    updateDims(text);
  }

  function showAnsiOutput(text) {
    document.getElementById('playground-output').innerHTML = AnsiToHtml.convert(text);
    updateDims(stripAnsi(text));
  }

  function hideSkeleton() {
    const el = document.getElementById('playground-skeleton');
    if (el) el.remove();
  }

  function setupToolbar() {
    const wrapToggle = document.getElementById('pg-wrap');
    const fontSlider = document.getElementById('pg-fontsize');
    const output = document.getElementById('playground-output');

    if (wrapToggle) {
      wrapToggle.addEventListener('change', () => {
        output.classList.toggle('wrap', wrapToggle.checked);
      });
    }

    if (fontSlider) {
      fontSlider.addEventListener('input', () => {
        const size = fontSlider.value;
        output.style.fontSize = size + 'px';
        document.getElementById('pg-fontsize-val').textContent = size;
      });
    }
  }

  function updateDims(text) {
    const el = document.getElementById('pg-dims');
    if (!el || !text) { if (el) el.textContent = '0 x 0'; return; }
    const lines = text.split('\n');
    const rows = lines.length;
    const cols = Math.max(...lines.map(l => l.length));
    el.textContent = cols + ' x ' + rows;
  }

  function stripAnsi(text) {
    return text.replace(/\x1b\[[0-9;]*m/g, '');
  }

  function showError(msg) {
    document.getElementById('playground-output').innerHTML =
      '<span style="color:#f55">' + escapeHtml(msg) + '</span>';
  }

  function updateCli() {
    const display = document.getElementById('pg-cli-display');
    if (!display.classList.contains('visible')) return;

    let cmd = 'moji ';
    switch (currentTab) {
      case 'kaomoji': {
        const search = val('kaomoji-search');
        cmd += search ? 'kaomoji ' + search : 'kaomoji --random';
        break;
      }
      case 'banner': {
        const text = val('banner-text') || 'Hello';
        const font = val('banner-font') || 'standard';
        cmd += 'banner "' + text + '" --font ' + font;
        break;
      }
      case 'effects': {
        const text = val('effects-text') || 'Hello World';
        const effect = val('effects-type') || 'bubble';
        cmd += 'effect ' + effect + ' "' + text + '"';
        break;
      }
      case 'filters': {
        const text = val('filters-text') || 'Hello World';
        const filter = val('filters-type') || 'rainbow';
        cmd += 'banner "' + text + '" --filter ' + filter;
        break;
      }
      case 'gradient': {
        const text = val('gradient-text') || 'Hello World';
        const theme = val('gradient-theme') || 'rainbow';
        const mode = val('gradient-mode') || 'horizontal';
        cmd += 'gradient "' + text + '" --theme ' + theme + ' --mode ' + mode;
        break;
      }
      case 'styles': {
        const text = val('styles-text') || 'Hello World';
        const style = val('styles-type') || 'rainbow';
        cmd += 'banner "' + text + '" --style ' + style;
        const border = val('styles-border');
        if (border) cmd += ' --border ' + border;
        break;
      }
      case 'qr': {
        const text = val('qr-text') || 'https://github.com';
        const charset = val('qr-charset') || 'blocks';
        cmd += 'qr "' + text + '" --charset ' + charset;
        if (checked('qr-invert')) cmd += ' --invert';
        if (checked('qr-compact-mode')) cmd += ' --compact';
        break;
      }
      case 'patterns': {
        const text = val('patterns-text') || 'Hello World';
        const border = val('patterns-border') || 'single';
        cmd += 'pattern border "' + text + '" --style ' + border;
        break;
      }
      case 'speech': {
        const text = val('speech-text') || 'Hello World!';
        const style = val('speech-style') || 'round';
        cmd += 'speech "' + text + '" --style ' + style;
        break;
      }
      case 'art': {
        const search = val('art-search');
        cmd += search ? 'art ' + search : 'art --random';
        break;
      }

      case 'calendar': {
        const view = val('calendar-view');
        cmd += 'calendar';
        if (view === 'year') cmd += ' --year ' + val('calendar-year');
        if (view === 'week') cmd += ' --week';
        break;
      }
      case 'chain': {
        const text = val('chain-text') || 'Hello World';
        cmd += 'banner "' + text + '"';
        const g = val('chain-gradient');
        const b = val('chain-border');
        const e = val('chain-effect');
        const bub = val('chain-bubble');
        if (g) cmd += ' --gradient ' + g;
        if (b) cmd += ' --border ' + b;
        if (e) cmd += ' --effect ' + e;
        if (bub) cmd += ' --bubble ' + bub;
        break;
      }
    }
    display.textContent = '$ ' + cmd;
  }

  function val(id) {
    const el = document.getElementById(id);
    return el ? el.value : '';
  }

  function checked(id) {
    const el = document.getElementById(id);
    return el ? el.checked : false;
  }

  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  function saveState() {
    try {
      const state = { tab: currentTab };
      const fields = [
        'banner-text', 'banner-font', 'effects-text', 'effects-type',
        'filters-text', 'filters-type', 'gradient-text', 'gradient-theme',
        'gradient-mode', 'qr-text', 'qr-charset', 'patterns-text',
        'patterns-border', 'styles-text', 'styles-type', 'styles-border',
        'styles-align', 'speech-text', 'speech-style', 'speech-width',
        'art-search', 'art-category',
        'calendar-view', 'calendar-year', 'calendar-month',
        'chain-text', 'chain-effect', 'chain-gradient', 'chain-border', 'chain-bubble'
      ];
      fields.forEach(id => { state[id] = val(id); });
      localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
    } catch (e) {}
  }

  function loadSavedState() {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (!raw) return;
      const state = JSON.parse(raw);
      if (state.tab) switchTab(state.tab);
      Object.keys(state).forEach(id => {
        if (id === 'tab') return;
        const el = document.getElementById(id);
        if (el && state[id]) el.value = state[id];
      });
    } catch (e) {}
  }

  return { init };
})();

document.addEventListener('DOMContentLoaded', Playground.init);
