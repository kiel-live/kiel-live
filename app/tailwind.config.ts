import type { Config } from 'tailwindcss';

export default {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        green: '#4caf50',
        link: '#60a5fa', // blue-400 equivalent

        // status colors
        'status-red': '#f87171', // red-400 equivalent
        'status-gray': '#9ca3af', // gray-400 equivalent
        'status-blue': '#60a5fa', // blue-400 equivalent
        'status-green': '#4caf50',

        // override red-500 to fix accessibility issues
        'red-500': 'rgb(250 68 68)',
      },
      boxShadow: {
        right: '4px 0 8px 0 rgba(17, 24, 39, 0.25)',
        top: '-4px 0 8px 0 rgba(17, 24, 39, 0.25)',
      },
      stroke: ({ theme }) => theme('colors'),
      fill: ({ theme }) => theme('colors'),
    },
  },
} satisfies Config;
