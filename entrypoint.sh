#!/bin/sh
set -euo pipefail

SECRET_FILE="simply-do-firebase-adminsdk.json"
echo ${FIREBASE_SERVICE_ACCOUNT_FILE} > $SECRET_FILE

/app/server
