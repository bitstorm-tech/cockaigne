/** @type {import("tailwindcss").Config} */
module.exports = {
  content: ["./internal/view/**/*.templ"],
  theme: {
    extend: {
      colors: {
        dark: "#1b2123"
      }
    }
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        dark: {
          ...require("daisyui/src/theming/themes")["dark"],
          "base-300": "#1b2123",
          "base-100": "#1b2123",
          primary: "#2c363a",
          warning: "#751b37"
        }
      }
    ]
  }
};
