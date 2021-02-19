let mix = require('laravel-mix');

mix.js('./js/app.js', 'js')
    .sass('./css/app.scss', 'css')
    .setPublicPath('../static')
    .vue();

if (mix.inProduction()) {
    mix.version()
} else {
    mix.sourceMaps();
}