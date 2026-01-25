/**
 * Landing page animations - hero kaomoji cycling, install tabs, copy buttons.
 */
(function() {
  const kaomojis = [
    '(\\u30CE\\u25D5\\u30EE\\u25D5)\\u30CE*:\\uFF65\\uFF9F\\u2727',
    '(\\u2565\\u203F\\u2565)',
    '\\u0CA5_\\u0CA5',
    '(\\u256F\\u00B0\\u25A1\\u00B0)\\u256F\\uFE35 \\u253B\\u2501\\u253B',
    '(\\u3065\\uFFE3 \\u00B3\\uFFE3)\\u3065',
    '\\u2764\\uFE0F (\\u0E51\\u2022\\u0300\\u2323\\u2022\\u0301\\u0E51)',
    '\\u10DA(\\u30FB\\u2200\\u30FB\\u10DA)',
    '(\\u3000\\u30FB\\u03C9\\u30FB)\\u3064\\u2661',
  ];

  // Decode unicode escapes for display
  function decodeKaomoji(str) {
    return str.replace(/\\u([0-9a-fA-F]{4})/g, (_, hex) =>
      String.fromCharCode(parseInt(hex, 16))
    );
  }

  let kaomojiIndex = 0;
  const heroKaomoji = document.getElementById('hero-kaomoji');

  function cycleKaomoji() {
    if (!heroKaomoji) return;
    heroKaomoji.style.opacity = '0';
    setTimeout(() => {
      heroKaomoji.textContent = decodeKaomoji(kaomojis[kaomojiIndex]);
      heroKaomoji.style.opacity = '1';
      kaomojiIndex = (kaomojiIndex + 1) % kaomojis.length;
    }, 300);
  }

  // Start cycling
  if (heroKaomoji) {
    heroKaomoji.style.transition = 'opacity 0.3s';
    cycleKaomoji();
    setInterval(cycleKaomoji, 2500);
  }

  // Install tabs
  document.querySelectorAll('.install-tab').forEach(tab => {
    tab.addEventListener('click', () => {
      const target = tab.dataset.tab;

      document.querySelectorAll('.install-tab').forEach(t =>
        t.classList.toggle('active', t === tab)
      );

      document.querySelectorAll('.install-content').forEach(c => {
        c.classList.toggle('active', c.id === 'install-' + target);
      });
    });
  });

  // Copy buttons
  document.querySelectorAll('.copy-btn').forEach(btn => {
    btn.addEventListener('click', () => {
      const text = btn.dataset.copy;
      navigator.clipboard.writeText(text).then(() => {
        btn.textContent = 'Copied!';
        btn.classList.add('copied');
        setTimeout(() => {
          btn.textContent = 'Copy';
          btn.classList.remove('copied');
        }, 1500);
      });
    });
  });

  // Smooth scroll for nav links
  document.querySelectorAll('.nav-links a[href^="#"]').forEach(link => {
    link.addEventListener('click', e => {
      e.preventDefault();
      const target = document.querySelector(link.getAttribute('href'));
      if (target) {
        target.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });
  });

  // Hero ASCII art glow pulse
  const heroAscii = document.getElementById('hero-ascii');
  if (heroAscii) {
    let glowIntensity = 10;
    let glowDirection = 1;

    function pulseGlow() {
      glowIntensity += glowDirection * 0.5;
      if (glowIntensity >= 20) glowDirection = -1;
      if (glowIntensity <= 5) glowDirection = 1;
      heroAscii.style.textShadow = `0 0 ${glowIntensity}px rgba(0, 255, 65, 0.5)`;
      requestAnimationFrame(pulseGlow);
    }

    pulseGlow();
  }
})();
