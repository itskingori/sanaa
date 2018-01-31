# Sanaa

[![Build Status](https://travis-ci.org/itskingori/sanaa.svg?branch=master)](https://travis-ci.org/itskingori/sanaa)

A HTML to PDF/Image conversion microservice powered by `wkhtmltopdf` and `wkhtmltoimage`.

## Usage

```console
# Start the server (that will receive requests)
$ sanaa server \
  --verbose

# Start the workers (that will process requests)
$ sanaa worker \
  --s3-bucket="my-bucket" \
  --s3-region="us-east-1" --verbose
```

### Rendering Images

Make a `POST` request to `/render/image`.

```http
POST /render/image HTTP/1.1
Content-Type: application/json
Host: 127.0.0.1:8080

{
    "target": {
        "format": "png",
        "height": 480,
        "weight": 640
    },
    "source": {
        "url": "https://google.com"
    }
}
```

### Rendering PDFs

Make a `POST` request to `/render/pdf`.

```http
POST /render/pdf HTTP/1.1
Content-Type: application/json
Host: 127.0.0.1:8080

{
    "target": {
        "margin_top": 10,
        "margin_bottom": 10,
        "margin_left": 10,
        "margin_right": 10,
        "page_height": 210,
        "page_width": 300
    },
    "source": {
        "url": "https://google.com"
    }
}
```

## Development

Below instructions are only necessary if you intend to work on the source code.
For normal usage the above installation instruction should do.

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

[glide]: https://github.com/Masterminds/glide
[releases]: https://github.com/itskingori/sanaa/releases
