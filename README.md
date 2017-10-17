# Sanaa

A HTML to PDF/Image conversion microservice powered by `wkhtmltopdf`,
`wkhtmltoimage` and `xvfb`.

## Running

This tool is self documenting, so running the following will get you started:

```console
$ sanaa help
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
