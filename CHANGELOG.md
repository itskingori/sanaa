# Changelog

## 0.4.0

* Improve output of logs by presenting them as an array. Each log line will be
  an entry in the array and all the newlines are handled to improve the output.

## 0.3.0

* Configure sanaa to run as non-root user in Dockerfile.

## 0.2.0

* Add `/health/live` (liveness) and `/health/ready` (readiness) health endpoints
  on the server component.

## 0.1.0

* First prototype of the idea with basic features i.e. server and worker
  components that use redis as a data-store. Server receives requests (to render
  image or pdf from source URL) and the worker processes them (generates
  requested document and uploads it to S3). See [project
  homepage](https://kingori.co/sanaa) for details.
