/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../../internal/**/*.templ"],
  theme: {
    colors: {
      black: '#181818',
      gray: {
        '100': '#BBBBBB',
        '200': '#3A3A3A',
        '300': '#383838',
        '400': '#363636',
        '500': '#343434',
        '600': '#323232',
        '700': '#303030',
        '800': '#2e2e2e',
        '900': '#2d2d2d'
      },
      green: '#15803d',
      red: '#b91c1c'
    },
    extend: {},
  },
  plugins: [],
}

