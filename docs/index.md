---
title: Get
layout: default
---

## Synopsis 🔍

The objective of the project is to provide a HTTP API around `wkhtmltoimage` and
`wkhtmltopdf`. There's been no attempt to modify those two binaries. Sanaa
pretty much translates options passed in as JSON to flags, runs the command,
fetches the result and translates that results back to JSON which is served back
as a response.

The architecture of the project takes a server/worker architecture. This was
deemed ideal as it works well with scaling. You can scale the server part based
on incoming requests and the worker part based on jobs on the queue.

## Installation ⬇️

Sanaa is a single Go binary. All you need to do is download the binary for your
platform [from the releases page][releases] to any location in your `$PATH` and
you're good to go.

If using Docker 🐳, there's the `kingori/sanaa` image [on Docker Hub][dockerhub].
Check out the `examples/` folder for a docker-compose config.

## Dependencies 🖇️

Just make sure that `wkhtmltoimage` and `wkhtmltopdf`  are available in your
`$PATH` for sanaa to be able to autodetect them. Get [downloads
here][wkhtmltopdf].

## Configuration 🎛️

Most configuration is done via flags, see `sanaa --help`. But there's AWS
specific configuration, which are mostly secrets, that don't seem appropriate to
set via flags.

For example, Sanaa requires AWS credentials with permissions that give it access
to the S3 bucket it will use to store the results of rendering. These
credentials will be sourced automatically from the following locations (in order
of priority, first at the top):

1. Environment Credentials - via environment variables:

   ```bash
   export AWS_ACCESS_KEY_ID=SOME_KEY
   export AWS_SECRET_ACCESS_KEY=SOME_SECRET
   export AWS_REGION=us-east-1
   ```

2. Shared Credentials file - via `~/.aws/credentials`:

   ```text
   [default]
   aws_access_key_id = <SOME_KEY>
   aws_secret_access_key = <SOME_SECRET>
   aws_region = us-east-1
   ```

3. EC2 Instance Role Credentials - assigns credentials to application if it's
   running on an EC2 instance that's been given an EC2 Instance Role. This
   removes the need to manage credential files in production.

## Usage 💻

### Getting Started

Start the server (that will receive requests):

```console
$ sanaa server --verbose
INFO[0000] Starting the server
INFO[0000] Request TTL set to 86400 seconds
INFO[0000] Listening on http://0.0.0.0:8080
```

Start the worker (that will process requests):

```console
$ sanaa worker --s3-bucket=example-bucket-name --verbose
INFO[0000] Starting the worker
INFO[0000] Using wkhtmltoimage 0.12.4 (with patched qt)
INFO[0001] Using wkhtmltopdf 0.12.4 (with patched qt)
INFO[0001] Concurrency set to 2
INFO[0001] Registering 'convert' queue
INFO[0001] Waiting to pick up jobs placed on any registered queue
```

### Basic Usage

#### Rendering Images & PDFs

For images, make a `POST` request to `/render/image`:

```http
POST /render/image HTTP/1.1
Content-Type: application/json
Host: 127.0.0.1:8080
Connection: close
Content-Length: 172

{
    "target": {
        "format": "png",
        "height": 1080,
        "width": 1920
    },
    "source": {
        "url": "https://en.wikipedia.org/wiki/Kenya"
    }
}
```

For PDFs, make a `POST` request to `/render/pdf`:

```http
POST /render/pdf HTTP/1.1
Content-Type: application/json
Host: 127.0.0.1:8080
Connection: close
Content-Length: 127

{
    "target": {
        "page_size": "A4"
    },
    "source": {
        "url": "https://en.wikipedia.org/wiki/Kenya"
    }
}
```

If a render request was successful, expect a `201 Created` HTTP response
indicating that the server has acknowledged the request:

