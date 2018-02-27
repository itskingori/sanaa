---
title: Get
layout: default
---

## Synopsis 🔍

Sanaa provides a HTTP API around `wkhtmltoimage` and `wkhtmltopdf`. There's been
no attempt to modify those two binaries. It's [BYO][byo]-wkhtmltoX.

It translates options passed in as JSON to flags, runs the command, fetches the
command-output/generated-file and translates those results into a JSON response.
The generated file should have been uploaded to an S3 bucket (by that point) and
the API response should contain a signed link to it.

The current implementation assumes that Sanaa's _**role is purely to render what
you ask it to and provide you with a means to fetch it**_. This emphasis on a
single-responsibility made for a simple design. So, it's left up to you to use
the result as you wish.

That pretty much is it 💪 ... in a nutshell! 🥜🐚

## Features 🎉

* Server and worker components that are scalable separately. You can scale the
  server part based on incoming requests and the worker part based on jobs on
  the queue.
* Simple HTTP API with render request and status checking endpoints.
* Liveness and readiness endpoints for proper health checks.
* Cleans up after itself. Render requests (in redis) and their resulting files
  (in S3) expire after configurable TTL is exceeded.
* Configurable max retries on failure with built-in exponential backoff.
* Proper logging with unique id of each job on each line (where appropriate)
  makes it easy for filtering logs and therefore quick debugging.

## Installation ⬇️

Sanaa is a single Go binary. All you need to do is download the binary for your
platform [from the releases page][releases] to any location in your `$PATH` and
you're good to go.

If using Docker 🐳, there's the `kingori/sanaa` image [on Docker Hub][dockerhub].
Find [docker-compose][example1] and [kubernetes][example2] examples in the
`examples/` folder.

## Dependencies 🖇️

Just make sure that `wkhtmltoimage` and `wkhtmltopdf` binaries are available in
your `$PATH` for sanaa to be able to autodetect them. Fetch [downloads from
here][wkhtmltopdf].

## Configuration 🎛️

Most configuration is done via flags, see `sanaa --help`. But there's AWS
specific configuration, which are mostly secrets, that don't seem appropriate to
set via flags.

For example, Sanaa requires AWS credentials with permissions. The worker
requires upload access to the S3 bucket it will use to store the results of
rendering and the server will require access to generate signed URLs to download
from the same bucket.

These credentials will be sourced automatically from the following locations (in
order of priority, first at the top):

1. Environment Credentials - via environment variables:

   ```bash
   export AWS_ACCESS_KEY_ID=SOME_KEY
   export AWS_SECRET_ACCESS_KEY=SOME_SECRET
   ```

2. Shared Credentials file - via `~/.aws/credentials`:

   ```text
   [default]
   aws_access_key_id = <SOME_KEY>
   aws_secret_access_key = <SOME_SECRET>
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
INFO[0001] Maximum retries set to 1
INFO[0001] Registering 'convert' queue
INFO[0001] Waiting to pick up jobs placed on any registered queue
```

### Basic Usage

#### Rendering Images & PDFs

For images ([see reference][api-ref-image]), make a `POST` request to
`/render/image`:

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

For PDFs ([see reference][api-ref-pdf]), make a `POST` request to `/render/pdf`:

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
  "logs": [
    ""
  ]
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
conversion job. An example of one that's succeeded:

```http
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 Feb 2018 00:41:31 GMT
Connection: close
Transfer-Encoding: chunked

{
  "uuid": "21835d4a-5dfc-41a4-a798-21980baa43c9",
  "created_at": "2018-02-24T00:40:32Z",
  "started_at": "2018-02-24T00:40:36Z",
  "ended_at": "2018-02-24T00:40:57Z",
  "expires_in": 86400,
  "file_url": "https://s3.amazonaws.com/example-bucket-name/21835d4a-5dfc-41a4-a798-21980baa43c9/file.png?signed-url-signature",
  "status": "succeeded",
  "logs": [
    "Loading page (1/2)",
    "...",
    "Rendering (2/2)",
    "...",
    "Done"
  ]
}
```

Notably, several fields reflect the state of the conversion job:

