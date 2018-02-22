# Changelog

## 0.1.0

* First prototype of the idea with basic features i.e. server and worker
  components that use redis as a data-store. Server receives requests (to render
  image or pdf from source URL) and the worker processes them (generates
  requested document and uploads it to S3). See [project
  homepage](https://kingori.co/sanaa) for details.