```http
HTTP/1.1 201 Created
Content-Type: application/json
Date: Tue, 06 Feb 2018 05:19:09 GMT
Content-Length: 176
Connection: close

{
  "uuid": "640882bd-9441-48fb-8686-27286f399004",
  "created_at": "2018-02-06T05:19:09Z",
  "started_at": "",
  "ended_at": "",
  "expires_in": 86400,
  "file_url": "",
  "status": "pending",
  "logs": ""
}
```

In case of failure, expect an appropriate response as well. For example:

1. `400 Bad Request` - if unable to unmarshall the request JSON, or if you've
   requested for a render type apart from the supported types i.e. `image` or
   `pdf`.
2. `500 Internal Server Error` - if unable to enqueue the job for the workers to
   pick up e.g. if redis is down.

#### Checking Render Request Status

Each render request that has been enqueued is assigned a UUID (found in `uuid`
attribute of the response to a render request). Pass the UUID to the
`/status/{uuid}` endpoint via `GET` to get an update on the status of the
conversion job:

```http
GET /status/4c815816-1bfe-4790-b8d1-ee06c98b7d6d HTTP/1.1
Content-Type: application/json
Host: 127.0.0.1:8080
Connection: close

```

The status endpoint will return a `200 OK` HTTP response with details of the
conversion job. An example of one that's still in processing:

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 06 Feb 2018 06:33:07 GMT
Connection: close
Transfer-Encoding: chunked

{
  "uuid": "4c815816-1bfe-4790-b8d1-ee06c98b7d6d",
  "created_at": "2018-02-06T06:32:42Z",
  "started_at": "2018-02-06T06:32:45Z",
  "ended_at": "",
  "expires_in": 86400,
  "file_url": "",
  "status": "processing",
  "logs": "Loading page (1/2)\n[>                                                           ] 0%\r[======>                                                     ] 10%\r[==========>                                                 ] 17%\r[============>                                               ] 21%\r[=============>                                              ] 23%\r[===============>                                            ] 26%\r[=================>                                          ] 29%\r[=================>                                          ] 29%\r[=================>                                          ] 29%\r[==================>                                         ] 31%\r[====================>                                       ] 34%\r[======================>                                     ] 37%\r[=======================>                                    ] 39%\r[==========================>                                 ] 44%\r[============================>                               ] 47%\r[==============================>                             ] 50%\r[===============================>                            ] 52%\r[================================>                           ] 54%\r[=================================>                          ] 56%\r[==================================>                         "
}
```

Notably, several fields reflect the `processing` state of the conversion job:

1. `status` is set to `processing`,
2. `started_at` has been set to the time the processing started,
3. `ended_at` is still empty (obviously) and
4. `logs` has been populated with some output from the processing.

In case of failure, expect to recieve responses that communicate the problem.
For example:

1. `404 Not Found` - may happen if you render request has expired (based on TTL)
   or if there's no job found matching the UUID set.
2. `400 Bad Request` - if your identifier is not a valid UUID.
3. `500 Internal Server Error` - if the server is unable to fulfill your request
   i.e. if redis is down.

#### Attributes Of Response Objects

The `/render/{type}` and `/status/{uuid}` endpoints either return an object
representing an error or a conversion job.

For errors, the response body is simple and self-explanatory. It includes the
`uuid` of the request and a `message` explaining the error. For example, if you
send a bad JSON body during an image render request, the response would be
something like this:

```http
HTTP/1.1 500 Internal Server Error
Content-Type: application/json
Date: Tue, 06 Feb 2018 07:26:44 GMT
Content-Length: 91
Connection: close