1. `status` is set to `succeeded`,
2. `started_at` has been set to the time the processing started,
3. `ended_at` has been set to the time the processing ended. It should be empty
   if the it's still in processing state,
4. `logs` has been populated with output from the processing.

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

Timestamp fields are [RFC3339][rfc3339] and always in UTC.

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

For normal usage the above instructions should do. Below instructions are only
necessary if you intend to work on the source code (find [contributing
guidelines here][contributing], [plan here][plan] and the [milestones
here][milestones]).

### Building

1. Fetch the code with `go get -v github.com/itskingori/sanaa`.
1. Install the Go development tools via `make dependencies`.
1. Install application dependencies via `make install` (they'll be placed in
   `./vendor`). Requires [golang/dep][dep] package manager.
1. Build and install the binary with `make build`.
1. Run the command e.g. `./sanaa help` as a basic test.

### Testing

1. Install the Go testing tools via `make dependencies`.
1. Run linter using `make lint` and test using `make test`.

### Documentation

The home page is built using [Jekyll][jekyll] (a fun and easy to use static site
generator) and it is [hosted on GitHub Pages][github-page]. The code is in the
`docs/` folder if you want to have a peek.

### Releasing

1. Create a tag (`git tag`) and push the tags to remote (`git push --tags`).
2. CI pipeline (i.e. Travis CI) will detect the tag and create a [GitHub release
   here][releases]. To note:
   * Tags matching `x.y.z` will be marked as final releases.
   * Tags matching `x.y.z-*` will be marked as pre-releases.
   * Tags not matching either of the above, will be ignored and assumed to be
     normal tags.
   * Compressed binary with a shasum 256 file will be uploaded as attachments to
     the release.
3. Trigger will start a build on Docker Hub to publish two Docker images:
   `kingori/sanaa:latest` and `kingori/sanaa:x.y.z`.

## FAQ

**What does _Sanaa_ mean?**

It's the [Swahili][swahili] word for _"art"_ or more specifically a _"work of
beauty"_. I'm [Kenyan][kenya] so my bias to Swahili is obvious. 🤷

**How Can I Help?**

Write tests (or show me how to). I'm fairly new to Go and I'm of the opinion
that writing tests for a service like this is non-trivial. So far testing has
been manual but I plan to read on it and write some when I get time (as an
exercise in continuous learning).

Give feedback. Feel free to submit via [raising an issue][issue-new] or even
comment on [the open issues][issue-list].

## License 📜

[King'ori J. Maina][personal-site] © 2018. The [GNU Affero General Public
License v3.0 bundled therein][license], essentially says, if you make a
derivative work of this, and distribute it to others under certain
circumstances, then you have to provide the source code under this license. And
this still applies if you run the modified program on a server and let other
users communicate with it there.

[byo]: https://www.urbandictionary.com/define.php?term=BYO
[contributing]: https://github.com/itskingori/sanaa/blob/master/CONTRIBUTING.md
[dep]: https://golang.github.io/dep/
[dockerhub]: https://hub.docker.com/r/kingori/sanaa
[example1]: https://github.com/itskingori/sanaa/tree/master/examples/docker-compose
[example2]: https://github.com/itskingori/sanaa/tree/master/examples/kubernetes
[github-page]: https://pages.github.com/
[issue-list]: https://github.com/itskingori/sanaa/issues
[issue-new]: https://github.com/itskingori/sanaa/issues/new
[jekyll]: http://jekyllrb.com/
[kenya]: https://en.wikipedia.org/wiki/Kenya
[milestones]: https://github.com/itskingori/sanaa/milestones
[plan]: https://github.com/itskingori/sanaa/projects
[personal-site]: http://kingori.co/
[rfc3339]: https://www.ietf.org/rfc/rfc3339.txt
[swahili]: https://en.wikipedia.org/wiki/Swahili_language
[license]: https://raw.githubusercontent.com/itskingori/sanaa/master/LICENSE
[releases]: https://github.com/itskingori/sanaa/releases
[wkhtmltopdf]: https://wkhtmltopdf.org/downloads.html

[api-ref-image]: {{ site.baseurl }}/api-reference/image/
[api-ref-pdf]: {{ site.baseurl }}/api-reference/pdf/
