let mix = require('laravel-mix');

mix.js('./themes/orpheus/js/app.js', 'js')
    .sass('./themes/orpheus/css/app.scss', 'css')
    .setPublicPath('./static/orpheus')
    .vue();

if (mix.inProduction()) {
    mix.version()
} else {
    mix.sourceMaps();
}
