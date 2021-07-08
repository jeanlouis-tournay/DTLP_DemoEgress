#!/bin/bash

push_data() {
  echo "call product operator to push $1 at "http://localhost:8000/api/v1/$2
  return_code=$(curl -sw '%{http_code}' $CURL_PROXY -X POST -H "Content-Type: application/json" -d @$1 http://localhost:8000/api/v1/$2)
  echo " Push platform operation HTTP status: "$return_code
  if [ "$return_code" != "202" ]; then
    echo "unable to push descriptors"
    exit 1
  fi
  echo "Operation succeeded"
}


push_data "platform-descriptor.json"   "model/platform/descriptors"
push_data "platform-servicerail-api.json" "model/platform/platform-services"
push_data "contract.json" "model/contract/digital-products"
push_data "deploy.json" "model/contract/specification"


