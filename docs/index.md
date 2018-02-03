---
title: Get
layout: default
---

## Installation

Sanaa comes in a single binary. All you need to do is download the binary [from
the releases page][releases] to any location in your `$PATH` and you're good to
go. Shasum-256 files are provided if you want to validate your download.

Also note that `wkhtmltoimage` and `wkhtmltopdf` are expected to be available in
your `$PATH` for sanaa to be able to autodetect them.

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

[King'ori J. Maina][personal-site] © 2018. Under the [GNU General Public License
v3.0 bundled therein][license], you may copy, distribute and modify the software
as long as you track changes/dates in source files. Any modifications to or
software including (via compiler) GPL-licensed code must also be made available
under the GPL along with build & install instructions.

[contributing]: https://raw.githubusercontent.com/itskingori/sanaa/master/LICENSE
[github-page]: https://pages.github.com/
[glide]: https://github.com/Masterminds/glide
[jekyll]: http://jekyllrb.com/
[personal-site]: http://kingori.co/
[license]: https://raw.githubusercontent.com/itskingori/sanaa/master/LICENSE
[releases]: https://github.com/itskingori/sanaa/releases
