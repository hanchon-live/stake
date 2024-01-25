#!/bin/bash



# Template to go file
echo "creating the go code using the templates"
templ generate

# Clear old css files
echo "cleaning up old css files"
rm ./public/assets/*

# Compile tailwind and update the css dependecy to clear browser cache
echo "building new css using tailwind"
random=$RANDOM
node ./node_modules/tailwindcss/lib/cli/index.js -i ./assets/app.css -o ./public/assets/app-$random.css
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/app.css/app-$random.css/g" src/components/layout_templ.go
else
    sed -i "s/app.css/app-$random.css/g" src/components/layout_templ.go
fi

# Build the server
echo "building the server"
go build -o ./tmp/main .
