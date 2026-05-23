<template>
  <svg
    class="dashboard-workbench-illustration__svg w-full"
    viewBox="0 0 480 320"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    :aria-label="$t('dashboardIllustrationAria')"
  >
    <!-- 背景光晕 -->
    <ellipse cx="240" cy="168" rx="170" ry="118" fill="url(#db-glow)" opacity="0.65" />
    <ellipse cx="380" cy="88" rx="64" ry="48" fill="#DDD6FE" opacity="0.45" />
    <ellipse cx="96" cy="248" rx="56" ry="40" fill="#F3E8FF" opacity="0.75" />

    <!-- 装饰节点 -->
    <circle cx="72" cy="56" r="5" fill="#A78BFA" opacity="0.55" />
    <circle cx="408" cy="44" r="4" fill="#C4B5FD" />
    <circle cx="432" cy="268" r="6" fill="#8B5CF6" opacity="0.3" />
    <path d="M72 62 L108 96" stroke="#C4B5FD" stroke-width="1.5" stroke-dasharray="4 4" />
    <path d="M408 48 L372 82" stroke="#DDD6FE" stroke-width="1.5" />

    <!-- 主面板 · 每日作战台 -->
    <g transform="translate(56 36) rotate(-1.5 184 130)">
      <rect x="0" y="0" width="368" height="248" rx="18" fill="white" stroke="#E9D5FF" stroke-width="1.5" />

      <!-- 顶栏 -->
      <rect x="0" y="0" width="368" height="46" rx="18" fill="url(#db-header)" />
      <rect x="0" y="30" width="368" height="16" fill="url(#db-header)" />
      <circle cx="24" cy="23" r="4.5" fill="#FCA5A5" opacity="0.85" />
      <circle cx="40" cy="23" r="4.5" fill="#FDE68A" opacity="0.85" />
      <circle cx="56" cy="23" r="4.5" fill="#86EFAC" opacity="0.85" />

      <!-- 闪电 + 标题 -->
      <rect x="76" y="14" width="18" height="18" rx="6" fill="url(#db-bolt)" />
      <path d="M85 18v4l-2 3h3l-2 4" stroke="white" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
      <rect x="102" y="17" width="72" height="6" rx="3" fill="#C4B5FD" opacity="0.7" />
      <rect x="102" y="27" width="48" height="4" rx="2" fill="#DDD6FE" />

      <!-- 优先卡 1 · 预警 -->
      <g transform="translate(16 58)">
        <rect width="336" height="52" rx="10" fill="#FAF5FF" stroke="#FECACA" stroke-width="1" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#EF4444" />
        <circle cx="28" cy="26" r="14" fill="#FEE2E2" />
        <path d="M28 20v4M28 30h.01" stroke="#EF4444" stroke-width="1.5" stroke-linecap="round" />
        <rect x="50" y="14" width="88" height="6" rx="3" fill="#5B21B6" opacity="0.85" />
        <rect x="50" y="26" width="120" height="4" rx="2" fill="#DDD6FE" />
        <rect x="50" y="34" width="96" height="4" rx="2" fill="#EDE9FE" />
        <!-- sparkline -->
        <path d="M248 34 L258 28 L268 30 L278 24 L288 26" stroke="#F97316" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#EA580C" font-size="11" font-family="system-ui,sans-serif" font-weight="700">28</text>
        <rect x="312" y="18" width="16" height="16" rx="8" fill="url(#db-bolt)" opacity="0.9" />
        <path d="M318 24h4M320 22v4" stroke="white" stroke-width="1.2" stroke-linecap="round" />
      </g>

      <!-- 优先卡 2 · 观察 -->
      <g transform="translate(16 118)">
        <rect width="336" height="52" rx="10" fill="#FAF5FF" stroke="#FDE68A" stroke-width="1" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#F59E0B" />
        <circle cx="28" cy="26" r="14" fill="#FEF3C7" />
        <circle cx="28" cy="26" r="5" fill="#F59E0B" opacity="0.35" />
        <rect x="50" y="14" width="72" height="6" rx="3" fill="#5B21B6" opacity="0.85" />
        <rect x="50" y="26" width="108" height="4" rx="2" fill="#DDD6FE" />
        <path d="M248 34 L260 30 L272 32 L284 28" stroke="#8B5CF6" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#7C3AED" font-size="11" font-family="system-ui,sans-serif" font-weight="700">62</text>
      </g>

      <!-- 优先卡 3 · 健康 -->
      <g transform="translate(16 178)">
        <rect width="336" height="52" rx="10" fill="#FAF5FF" stroke="#BBF7D0" stroke-width="1" opacity="0.85" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#22C55E" />
        <circle cx="28" cy="26" r="14" fill="#DCFCE7" />
        <path d="M23 26l3 3 6-6" stroke="#16A34A" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        <rect x="50" y="14" width="64" height="6" rx="3" fill="#DDD6FE" />
        <rect x="50" y="26" width="92" height="4" rx="2" fill="#EDE9FE" />
        <path d="M248 32 L262 28 L276 30 L290 24" stroke="#22C55E" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#16A34A" font-size="11" font-family="system-ui,sans-serif" font-weight="700">88</text>
      </g>
    </g>

    <!-- 浮动 · 需关注徽章 -->
    <g transform="translate(8 12)" filter="url(#db-shadow)">
      <rect width="96" height="40" rx="12" fill="white" stroke="#E9D5FF" stroke-width="1.5" />
      <circle cx="22" cy="20" r="12" fill="#FEE2E2" />
      <text x="22" y="24" text-anchor="middle" fill="#DC2626" font-size="11" font-family="system-ui,sans-serif" font-weight="700">3</text>
      <text x="40" y="18" fill="#7C3AED" font-size="8" font-family="system-ui,sans-serif" font-weight="600">Priority</text>
      <text x="40" y="30" fill="#A78BFA" font-size="7" font-family="system-ui,sans-serif">Today</text>
    </g>

    <!-- 浮动 · 跟进 CTA -->
    <g transform="translate(352 248)" filter="url(#db-shadow)">
      <rect width="112" height="44" rx="12" fill="url(#db-bolt)" />
      <text x="16" y="20" fill="white" font-size="8" font-family="system-ui,sans-serif" font-weight="600" opacity="0.85">Follow up</text>
      <text x="16" y="34" fill="white" font-size="11" font-family="system-ui,sans-serif" font-weight="700">→ Timeline</text>
      <circle cx="88" cy="22" r="14" fill="rgba(255,255,255,0.2)" />
      <path d="M84 22h8M88 18v8" stroke="white" stroke-width="1.5" stroke-linecap="round" opacity="0.9" />
    </g>

    <!-- 连接线 -->
    <path d="M180 290 Q240 260 300 230" stroke="url(#db-line)" stroke-width="1.5" opacity="0.35" stroke-dasharray="5 4" />

    <defs>
      <linearGradient id="db-glow" x1="0" y1="0" x2="1" y2="1">
        <stop stop-color="#EDE9FE" />
        <stop offset="1" stop-color="#F5F3FF" />
      </linearGradient>
      <linearGradient id="db-header" x1="0" y1="0" x2="1" y2="0">
        <stop stop-color="#F5F3FF" />
        <stop offset="1" stop-color="#EDE9FE" />
      </linearGradient>
      <linearGradient id="db-bolt" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#A78BFA" />
        <stop offset="1" stop-color="#7C3AED" />
      </linearGradient>
      <linearGradient id="db-line" x1="0" y1="0" x2="1" y2="0">
        <stop stop-color="#C4B5FD" />
        <stop offset="1" stop-color="#8B5CF6" />
      </linearGradient>
      <filter id="db-shadow" x="-20%" y="-20%" width="140%" height="140%">
        <feDropShadow dx="0" dy="6" stdDeviation="10" flood-color="#7C3AED" flood-opacity="0.14" />
      </filter>
    </defs>
  </svg>
</template>
