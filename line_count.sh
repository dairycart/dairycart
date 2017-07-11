command -v golocc >/dev/null 2>&1 || { go get github.com/warmans/golocc; }
(cd api; golocc --no-vendor ./...)