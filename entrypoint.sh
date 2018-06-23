#!/bin/sh
set -euo pipefail

SECRET_FILE="simply-do-firebase-adminsdk.json"
echo ${FIREBASE_SERVICE_ACCOUNT_FILE_BASE64} | base64 -d > ${SECRET_FILE}

"$@"
