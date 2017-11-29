import gulp from "gulp";
import jsonminify from "gulp-jsonminify";
import less from "gulp-less";
import minifyCSS from "gulp-minify-css";
import rename from "gulp-rename";
import uglify from "gulp-uglify";
import watch from "gulp-watch";
import browserify from "browserify";
import babelify from "babelify";
import source from "vinyl-source-stream";
import buffer from "vinyl-buffer";
import del from "del";
import runSequence from "run-sequence";

const sourceDir = "src";
const buildDir = "dist";

gulp.task("clean", done => {
  del(buildDir).then(() => done()).catch(err => done(err));
});

gulp.task("build-content-script-js", () => {
  return browserify(`./${sourceDir}/content-script/index.js`)
    .transform(babelify)
    .bundle()
    .pipe(source("content-script.js")) // Convert from Browserify stream to vinyl stream.
    .pipe(buffer()) // Convert from streaming mode to buffered mode.
    .pipe(uglify({ mangle: false }))
    .pipe(gulp.dest(buildDir));
});

gulp.task("build-background-script-js", () => {
  return browserify(`./${sourceDir}/content-script/background.js`)
    .transform(babelify)
    .bundle()
    .pipe(source("background.js")) // Convert from Browserify stream to vinyl stream.
    .pipe(buffer()) // Convert from streaming mode to buffered mode.
    .pipe(uglify({ mangle: false }))
    .pipe(gulp.dest(buildDir));
});

gulp.task("build-content-script-css", () => {
  return gulp
    .src(`${sourceDir}/content-script/index.css`)
    .pipe(less())
    .pipe(minifyCSS())
    .pipe(rename("content-script.css"))
    .pipe(gulp.dest(buildDir));
});

gulp.task("build-manifest", () => {
  return gulp
    .src(`${sourceDir}/manifest.json`)
    .pipe(jsonminify())
    .pipe(gulp.dest(buildDir));
});

gulp.task("build", ["clean"], done => {
  runSequence(["build-content-script-js", "build-content-script-css", "build-background-script-js", "build-manifest"], done);
});

gulp.task("watch", ["build"], () => {
  return watch(`${sourceDir}/**/*`, () => { runSequence("build") });
});

gulp.task("default", ["build"]);
