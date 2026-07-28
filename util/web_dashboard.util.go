package util

import (
	"fmt"
	"net/http"
)

const gitDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Re-writing History — Git History & Commit Graph Telemetry</title>
  <meta name="description" content="Production-grade Git history rewriting, commit manipulation, and AI timestamp graph inspector.">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  <style>
    :root {
      color-scheme: dark;
      --bg-base: #09090b;
      --bg-raised: #121318;
      --bg-overlay: #1c1d24;
      --border-default: #27272a;
      --border-subtle: #1f1f23;
      --border-focus: #6366f1;
      --text-primary: #f4f4f5;
      --text-secondary: #a1a1aa;
      --text-muted: #71717a;
      --accent: #6366f1;
      --accent-muted: rgba(99, 102, 241, 0.15);
      --success: #10b981;
      --success-bg: rgba(16, 185, 129, 0.12);
      --font-sans: 'Inter', system-ui, sans-serif;
      --font-mono: 'JetBrains Mono', monospace;
      --radius-sm: 6px;
      --radius-md: 10px;
      --radius-lg: 14px;
    }
    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: var(--bg-base);
      color: var(--text-primary);
      font-family: var(--font-sans);
      font-size: 14px;
      line-height: 1.6;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }
    .header {
      background: var(--bg-raised);
      border-bottom: 1px solid var(--border-default);
      position: sticky;
      top: 0;
      z-index: 100;
    }
    .header-inner {
      max-width: 1400px;
      margin: 0 auto;
      padding: 16px 24px;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
    .brand { display: flex; align-items: center; gap: 12px; }
    .brand-icon {
      width: 32px;
      height: 32px;
      background: linear-gradient(135deg, #8b5cf6, #6366f1);
      border-radius: var(--radius-sm);
      display: grid;
      place-items: center;
      font-weight: 700;
      color: #fff;
    }
    .brand-title { font-size: 16px; font-weight: 700; }
    .main {
      max-width: 1400px;
      margin: 0 auto;
      width: 100%;
      padding: 32px 24px;
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 28px;
    }
    .hero {
      background: var(--bg-raised);
      border: 1px solid var(--border-default);
      border-radius: var(--radius-lg);
      padding: 24px 28px;
      animation: slideFadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards;
    }
    .hero-title { font-size: 22px; font-weight: 700; margin-bottom: 6px; }
    .hero-sub { color: var(--text-secondary); font-size: 14px; }
    .commit-list {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }
    .commit-card {
      background: var(--bg-raised);
      border: 1px solid var(--border-default);
      border-radius: var(--radius-md);
      padding: 16px 20px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      animation: slideFadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
      transition: transform 0.2s ease, border-color 0.2s ease;
    }

    .commit-card:hover {
      transform: translateX(4px);
      border-color: var(--border-focus);
    }

    .commit-card:nth-child(1) { animation-delay: 0.05s; }
    .commit-card:nth-child(2) { animation-delay: 0.10s; }
    .commit-card:nth-child(3) { animation-delay: 0.15s; }

    @keyframes slideFadeIn {
      from {
        opacity: 0;
        transform: translateY(14px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    @keyframes pulseEngine {
      0%, 100% { opacity: 1; transform: scale(1); }
      50% { opacity: 0.5; transform: scale(1.1); }
    }

    .engine-dot {
      display: inline-block;
      width: 8px;
      height: 8px;
      background: var(--success);
      border-radius: 50%;
      margin-right: 6px;
      animation: pulseEngine 2s infinite ease-in-out;
    }

    .commit-sha {
      font-family: var(--font-mono);
      font-size: 12px;
      background: var(--accent-muted);
      color: var(--accent);
      padding: 4px 8px;
      border-radius: 4px;
      font-weight: 600;
      transition: transform 0.15s ease, background-color 0.15s ease;
    }

    .commit-sha:hover {
      transform: scale(1.05);
      background-color: var(--accent);
      color: #ffffff;
    }

    .commit-msg { font-size: 14px; font-weight: 500; }
    .commit-date { font-family: var(--font-mono); font-size: 12px; color: var(--text-muted); }

    @media (prefers-reduced-motion: reduce) {
      .hero, .commit-card, .engine-dot, .commit-sha {
        animation: none !important;
        transition: none !important;
        transform: none !important;
      }
    }
  </style>
</head>
<body>
  <header class="header">
    <div class="header-inner">
      <div class="brand">
        <div class="brand-icon">📜</div>
        <span class="brand-title">Git History Rewriter</span>
      </div>
      <span style="font-family: var(--font-mono); font-size: 12px; color: var(--text-muted);"><span class="engine-dot"></span>Status: Active Engine</span>
    </div>
  </header>

  <main class="main">
    <section class="hero">
      <h1 class="hero-title">Commit Graph & Timestamp Manipulation Engine</h1>
      <p class="hero-sub">AI-assisted Git history synthesizer, commit date normalization, and automated remote synchronization platform.</p>
    </section>

    <section>
      <h2 style="font-size: 16px; font-weight: 600; margin-bottom: 16px;">Rewritten Commit History Log</h2>
      <div class="commit-list">
        <div class="commit-card">
          <div style="display: flex; align-items: center; gap: 14px;">
            <span class="commit-sha">a1b2c3d</span>
            <span class="commit-msg">feat(core): initialize automated git history rewriting workflow</span>
          </div>
          <span class="commit-date">2026-06-17 10:14:00 +0530</span>
        </div>
        <div class="commit-card">
          <div style="display: flex; align-items: center; gap: 14px;">
            <span class="commit-sha">e4f5g6h</span>
            <span class="commit-msg">refactor(ai): sanitize commit timestamps and author signatures</span>
          </div>
          <span class="commit-date">2026-06-18 14:30:00 +0530</span>
        </div>
        <div class="commit-card">
          <div style="display: flex; align-items: center; gap: 14px;">
            <span class="commit-sha">7890xyz</span>
            <span class="commit-msg">chore(sync): force update SSH target repository graph</span>
          </div>
          <span class="commit-date">2026-06-19 18:45:00 +0530</span>
        </div>
      </div>
    </section>
  </main>
</body>
</html>`

func ServeGitDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, gitDashboardHTML)
}
