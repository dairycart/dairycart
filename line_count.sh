command -v golocc >/dev/null 2>&1 || { go get github.com/warmans/golocc; }
(cd .; golocc --no-vendor ./...)