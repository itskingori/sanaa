# Changelog

## 0.9.0

* Change all logging output to start with lowercase letters.

## 0.8.0

* Configure commands to not accept arguments i.e. `server`, `worker` and
  `version`.
* Improve `version` command to include build SHA.

## 0.7.0

* Remove explicit region configuration, no need for `--s3-region` flag since
  it's picked from AWS configuration.

## 0.6.0

* Add worker failure `--max-retries` configuration option.
* Upgrade `github.com/itskingori/go-wkhtml` to v1.0.0.

## 0.5.0

* Add validation for worker `--s3-bucket` flag.
* Add configuration of region of the S3 bucket via a `--s3-region` flag on the
  worker.
* Fix issue with server component not returning the appropriate response if
  unable to enqueue a job i.e. if redis is down. Previously it would return a
  `201 Created` instead of a `500 Internal Server Error`.
* Fix issue where some methods would not return an error correctly which could
  possibly affect other components that rely on the returned error to apply the
  right logic.
* Fix issue with the `/status` endpoint where it would return with a `200 OK`
  even if there was an issue generating the pre-signed URL to the rendered file.

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
