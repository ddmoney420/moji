/**
 * WASM bridge - loads the moji WASM binary and exposes a Promise-based API.
 */
const MojiBridge = (() => {
  let ready = false;
  let readyPromise = null;
  let loadError = null;

  function init() {
    if (readyPromise) return readyPromise;

    readyPromise = new Promise((resolve, reject) => {
      if (typeof Go === 'undefined') {
        reject(new Error('wasm_exec.js not loaded'));
        return;
      }

      const go = new Go();
      const wasmPath = 'wasm/moji.wasm';

      if (typeof WebAssembly.instantiateStreaming === 'function') {
        WebAssembly.instantiateStreaming(fetch(wasmPath), go.importObject)
          .then(result => {
            go.run(result.instance);
            ready = true;
            resolve();
          })
          .catch(() => {
            fetchAndInstantiate(go, wasmPath, resolve, reject);
          });
      } else {
        fetchAndInstantiate(go, wasmPath, resolve, reject);
      }
    });

    readyPromise.catch(err => {
      loadError = err;
    });

    return readyPromise;
  }

  function fetchAndInstantiate(go, wasmPath, resolve, reject) {
    fetch(wasmPath)
      .then(response => response.arrayBuffer())
      .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
      .then(result => {
        go.run(result.instance);
        ready = true;
        resolve();
      })
      .catch(reject);
  }

  function isReady() { return ready; }
  function getError() { return loadError; }

  // --- Banner ---
  function banner(text, font) {
    if (!ready) return 'Loading WASM...';
    return window.mojiBanner(text, font);
  }
  function bannerFonts() {
    if (!ready) return '[]';
    return window.mojiBannerFonts();
  }

  // --- Kaomoji ---
  function kaomojiGet(name) {
    if (!ready) return '';
    return window.mojiKaomoji(name);
  }
  function kaomojiList(search, category) {
    if (!ready) return '[]';
    return window.mojiKaomojiList(search || '', category || '');
  }
  function kaomojiRandom() {
    if (!ready) return '{}';
    return window.mojiKaomojiRandom();
  }
  function kaomojiCategories() {
    if (!ready) return '[]';
    return window.mojiKaomojiCategories();
  }
  function kaomojiArt(name) {
    if (!ready) return '';
    return window.mojiKaomojiArt(name);
  }
  function kaomojiArtList(category) {
    if (!ready) return '[]';
    return window.mojiKaomojiArtList(category || '');
  }
  function kaomojiSuggest(name) {
    if (!ready) return '[]';
    return window.mojiKaomojiSuggest(name);
  }
  function smileyToEmoji(smiley) {
    if (!ready) return '';
    return window.mojiSmileyToEmoji(smiley);
  }

  // --- Effects ---
  function effect(effectName, text) {
    if (!ready) return '';
    return window.mojiEffect(effectName, text);
  }
  function effectList() {
    if (!ready) return '[]';
    return window.mojiEffectList();
  }

  // --- Filters ---
  function filter(name, text) {
    if (!ready) return '';
    return window.mojiFilter(name, text);
  }
  function filterList() {
    if (!ready) return '[]';
    return window.mojiFilterList();
  }
  function filterChain(text, chainSpec) {
    if (!ready) return '';
    return window.mojiFilterChain(text, chainSpec);
  }

  // --- Gradient ---
  function gradientApply(text, theme, mode) {
    if (!ready) return '';
    return window.mojiGradient(text, theme, mode);
  }
  function gradientThemes() {
    if (!ready) return '[]';
    return window.mojiGradientThemes();
  }

  // --- QR ---
  function qr(text, charset, invert) {
    if (!ready) return '';
    return window.mojiQR(text, charset, !!invert);
  }
  function qrCompact(text, invert) {
    if (!ready) return '';
    return window.mojiQRCompact(text, !!invert);
  }

  // --- Patterns ---
  function pattern(text, style, padding) {
    if (!ready) return '';
    return window.mojiPattern(text, style, padding || 1);
  }
  function divider(style, width) {
    if (!ready) return '';
    return window.mojiDivider(style, width);
  }
  function patternCreate(name, width, height) {
    if (!ready) return '';
    return window.mojiPatternCreate(name, width, height);
  }
  function patternPreset(name, width, height) {
    if (!ready) return '';
    return window.mojiPatternPreset(name, width, height);
  }
  function patternSymmetric(half) {
    if (!ready) return '';
    return window.mojiPatternSymmetric(half);
  }
  function patternListBorders() {
    if (!ready) return '[]';
    return window.mojiPatternListBorders();
  }
  function patternListDividers() {
    if (!ready) return '[]';
    return window.mojiPatternListDividers();
  }
  function patternListPatterns() {
    if (!ready) return '[]';
    return window.mojiPatternListPatterns();
  }


  // --- Styles ---
  function style(text, styleName) {
    if (!ready) return '';
    return window.mojiStyle(text, styleName);
  }
  function styleList() {
    if (!ready) return '[]';
    return window.mojiStyleList();
  }
  function styleBorder(text, borderName) {
    if (!ready) return '';
    return window.mojiStyleBorder(text, borderName);
  }
  function styleBorderList() {
    if (!ready) return '[]';
    return window.mojiStyleBorderList();
  }
  function styleAlign(text, align, width) {
    if (!ready) return '';
    return window.mojiStyleAlign(text, align, width);
  }

  // --- Speech ---
  function speechWrap(text, styleName, maxWidth) {
    if (!ready) return '';
    return window.mojiSpeech(text, styleName, maxWidth);
  }
  function speechStyles() {
    if (!ready) return '[]';
    return window.mojiSpeechStyles();
  }
  function speechCombine(bubble, art) {
    if (!ready) return '';
    return window.mojiSpeechCombine(bubble, art);
  }

  // --- Art DB ---
  function artList(category) {
    if (!ready) return '[]';
    return window.mojiArtList(category || '');
  }
  function artGet(name) {
    if (!ready) return '';
    return window.mojiArtGet(name);
  }
  function artSearch(query) {
    if (!ready) return '[]';
    return window.mojiArtSearch(query);
  }
  function artCategories() {
    if (!ready) return '[]';
    return window.mojiArtCategories();
  }
  function artRandom() {
    if (!ready) return '{}';
    return window.mojiArtRandom();
  }

  // --- Chain ---
  function chainApply(text, optsJSON) {
    if (!ready) return '';
    return window.mojiChain(text, optsJSON);
  }

  // --- Calendar ---
  function calendarMonth(year, month, mondayFirst) {
    if (!ready) return '';
    return window.mojiCalendarMonth(year, month, !!mondayFirst);
  }
  function calendarYear(year) {
    if (!ready) return '';
    return window.mojiCalendarYear(year);
  }
  function calendarCurrent() {
    if (!ready) return '';
    return window.mojiCalendarCurrent();
  }
  function calendarToday() {
    if (!ready) return '';
    return window.mojiCalendarToday();
  }
  function calendarArt() {
    if (!ready) return '';
    return window.mojiCalendarArt();
  }
  function calendarWeek() {
    if (!ready) return '';
    return window.mojiCalendarWeek();
  }

  return {
    init, isReady, getError,
    // Banner
    banner, bannerFonts,
    // Kaomoji
    kaomojiGet, kaomojiList, kaomojiRandom, kaomojiCategories,
    kaomojiArt, kaomojiArtList, kaomojiSuggest, smileyToEmoji,
    // Effects
    effect, effectList,
    // Filters
    filter, filterList, filterChain,
    // Gradient
    gradientApply, gradientThemes,
    // QR
    qr, qrCompact,
    // Patterns
    pattern, divider, patternCreate, patternPreset, patternSymmetric,
    patternListBorders, patternListDividers, patternListPatterns,

    // Styles
    style, styleList, styleBorder, styleBorderList, styleAlign,
    // Speech
    speechWrap, speechStyles, speechCombine,
    // Art DB
    artList, artGet, artSearch, artCategories, artRandom,
    // Chain
    chainApply,
    // Calendar
    calendarMonth, calendarYear, calendarCurrent, calendarToday, calendarArt, calendarWeek,
  };
})();
