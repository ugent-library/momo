let mix = require('laravel-mix');

mix.js('./themes/opale/js/app.js', 'js')
    .sass('./themes/opale/css/app.scss', 'css')
    .setPublicPath('./static/opale')
    .vue();

if (mix.inProduction()) {
    mix.version()
} else {
    mix.sourceMaps();
}
