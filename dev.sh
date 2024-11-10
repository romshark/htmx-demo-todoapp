#!/bin/sh

if [ ! -d "node_modules" ]; then
  echo "node_modules not found. Installing dependencies..."
  npm i
fi

echo "\n--- Build JavaScript"
npm run build:js

echo "\n--- Build CSS"
npm run build:css

echo "\n--- Generate Templ Templates"
templ generate

go run github.com/romshark/templier
