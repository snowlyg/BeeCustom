var {task, watch, src, dest, series,parallel } = require('gulp'),
    minifycss = require('gulp-minify-css'),
    jshint = require('gulp-jshint'),
    uglify = require('gulp-uglify'),
    imagemin = require('gulp-imagemin'),
    notify = require('gulp-notify'),
    cache = require('gulp-cache'),
    livereload = require('gulp-livereload'),
    babel = require('gulp-babel'),
    del = require('del');

//config 代码校验、合并和压缩
task('config', async function () {
    return src('customAdmin/config.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//login 代码校验、合并和压缩
task('login', async function () {
    return src('customAdmin/modules/login.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/modules'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//reset 代码校验、合并和压缩
task('reset', async function () {
    return src('customAdmin/modules/reset.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/modules'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//treetable 代码校验、合并和压缩
task('treetablejs', async function () {
    return src('customAdmin/modules/treetable-lay/treetable.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/modules/treetable-lay'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//treetable 代码校验、合并和压缩
task('treetablecss', async function () {
    return src('customAdmin/modules/treetable-lay/treetable.css')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/modules/treetable-lay'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//admin 代码校验、合并和压缩
task('admin', async function () {
    return src('customAdmin/lib/admin.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//AutoComplete 代码校验、合并和压缩
task('AutoComplete', async function () {
    return src('customAdmin/lib/AutoComplete.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//index 代码校验、合并和压缩
task('index', async function () {
    return src('customAdmin/lib/index.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//view 代码校验、合并和压缩
task('view', async function () {
    return src('customAdmin/lib/view.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//echarts 代码校验、合并和压缩
task('echarts', async function () {
    return src('customAdmin/lib/extend/echarts.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/extend'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//echartsTheme 代码校验、合并和压缩
task('echartsTheme', async function () {
    return src('customAdmin/lib/extend/echartsTheme.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/extend'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//tablePlug 代码校验、合并和压缩
task('tablePlugjs', async function () {
    return src('customAdmin/lib/tablePlug/tablePlug.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//optimizeSelectOption 代码校验、合并和压缩
task('optimizeSelectOptioncss', async function () {
    return src('customAdmin/lib/tablePlug/optimizeSelectOption/optimizeSelectOption.css')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug/optimizeSelectOption'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//optimizeSelectOption 代码校验、合并和压缩
task('optimizeSelectOptionjs', async function () {
    return src('customAdmin/lib/tablePlug/optimizeSelectOption/optimizeSelectOption.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug/optimizeSelectOption'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//tablePlug 代码校验、合并和压缩
task('tablePlugcss', async function () {
    return src('customAdmin/lib/tablePlug/tablePlug.css')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//layuiXtree 代码校验、合并和压缩
task('layuiXtree', async function () {
    return src('customAdmin/lib/extend/layuiXtree.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/extend'))
        .pipe(notify({message: 'Scripts task complete'}));
});


//压缩图片
task('images', function () {
    return src('src/images/**/*')
        .pipe(cache(imagemin({optimizationLevel: 5, progressive: true, interlaced: true})))
        .pipe(dest('dist/assets/img'))
        .pipe(notify({message: 'Images task complete'}));
});

//清除文件
task('clean', async function () {
    await del([
        'static/customAdmin/',
        'static/customAdmin/lib',
        'static/customAdmin/lib/extend',
        'static/customAdmin/lib/tablePlug',
        'static/customAdmin/modules',
        'static/customAdmin/modules/treetable-lay',
        'static/customAdmin/lib/tablePlug/optimizeSelectOption/',
    ])
});

task('watch', function () {
    // Watch .scss files
    watch('src/styles/**/*.scss', ['styles']);
    // Watch .js files
    watch('src/scripts/**/*.js', ['scripts']);
    // Watch image files
    watch('src/images/**/*', ['images']);
    // Create LiveReload server
    livereload.listen();
    // Watch any files in dist/, reload on change
    watch(['dist/**']).on('change', livereload.changed);
});

task('default', series(
  // 'clean',
  'config',
  'login',
  'reset',
  'treetablejs',
  // 'treetablecss',
  'admin',
  'AutoComplete',
  'index',
  'view',
  'echarts',
  'echartsTheme',
  'tablePlugjs',
  // 'optimizeSelectOptioncss',
  'optimizeSelectOptionjs',
  // 'tablePlugcss',
  'layuiXtree',
));