{
  "uuid": "536d3847-64b8-497a-8d8a-ac541dfa9c9e",
  "message": "Unable to unmarshal json to image type"
}
```

For render requests, the returned object represents a conversion job which has
the following attributes:

| Attribute     |  Description |
|---------------|--------------|
| `uuid`        | Unique identifier of the request |
| `created_at`  | When the request was initiated |
| `started_at`  | When the request was picked by a worker for processing |
| `ended_at`    | When a worker completed processing the request after picking it up |
| `expires_in`  | How long to persist the request and any of it's data |
| `file_url`    | URL to fetch the artefact generated by the request after processing |
| `status`      | Status of the job i.e. `pending`, `processing`, `failed`, `succeeded` |
| `logs`        | Output of processing by the worker, useful when debugging |

### Advanced Usage

#### Health Endpoints

The server component has two health endpoints available:

* `/health/live` -  liveness endpoint, indicates that the server is up.
* `/health/ready` - readiness endpoint, indicates that server is ready to
  receive requests.

Pass the `?full=1` query parameter to expose the details of the check in the
JSON response. These are omitted by default for performance.

Also note that both endpoints return the appropriate response conveying the
health of the service. To demonstrate this, make a request to the readiness
endpoint:

```http
GET /health/ready?full=1 HTTP/1.1
Host: 127.0.0.1:8080
Connection: close
```

If redis is up, you should get a `200 OK` HTTP response:

```http
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 22 Feb 2018 19:54:41 GMT
Content-Length: 37
Connection: close

{
    "redis-tcp-connection": "OK"
}
```

If redis is down, you should get a `503 Service Unavailable` HTTP response:

```http
HTTP/1.1 503 Service Unavailable
Content-Type: application/json; charset=utf-8
Date: Thu, 22 Feb 2018 20:00:27 GMT
Content-Length: 87
Connection: close

{
    "redis-tcp-connection": "dial tcp 127.0.0.1:6379: connect: connection refused"
}
```

## Development ⚒️

Below instructions are only necessary if you intend to work on the source code
(find [contributing guidelines][contributing] here). For normal usage the above
instructions should do.

### Building

1. Fetch the code with `go get github.com/itskingori/sanaa`.
1. Install the Go development tools via `make dependencies`.
1. Install application dependencies via `make install` (they'll be placed in
   `./vendor`). Requires [golang/dep][dep] package manager.
1. Build and install the binary with `make build`.
1. Run the command e.g. `./sanaa help`.

### Testing

1. Install the Go testing tools via `make dependencies`.
1. Run linter using `make lint` and test using `make test`.

### Documentation

The home page is built using [Jekyll][jekyll] (a fun and easy to use static site
generator) and it is [hosted on GitHub Pages][github-page]. The code is in the
`docs/` folder if you want to have a peek.

### Releasing

1. Create a tag (`git tag`) and push the tags to remote (`git push --tags`).
2. CI pipeline will detect the tag and create a [GitHub release here][releases].
   To note:
   * Tags matching `x.y.z` will be marked as final releases.
   * Tags matching `x.y.z-*` will be marked as pre-releases.
   * Tags not matching either of the above, will be ignored and assumed to be
     normal tags.
   * Compressed binary with a shasum 256 file will be uploaded as attachments to
     the release.
3. Trigger will start a build on Docker Hub to publish two Docker images:
   `kingori/sanaa:latest` and `kingori/sanaa:x.y.z`.

## License 📜

[King'ori J. Maina][personal-site] © 2018. The [GNU Affero General Public
License v3.0 bundled therein][license], essentially says, if you make a
derivative work of this, and distribute it to others under certain
circumstances, then you have to provide the source code under this license. And
this still applies if you run the modified program on a server and let other
users communicate with it there.

[contributing]: https://github.com/itskingori/sanaa/blob/master/CONTRIBUTING.md
[dep]: https://golang.github.io/dep/
[dockerhub]: https://hub.docker.com/r/kingori/sanaa
[github-page]: https://pages.github.com/
[jekyll]: http://jekyllrb.com/
[personal-site]: http://kingori.co/
[license]: https://raw.githubusercontent.com/itskingori/sanaa/master/LICENSE
[releases]: https://github.com/itskingori/sanaa/releases
[wkhtmltopdf]: https://wkhtmltopdf.org/downloads.html
