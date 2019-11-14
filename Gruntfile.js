module.exports = function(grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        clean: {
            build: ['build']
        },
        copy: {
            main: {
                files: [
                    {
                        expand: true,
                        cwd: 'www',
                        src: '**',
                        dest: 'build',
                        flatten: false,
                        filter: 'isFile'
                    }
                ]
            }
        },
        uglify: {
            main: {
                options: {
                    sourceMap: false
                },
                files: [
                    {
                        expand: true,
                        cwd: 'build/src/js/',
                        src: ['**/*.js', '!**/*.min.js'],
                        dest: 'build/src/js/'
                    },
                    {
                        expand: true,
                        cwd: 'build/src/libs/',
                        src: ['*.js', '!*.min.js'],
                        dest: 'build/src/libs/'
                    },
                    {
                        expand: true,
                        cwd: 'build/src/libs/jquery',
                        src: ['*.js', '!*.min.js'],
                        dest: 'build/src/libs/jquery'
                    },
                    {
                        expand: true,
                        cwd: 'build/src/libs/layui/lay/modules',
                        src: ['*.js', '!*.min.js'],
                        dest: 'build/src/libs/layui/lay/modules'
                    },
                    {
                        expand: true,
                        cwd: 'build/src/libs/layui',
                        src: ['*.js', '!*.min.js'],
                        dest: 'build/src/libs/layui'
                    }
                ]
            }
        },
        cssmin: {
            /* minify: {
                expand: true,
                cwd: 'demo/resources/css',
                src: ['*.css', '!*.min.css'],
                dest: 'build/resources/css'
            } */
            main: {
                files: [
                    {
                        expand: true,
                        cwd: 'build/src/css',
                        src: ['*.css', '!*.min.css'],
                        dest: 'build/src/css'
                    }
                ]
            },
            easyUI: {
                files: [
                    {
                        expand: true,
                        cwd: 'build/src/libs/easyUI',
                        src: ['*.css', '!*.min.css'],
                        dest: 'build/src/libs/easyUI'
                    }
                ]
            }
        },
        watch: {
            options: {
                livereload: true
            },
            build: {
                files: ['www/src/*.html', 'www/src/js/**/*.js', 'www/src/css/*.css'],
                tasks: ['uglify', 'cssmin:main'],
                options: {
                    spawn: false
                }
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-clean');
    grunt.loadNpmTasks('grunt-contrib-concat');
    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-watch');

    grunt.registerTask('default', ['clean', 'copy', 'cssmin', 'watch']);
    //grunt.registerTask('default', ['clean']);
};