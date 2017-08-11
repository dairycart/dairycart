rm dist/css/*.css*
rm dist/js/*.js
tsc src/typescript/*.ts --outFile dist/js/app.js
sass src/sass/*.sass dist/css/app.css