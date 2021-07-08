#!/bin/bash

echo "deploy product $1 specification $2 at "http://localhost:8000/api/v1/deployment/$1/$2
return_code=$(curl -sw '%{http_code}' $CURL_PROXY -X POST -H "Content-Type: application/json" http://localhost:8000/api/v1/deployment/$1/$2)
echo " Deployment operation HTTP status: "$return_code