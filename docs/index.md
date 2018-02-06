---
title: Get
layout: default
---

## Installation

Sanaa is a single Go binary. All you need to do is download the binary for your
platform [from the releases page][releases] to any location in your `$PATH` and
you're good to go.

If using Docker, there's the `kingori/sanaa` [image on Docker Hub][dockerhub].
For examples and more information, checkout out [the docker image's
repository][dockerrepo].

## Dependencies

Just make sure that `wkhtmltoimage` and `wkhtmltopdf`  are available in your
`$PATH` for sanaa to be able to autodetect them. Get [downloads
here][wkhtmltopdf].

## Usage

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

1. `400 Bad Request` - unable to unmarshall the request JSON, or you've
   requested for a render type apart from the supported types i.e. `image` or
   `pdf`.
2. `500 Internal Server Error` - unable to enqueue the job for the workers to
   pick up e.g. if redis is down.

For each case, the body of the response should include a message explaining the
reason for failure. For example, if you send a bad JSON body during an image
render request, the response would be something like this:

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


## Development

Below instructions are only necessary if you intend to work on the source code
(find [contributing guidelines][contributing] here). For normal usage the above
instructions should do.

### Building

1. Fetch the code with `go get github.com/itskingori/sanaa`.
1. Install the Go development tools via `make dependencies`.
1. Install application dependencies via `make install` (they'll be placed in
   `./vendor`). Requires [Glide package manager][glide].
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

## License

[King'ori J. Maina][personal-site] © 2018. The [GNU Affero General Public
License v3.0 bundled therein][license], essentially says, if you make a
derivative work of this, and distribute it to others under certain
circumstances, then you have to provide the source code under this license. And
this still applies if you run the modified program on a server and let other
users communicate with it there.

[contributing]: https://github.com/itskingori/sanaa/blob/master/CONTRIBUTING.md
[dockerhub]: https://hub.docker.com/r/kingori/sanaa
[dockerrepo]: https://github.com/itskingori/docker-sanaa
[github-page]: https://pages.github.com/
[glide]: https://github.com/Masterminds/glide
[jekyll]: http://jekyllrb.com/
[personal-site]: http://kingori.co/
[license]: https://raw.githubusercontent.com/itskingori/sanaa/master/LICENSE
[releases]: https://github.com/itskingori/sanaa/releases
[wkhtmltopdf]: https://wkhtmltopdf.org/downloads.html
