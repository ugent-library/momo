const mix = require('laravel-mix')
const ESLintPlugin = require('eslint-webpack-plugin')

mix.js('./themes/orpheus/js/app.js', 'js')
  .sass('./themes/orpheus/css/app.scss', 'css')
  .setPublicPath('./static/orpheus')
  .vue()

if (mix.inProduction()) {
  mix.version()
} else {
  mix.sourceMaps()
}

// Copy font files to the ./static/ugent/fonts folder
// webpackConfig does a deep merge. The "test:<value>" key needs to match with the
// exact line as generated by mix in order to replace the generated config at that location.
mix.webpackConfig({
  plugins: [
    new ESLintPlugin({})
  ]
})

// Uncomment this if you want to see generated webpack.config.js
// mix.dump();
