#!/bin/sh

# Inject ENV vars into non-bundled config file
echo "Injecting ENV vars into non-bundled JS config..."
envsubst < /usr/share/nginx/html/conf.js.template > /usr/share/nginx/html/conf.js

# Remove any JS comment (only // supported, not /*...*/) lines from this un-
# bundled JS. Use temporary file since we can't stdin and stdout WRT same file
grep -v '^\s*//' /usr/share/nginx/html/conf.js > /usr/share/nginx/html/conf.tmp.js
mv /usr/share/nginx/html/conf.tmp.js /usr/share/nginx/html/conf.js

# Running server
echo "Starting nginx..."
nginx -g 'daemon off;'