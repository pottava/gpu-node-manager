#!/bin/sh

echo 'target/\ntest-results/\ntmp/\nroutes/\nkey.json' > src/.gitignore

url=$( gcloud run services describe "$1" --region "$2" \
    --format 'value(status.address.url)' )
echo "window.apiBaseURL = \"${url}\";" > src/public/js/app/config.js
