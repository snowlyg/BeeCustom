var {task, watch, src, dest, series,parallel } = require('gulp'),
    sass = require('gulp-ruby-sass'),
    autoprefixer = require('gulp-autoprefixer'),
    minifycss = require('gulp-minify-css'),
    jshint = require('gulp-jshint'),
    uglify = require('gulp-uglify'),
    imagemin = require('gulp-imagemin'),
    rename = require('gulp-rename'),
    concat = require('gulp-concat'),
    notify = require('gulp-notify'),
    cache = require('gulp-cache'),
    livereload = require('gulp-livereload'),
    babel = require('gulp-babel'),
    del = require('del');

//customAdmin_lib 代码校验、合并和压缩
task('customAdmin_lib', async function () {
    return src('static/customAdmin/lib/*.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        // .pipe(concat('main.js'))
        // .pipe(dest('dist/assets/js'))
        // .pipe(rename({suffix: '.min'}))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/dest'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//autocomplete 代码校验、合并和压缩
task('autocomplete', async function () {
    return src('static/customAdmin/lib/autocomplete/*.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        // .pipe(concat('main.js'))
        // .pipe(dest('dist/assets/js'))
        // .pipe(rename({suffix: '.min'}))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/autocomplete/dest'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//extend 代码校验、合并和压缩
task('extend', async function () {
    return src('static/customAdmin/lib/extend/*.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        // .pipe(concat('main.js'))
        // .pipe(dest('dist/assets/js'))
        // .pipe(rename({suffix: '.min'}))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/extend/dest'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//tablePlug 代码校验、合并和压缩
task('tablePlug', async function () {
    return src('static/customAdmin/lib/tablePlug/*.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        // .pipe(concat('main.js'))
        // .pipe(dest('dist/assets/js'))
        // .pipe(rename({suffix: '.min'}))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug/dest'))
        .pipe(notify({message: 'Scripts task complete'}));
});

//optimizeSelectOption 代码校验、合并和压缩
task('optimizeSelectOption', async function () {
    return src('static/customAdmin/lib/tablePlug/optimizeSelectOption/*.js')
        .pipe(babel())
        .pipe(jshint())
        .pipe(jshint.reporter('default'))
        // .pipe(concat('main.js'))
        // .pipe(dest('dist/assets/js'))
        // .pipe(rename({suffix: '.min'}))
        .pipe(uglify())
        .pipe(dest('static/customAdmin/lib/tablePlug/optimizeSelectOption/dest'))
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
        'static/customAdmin/lib/dest',
        'static/customAdmin/lib/autocomplete/dest',
        'static/customAdmin/lib/extend/dest',
        'static/customAdmin/lib/tablePlug/dest',
        'static/customAdmin/lib/tablePlug/optimizeSelectOption/dest',
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

task('default', series('clean','optimizeSelectOption','tablePlug','extend','autocomplete','customAdmin_lib'));