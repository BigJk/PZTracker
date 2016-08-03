var exec = require('child_process').exec;

var gulp = require('gulp');

var minify = require('gulp-minify');
var zip = require('gulp-zip');

// Building Release
gulp.task('default', function () {
  // Minifiey the inject.js
  gulp.src('inject.js')
    .pipe(minify({
      noSource: true
    }))
    .pipe(gulp.dest('release'));
  
  // Builds PZTracker
  exec('go build', function (err, stdout, stderr) {
    gulp.src(['PZTracker', 'PZTracker.exe'])
      .pipe(gulp.dest('release'));
  });

  // Copys the PZTrackerMod
  gulp.src(['PZTrackerMod/**/*'])
      .pipe(gulp.dest('release/PZTrackerMod'));
});

// Zip the Release
gulp.task('zip', function () {
  gulp.src('release/**/*')
    .pipe(zip('pztracker_release.zip'))
    .pipe(gulp.dest('release'));
});