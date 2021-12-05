import { defineConfig } from 'windicss/helpers';
import typography from 'windicss/plugin/typography';
import colors from 'windicss/colors';

export default defineConfig({
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        green: '#4caf50',
        link: colors.blue[400],

        // status colors
        'status-red': colors.red[400],
        'status-gray': colors.gray[400],
        'status-blue': colors.blue[400],
        'status-green': '#4caf50',
      },
      boxShadow: {
        all: '0px 0px 8px 4px rgba(17, 24, 39, 0.25)',
      },
      stroke: (theme) => theme('colors'),
      fill: (theme) => theme('colors'),
    },
  },
  plugins: [typography],
});
