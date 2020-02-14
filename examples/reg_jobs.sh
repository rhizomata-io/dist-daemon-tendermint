#!/bin/bash

curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data '{"interval":"10ms","greet":"hello 10ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/veryfast

curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data '{"interval":"200ms","greet":"hello 200ms" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi1

curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data '{"interval":"0.5s","greet":"hello 0.5s" }' http://localhost:7777/v1/daemon/job/add/factory/hello-worker/jobid/hi2
