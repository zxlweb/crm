<template>
  <svg class="w-full max-w-lg opacity-95" viewBox="0 0 520 380" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
    <ellipse cx="260" cy="200" rx="200" ry="140" fill="#7C3AED" opacity="0.08" />
    <ellipse cx="400" cy="100" rx="90" ry="70" fill="#A78BFA" opacity="0.06" />

    <g transform="translate(72 48)">
      <rect x="0" y="0" width="376" height="268" rx="20" fill="#141414" stroke="rgba(255,255,255,0.08)" stroke-width="1" />
      <rect x="0" y="0" width="376" height="52" rx="20" fill="#1A1A1A" />
      <rect x="0" y="36" width="376" height="16" fill="#1A1A1A" />
      <circle cx="28" cy="26" r="5" fill="#EF4444" opacity="0.7" />
      <circle cx="48" cy="26" r="5" fill="#EAB308" opacity="0.7" />
      <circle cx="68" cy="26" r="5" fill="#22C55E" opacity="0.7" />

      <rect x="16" y="64" width="72" height="188" rx="12" fill="#111111" />
      <rect x="28" y="80" width="48" height="8" rx="4" fill="#2A2A2A" />
      <rect x="28" y="168" width="48" height="32" rx="8" fill="url(#darkAccent)" />

      <rect x="104" y="64" width="96" height="56" rx="12" fill="#111111" stroke="rgba(255,255,255,0.05)" />
      <text x="116" y="86" fill="#A78BFA" font-size="10" font-family="system-ui">Revenue</text>
      <text x="116" y="108" fill="#F5F3FF" font-size="18" font-family="system-ui" font-weight="700">$88.6k</text>

      <rect x="104" y="132" width="256" height="120" rx="12" fill="#111111" stroke="rgba(255,255,255,0.05)" />
      <path d="M118 220 L158 175 L198 190 L238 160 L278 170 L318 130 L358 145" stroke="url(#darkLine)" stroke-width="2.5" fill="none" stroke-linecap="round" />
      <path d="M118 220 L358 145 L358 220 L118 220 Z" fill="url(#darkArea)" opacity="0.4" />
      <circle cx="318" cy="130" r="4" fill="#A78BFA" stroke="#141414" stroke-width="2" />
    </g>

    <defs>
      <linearGradient id="darkAccent" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#7C3AED" />
        <stop offset="1" stop-color="#5B21B6" />
      </linearGradient>
      <linearGradient id="darkLine" x1="0" y1="0" x2="1" y2="0">
        <stop stop-color="#7C3AED" />
        <stop offset="1" stop-color="#C4B5FD" />
      </linearGradient>
      <linearGradient id="darkArea" x1="0" y1="0" x2="0" y2="1">
        <stop stop-color="#7C3AED" stop-opacity="0.5" />
        <stop offset="1" stop-color="#7C3AED" stop-opacity="0" />
      </linearGradient>
    </defs>
  </svg>
</template>
