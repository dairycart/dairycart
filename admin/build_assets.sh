sass src/sass/app.sass assets/css/app.css
cd src/elm
elm-format --yes *.elm
elm-make --yes Main.elm --output ../../assets/js/elm.js --warn --debug
cd ../../