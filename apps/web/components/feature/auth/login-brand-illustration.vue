<template>
  <svg
    class="w-full max-w-lg drop-shadow-2xl"
    viewBox="0 0 520 380"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    role="img"
    aria-label="CRM analytics illustration"
  >
    <!-- 背景光晕 -->
    <ellipse cx="260" cy="200" rx="200" ry="140" fill="url(#glow)" opacity="0.6" />
    <ellipse cx="380" cy="120" rx="80" ry="60" fill="#DDD6FE" opacity="0.5" />
    <ellipse cx="120" cy="280" rx="70" ry="50" fill="#F3E8FF" opacity="0.8" />

    <!-- 浮动装饰节点 -->
    <circle cx="88" cy="72" r="6" fill="#A78BFA" opacity="0.6" />
    <circle cx="430" cy="58" r="4" fill="#C4B5FD" />
    <circle cx="460" cy="300" r="8" fill="#8B5CF6" opacity="0.35" />
    <path d="M88 78 L140 120" stroke="#C4B5FD" stroke-width="1.5" stroke-dasharray="4 4" />
    <path d="M430 62 L390 100" stroke="#DDD6FE" stroke-width="1.5" />

    <!-- 主面板（等角透视感） -->
    <g transform="translate(72 48) rotate(-2 196 160)">
      <rect x="0" y="0" width="376" height="268" rx="20" fill="white" stroke="#E9D5FF" stroke-width="1.5" />
      <rect x="0" y="0" width="376" height="52" rx="20" fill="url(#headerGrad)" />
      <rect x="0" y="36" width="376" height="16" fill="url(#headerGrad)" />

      <!-- 窗口按钮 -->
      <circle cx="28" cy="26" r="5" fill="#FCA5A5" opacity="0.8" />
      <circle cx="48" cy="26" r="5" fill="#FDE68A" opacity="0.8" />
      <circle cx="68" cy="26" r="5" fill="#86EFAC" opacity="0.8" />

      <!-- 侧栏 -->
      <rect x="16" y="64" width="72" height="188" rx="12" fill="#FAF5FF" />
      <rect x="28" y="80" width="48" height="8" rx="4" fill="#E9D5FF" />
      <rect x="28" y="98" width="40" height="6" rx="3" fill="#DDD6FE" />
      <rect x="28" y="116" width="44" height="6" rx="3" fill="#C4B5FD" />
      <rect x="28" y="134" width="36" height="6" rx="3" fill="#DDD6FE" />
      <rect x="28" y="168" width="48" height="32" rx="8" fill="url(#sidebarAccent)" />

      <!-- KPI 卡片 -->
      <rect x="104" y="64" width="96" height="56" rx="12" fill="#FAF5FF" stroke="#F3E8FF" />
      <text x="116" y="86" fill="#6B21A8" font-size="10" font-family="system-ui,sans-serif" font-weight="600">Leads</text>
      <text x="116" y="108" fill="#4C1D95" font-size="18" font-family="system-ui,sans-serif" font-weight="700">2,847</text>
      <rect x="168" y="92" width="24" height="14" rx="7" fill="#D1FAE5" />
      <text x="174" y="102" fill="#059669" font-size="8" font-family="system-ui,sans-serif" font-weight="600">+18%</text>

      <rect x="212" y="64" width="96" height="56" rx="12" fill="#FAF5FF" stroke="#F3E8FF" />
      <text x="224" y="86" fill="#6B21A8" font-size="10" font-family="system-ui,sans-serif" font-weight="600">Deals</text>
      <text x="224" y="108" fill="#4C1D95" font-size="18" font-family="system-ui,sans-serif" font-weight="700">¥1.2M</text>

      <rect x="320" y="64" width="40" height="56" rx="12" fill="url(#purpleBtn)" />
      <path d="M332 88 L348 88 M340 80 L340 96" stroke="white" stroke-width="2" stroke-linecap="round" />

      <!-- 折线图区域 -->
      <rect x="104" y="132" width="256" height="120" rx="12" fill="#FAF5FF" stroke="#F3E8FF" />
      <text x="118" y="152" fill="#7C3AED" font-size="10" font-family="system-ui,sans-serif" font-weight="600">Pipeline</text>

      <!-- 网格线 -->
      <line x1="118" y1="220" x2="344" y2="220" stroke="#EDE9FE" stroke-width="1" />
      <line x1="118" y1="195" x2="344" y2="195" stroke="#EDE9FE" stroke-width="1" />
      <line x1="118" y1="170" x2="344" y2="170" stroke="#EDE9FE" stroke-width="1" />

      <!-- 折线 -->
      <path
        d="M118 210 L158 185 L198 195 L238 165 L278 175 L318 140 L358 155"
        stroke="url(#lineGrad)"
        stroke-width="3"
        stroke-linecap="round"
        stroke-linejoin="round"
        fill="none"
      />
      <path
        d="M118 210 L158 185 L198 195 L238 165 L278 175 L318 140 L358 155 L358 220 L118 220 Z"
        fill="url(#areaGrad)"
        opacity="0.5"
      />

      <!-- 数据点 -->
      <circle cx="318" cy="140" r="5" fill="#8B5CF6" stroke="white" stroke-width="2" />
    </g>

    <!-- 前景浮动卡片 -->
    <g transform="translate(8 260)">
      <rect width="140" height="72" rx="14" fill="white" stroke="#E9D5FF" stroke-width="1.5" filter="url(#cardShadow)" />
      <circle cx="28" cy="36" r="16" fill="url(#avatarGrad)" />
      <text x="52" y="32" fill="#5B21B6" font-size="11" font-family="system-ui,sans-serif" font-weight="600">New lead</text>
      <text x="52" y="48" fill="#A78BFA" font-size="9" font-family="system-ui,sans-serif">Acme Corp · Qualified</text>
      <rect x="52" y="54" width="72" height="6" rx="3" fill="#EDE9FE" />
      <rect x="52" y="54" width="48" height="6" rx="3" fill="#A78BFA" />
    </g>

    <g transform="translate(360 24)">
      <rect width="120" height="56" rx="14" fill="white" stroke="#E9D5FF" stroke-width="1.5" filter="url(#cardShadow)" />
      <text x="16" y="24" fill="#7C3AED" font-size="9" font-family="system-ui,sans-serif">Retention</text>
      <text x="16" y="42" fill="#4C1D95" font-size="16" font-family="system-ui,sans-serif" font-weight="700">94.2%</text>
      <circle cx="96" cy="28" r="18" fill="none" stroke="#E9D5FF" stroke-width="4" />
      <circle cx="96" cy="28" r="18" fill="none" stroke="#8B5CF6" stroke-width="4" stroke-dasharray="85 113" stroke-linecap="round" transform="rotate(-90 96 28)" />
    </g>

    <!-- 科技连接线 -->
    <path d="M200 340 Q260 300 320 260" stroke="url(#lineGrad)" stroke-width="2" opacity="0.4" stroke-dasharray="6 4" />

    <defs>
      <linearGradient id="glow" x1="0" y1="0" x2="1" y2="1">
        <stop stop-color="#EDE9FE" />
        <stop offset="1" stop-color="#F5F3FF" />
      </linearGradient>
      <linearGradient id="headerGrad" x1="0" y1="0" x2="1" y2="0">
        <stop stop-color="#F5F3FF" />
        <stop offset="1" stop-color="#EDE9FE" />
      </linearGradient>
      <linearGradient id="sidebarAccent" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#A78BFA" />
        <stop offset="1" stop-color="#7C3AED" />
      </linearGradient>
      <linearGradient id="purpleBtn" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#A78BFA" />
        <stop offset="1" stop-color="#7C3AED" />
      </linearGradient>
      <linearGradient id="lineGrad" x1="0" y1="0" x2="1" y2="0">
        <stop stop-color="#C4B5FD" />
        <stop offset="0.5" stop-color="#8B5CF6" />
        <stop offset="1" stop-color="#6D28D9" />
      </linearGradient>
      <linearGradient id="areaGrad" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#C4B5FD" stop-opacity="0.5" />
        <stop offset="1" stop-color="#EDE9FE" stop-opacity="0" />
      </linearGradient>
      <linearGradient id="avatarGrad" x1="0" y1="0" x2="1" y2="1">
        <stop stop-color="#DDD6FE" />
        <stop offset="1" stop-color="#8B5CF6" />
      </linearGradient>
      <filter id="cardShadow" x="-20%" y="-20%" width="140%" height="140%">
        <feDropShadow dx="0" dy="8" stdDeviation="12" flood-color="#7C3AED" flood-opacity="0.12" />
      </filter>
    </defs>
  </svg>
</template>
