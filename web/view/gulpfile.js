  var gulp = require('gulp'),
  typescript = require("gulp-typescript"),
  sass = require('gulp-sass'),
  watch = require("gulp-watch"),
  livereload = require('gulp-livereload'),
  minify_css = require('gulp-clean-css'),
  minify_js = require('gulp-minify'),
  player = require('node-wav-player');




  gulp.task('sass', function() {
    gulp.src('style/**/main.scss')
      .pipe(sass()) //Changing from SCSS to CSS
      .pipe(minify_css({compatibility: 'ie8'})) //Minify CSS
      .pipe(gulp.dest('./style'))
  });



  gulp.task('typescript', function () {
    let isSuccess = true;

     return gulp.src('src/**/script.ts')
         .pipe(typescript({
              noImplicitAny: true,
              noEmitOnError: true,
              outFile: 'script.js'
         })).once("error", function () {
           isSuccess = false;
           die();
         }).once("end", function () {
           //if(isSuccess) kick();
         })

         .pipe(minify_js({
             ext:{
                 src:'.js',
                 min:'.min.js'
             } /* ,
            exclude: ['tasks'], Those files won't minify
           ignoreFiles: ['.combo.js', '-min.js']
             Won't minify files which matches pattern. */
         }))
          .pipe(gulp.dest('./src'))
  });


gulp.task('watch', function(done) { 
        gulp.watch('src/**/*.ts', gulp.series('typescript'));
        gulp.watch('style/**/*.scss', gulp.series('sass'));
        done();
       });
 gulp.task('default', gulp.series('watch'));

      
       function kick(){
         player.play({
           path: './gulp_special_effects/smb_kick.wav',
         })

         setTimeout(() => {
           player.stop();
         }, 200);
       }

       function die(){
         player.play({
           path: './gulp_special_effects/smbmariodie.wav',
         })

         setTimeout(() => {
           player.stop();
         }, 600);
       }

/*
Jeśli nie zadziała, pobrać pliki zawarte w "package.json" za pomocą
komendy "npm install" albo "sudo npm install"

Bootstrap jest w folderze "./node_modules/bootstrap/scss/bootstrap.scss".
Tam są podpięte wszystkie pliki. Większość przydatnych zależności znajduje
sie w pliku "variables.scss". Go (jeśli chcemy coś edytować)
należy SKOPIOWAĆ do innej ścieżki i dopiero wtedy edytować.
*/
