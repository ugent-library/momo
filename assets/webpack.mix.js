let mix = require('laravel-mix');

mix.js('./js/app.js', 'js')
    .sass('./css/app.scss', 'css')
    .setPublicPath('../static')
    .version();

// if (mix.inProduction) {
// } else {
// }