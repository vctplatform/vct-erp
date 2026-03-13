/* ═══════════════════════════════════════════
   VCT Platform — Shared JavaScript
   ═══════════════════════════════════════════ */

// ══ THEME TOGGLE ══
function toggleTheme() {
  const html = document.documentElement;
  const next = html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
  html.setAttribute('data-theme', next);
  const btn = document.getElementById('themeBtn');
  if (btn) btn.innerHTML = next === 'dark' ? '🌙' : '☀️';
  localStorage.setItem('vct-theme', next);
  // Re-init mermaid if present
  if (typeof initMermaid === 'function') initMermaid(next);
}

// Restore theme on load
(function () {
  const saved = localStorage.getItem('vct-theme');
  if (saved) {
    document.documentElement.setAttribute('data-theme', saved);
  }
  document.addEventListener('DOMContentLoaded', function () {
    const t = document.documentElement.getAttribute('data-theme');
    const btn = document.getElementById('themeBtn');
    if (btn) btn.innerHTML = t === 'dark' ? '🌙' : '☀️';
  });
})();

// ══ NAVIGATION ══
document.addEventListener('DOMContentLoaded', function () {
  // Sticky nav shadow on scroll
  const nav = document.querySelector('.site-nav');
  if (nav) {
    window.addEventListener('scroll', function () {
      nav.classList.toggle('scrolled', window.scrollY > 20);
    }, { passive: true });
  }

  // Hamburger menu toggle
  const hamburger = document.querySelector('.nav-hamburger');
  const navLinks = document.querySelector('.nav-links');
  if (hamburger && navLinks) {
    hamburger.addEventListener('click', function () {
      navLinks.classList.toggle('open');
      hamburger.innerHTML = navLinks.classList.contains('open') ? '✕' : '☰';
    });
    // Close menu on link click
    navLinks.querySelectorAll('a').forEach(function (a) {
      a.addEventListener('click', function () {
        navLinks.classList.remove('open');
        hamburger.innerHTML = '☰';
      });
    });
  }

  // Mark current page as active in nav
  const currentPage = window.location.pathname.split('/').pop() || 'index.html';
  document.querySelectorAll('.nav-links a').forEach(function (a) {
    const href = a.getAttribute('href');
    if (href === currentPage || (currentPage === '' && href === 'index.html')) {
      a.classList.add('active');
    }
  });

  // ══ TABS (for pitch.html) ══
  const tabBtns = document.querySelectorAll('.tab-btn');
  if (tabBtns.length > 0) {
    tabBtns.forEach(function (btn) {
      btn.addEventListener('click', function () {
        const tab = this.dataset.tab;
        document.querySelectorAll('.tab-btn').forEach(function (b) { b.classList.remove('active'); });
        document.querySelectorAll('.tab-content').forEach(function (c) { c.classList.remove('active'); });
        this.classList.add('active');
        const target = document.getElementById('tab-' + tab);
        if (target) target.classList.add('active');
      });
    });
  }

  // ══ FAQ ACCORDION ══
  document.querySelectorAll('.faq-question').forEach(function (q) {
    q.addEventListener('click', function () {
      const item = this.parentElement;
      const wasOpen = item.classList.contains('open');
      // Close all
      document.querySelectorAll('.faq-item').forEach(function (i) { i.classList.remove('open'); });
      if (!wasOpen) item.classList.add('open');
    });
  });

  // ══ SCROLL ANIMATIONS ══
  const observer = new IntersectionObserver(function (entries) {
    entries.forEach(function (entry) {
      if (entry.isIntersecting) {
        entry.target.classList.add('visible');
        observer.unobserve(entry.target);
      }
    });
  }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

  document.querySelectorAll('.animate-in').forEach(function (el) {
    observer.observe(el);
  });

  // ══ ANIMATED COUNTERS ══
  document.querySelectorAll('[data-count]').forEach(function (el) {
    const target = parseInt(el.dataset.count, 10);
    if (isNaN(target)) return;
    const io = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          animateCounter(el, target);
          io.unobserve(el);
        }
      });
    }, { threshold: 0.5 });
    io.observe(el);
  });
});

function animateCounter(el, target) {
  const duration = 1500;
  const start = performance.now();
  function tick(now) {
    const elapsed = now - start;
    const progress = Math.min(elapsed / duration, 1);
    const eased = 1 - Math.pow(1 - progress, 3); // ease-out cubic
    el.textContent = Math.floor(target * eased);
    if (progress < 1) requestAnimationFrame(tick);
    else el.textContent = target;
  }
  requestAnimationFrame(tick);
}
