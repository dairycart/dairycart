if [ ! -d dist/css ]; then
  mkdir -p dist/css;
fi

if [ ! -d dist/js ]; then
  mkdir -p dist/js;
fi

rm dist/css/*.css*
rm dist/js/*.js
tsc src/typescript/*.ts --outFile dist/js/app.js
sass src/sass/*.sass dist/css/app.css