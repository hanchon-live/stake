#!/bin/bash

# Template to go file
echo "creating the go code using the templates"
templ generate

# Clear old css files
echo "cleaning up old css files"
rm ./public/assets/css/*

# Compile tailwind and update the css dependency to clear browser cache
echo "building new css using tailwind"
random=$RANDOM
node ./node_modules/tailwindcss/lib/cli/index.js -i ./assets/app.css -o ./public/assets/css/app-$random.css --minify
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/app.css/app-$random.css/g" src/components/layout_templ.go
else
    sed -i "s/app.css/app-$random.css/g" src/components/layout_templ.go
fi

# Build and update typescript dependency
echo "cleaning up old js files"
rm ./public/assets/js/*
# rm ./public/tsconfig.build.tsbuildinfo
echo "building js files"
node ./node_modules/microbundle/dist/cli.js -i src/typescript/index.ts -o dist/bundle.js --no-pkg-main -f umd
# node ./node_modules/typescript/bin/tsc --build tsconfig.build.json
mv ./dist/bundle.js ./public/assets/js/index-$random.js
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/index.js/index-$random.js/g" src/components/layout_templ.go
else
    sed -i "s/index.js/index-$random.js/g" src/components/layout_templ.go
fi

# Build the server
echo "building the server"
go build -o ./tmp/main ./cmd/
