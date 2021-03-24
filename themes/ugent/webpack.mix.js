const mix = require('laravel-mix')
const ESLintPlugin = require('eslint-webpack-plugin');

mix.js('./themes/ugent/js/app.js', 'js')
  .sass('./themes/ugent/scss/screen.scss', 'css')
  .setPublicPath('./static/ugent')
  .vue()

// Set the resourceroot for fonts so it points to the static assets path
mix.setResourceRoot('/s/ugent/fonts/')

// Copy images
mix.copyDirectory('./themes/_common/ugent/images', './static/ugent/images')

if (mix.inProduction()) {
  mix.version()
} else {
  mix.sourceMaps()
}

// Copy font files to the ./static/ugent/fonts folder
// webpackConfig does a deep merge. The "test:<value>" key needs to match with the
// exact line as generated by mix in order to replace the generated config at that location.
mix.webpackConfig({
  module: {
    rules: [
      {
        test: '/(\\.(woff2?|ttf|eot|otf)$|font.*\\.svg$)/',
        use: [{
          loader: 'file-loader',
          options: {
            name: '[name].[ext]',
            outputPath: './fonts/'
          }
        }]
      }
    ]
  },
  plugins: [
    new ESLintPlugin({})
  ],
})

// Uncomment this if you want to see generated webpack.config.js
// mix.dump();
