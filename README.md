# Sanaa

A HTML to PDF/Image conversion microservice powered by `wkhtmltopdf`,
`wkhtmltoimage` and `xvfb`.

## Usage

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
2. Install dependencies with `glide install --strip-vendor`, which will be
   placed in `./vendor`.
3. Build and install the binary with `go build` from within the repository.
4. Run the command e.g. `./sanaa help`.
