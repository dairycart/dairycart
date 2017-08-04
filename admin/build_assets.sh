rm dist/css/*.css*
rm dist/js/*.js
tsc src/typescript/*.ts --outFile dist/js/app.js
sass src/sass/*.sass dist/css/app.css
cp node_modules/bulma/css/bulma.css dist/css/bulma.css