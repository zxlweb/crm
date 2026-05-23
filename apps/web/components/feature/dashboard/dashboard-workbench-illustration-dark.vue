<template>
  <svg
    class="dashboard-workbench-illustration__svg w-full opacity-95"
    viewBox="0 0 480 320"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    aria-hidden="true"
  >
    <ellipse cx="240" cy="168" rx="170" ry="118" fill="#7C3AED" opacity="0.1" />
    <ellipse cx="392" cy="84" rx="70" ry="52" fill="#A78BFA" opacity="0.06" />

    <g transform="translate(56 36) rotate(-1.5 184 130)">
      <rect x="0" y="0" width="368" height="248" rx="18" fill="#141414" stroke="rgba(255,255,255,0.08)" stroke-width="1" />
      <rect x="0" y="0" width="368" height="46" rx="18" fill="#1A1A1A" />
      <rect x="0" y="30" width="368" height="16" fill="#1A1A1A" />
      <circle cx="24" cy="23" r="4.5" fill="#EF4444" opacity="0.65" />
      <circle cx="40" cy="23" r="4.5" fill="#EAB308" opacity="0.65" />
      <circle cx="56" cy="23" r="4.5" fill="#22C55E" opacity="0.65" />

      <rect x="76" y="14" width="18" height="18" rx="6" fill="url(#db-dark-bolt)" />
      <path d="M85 18v4l-2 3h3l-2 4" stroke="white" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" />
      <rect x="102" y="17" width="72" height="6" rx="3" fill="#4C1D95" opacity="0.8" />
      <rect x="102" y="27" width="48" height="4" rx="2" fill="#2A2A2A" />

      <g transform="translate(16 58)">
        <rect width="336" height="52" rx="10" fill="#111111" stroke="rgba(239,68,68,0.35)" stroke-width="1" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#EF4444" />
        <circle cx="28" cy="26" r="14" fill="rgba(239,68,68,0.15)" />
        <path d="M28 20v4M28 30h.01" stroke="#F87171" stroke-width="1.5" stroke-linecap="round" />
        <rect x="50" y="14" width="88" height="6" rx="3" fill="#A78BFA" opacity="0.7" />
        <rect x="50" y="26" width="120" height="4" rx="2" fill="#2A2A2A" />
        <path d="M248 34 L258 28 L268 30 L278 24 L288 26" stroke="#FB923C" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#FB923C" font-size="11" font-family="system-ui,sans-serif" font-weight="700">28</text>
      </g>

      <g transform="translate(16 118)">
        <rect width="336" height="52" rx="10" fill="#111111" stroke="rgba(245,158,11,0.3)" stroke-width="1" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#F59E0B" />
        <circle cx="28" cy="26" r="14" fill="rgba(245,158,11,0.12)" />
        <rect x="50" y="14" width="72" height="6" rx="3" fill="#A78BFA" opacity="0.7" />
        <rect x="50" y="26" width="108" height="4" rx="2" fill="#2A2A2A" />
        <path d="M248 34 L260 30 L272 32 L284 28" stroke="#A78BFA" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#C4B5FD" font-size="11" font-family="system-ui,sans-serif" font-weight="700">62</text>
      </g>

      <g transform="translate(16 178)">
        <rect width="336" height="52" rx="10" fill="#111111" stroke="rgba(34,197,94,0.25)" stroke-width="1" opacity="0.85" />
        <rect x="0" y="0" width="4" height="52" rx="2" fill="#22C55E" />
        <circle cx="28" cy="26" r="14" fill="rgba(34,197,94,0.12)" />
        <path d="M23 26l3 3 6-6" stroke="#4ADE80" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
        <rect x="50" y="14" width="64" height="6" rx="3" fill="#3F3F46" />
        <path d="M248 32 L262 28 L276 30 L290 24" stroke="#4ADE80" stroke-width="1.5" fill="none" stroke-linecap="round" />
        <text x="302" y="32" fill="#4ADE80" font-size="11" font-family="system-ui,sans-serif" font-weight="700">88</text>
      </g>
    </g>

    <g transform="translate(8 12)">
      <rect width="96" height="40" rx="12" fill="#141414" stroke="rgba(255,255,255,0.08)" stroke-width="1" />
      <circle cx="22" cy="20" r="12" fill="rgba(239,68,68,0.2)" />
      <text x="22" y="24" text-anchor="middle" fill="#F87171" font-size="11" font-family="system-ui,sans-serif" font-weight="700">3</text>
      <text x="40" y="18" fill="#A78BFA" font-size="8" font-family="system-ui,sans-serif" font-weight="600">Priority</text>
      <text x="40" y="30" fill="#71717A" font-size="7" font-family="system-ui,sans-serif">Today</text>
    </g>

    <g transform="translate(352 248)">
      <rect width="112" height="44" rx="12" fill="url(#db-dark-bolt)" />
      <text x="16" y="20" fill="white" font-size="8" font-family="system-ui,sans-serif" font-weight="600" opacity="0.85">Follow up</text>
      <text x="16" y="34" fill="white" font-size="11" font-family="system-ui,sans-serif" font-weight="700">→ Timeline</text>
    </g>

    <defs>
      <linearGradient id="db-dark-bolt" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#7C3AED" />
        <stop offset="1" stop-color="#5B21B6" />
      </linearGradient>
    </defs>
  </svg>
</template>
