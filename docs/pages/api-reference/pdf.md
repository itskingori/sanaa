---
title: PDF Render Options
layout: page
permalink: /api-reference/pdf/
---

# {{ page.title }}

------------------

## Target

| Key                             | Type               | Mapped Flag             |
|---------------------------------|--------------------|-------------------------|
| `cache_dir`                     | `string`           | `--cache-dir` |
| `cookie`                        | `[]object` [1]     | `--cookie` |
| `custom_header`                 | `[]object` [2]     | `--custom-header` |
| `custom_header_propagation`     | `bool`             | `--custom-header-propagation`/`--no-custom-header-propagation` |
| `debug_javascript`              | `bool`             | `--debug-javascript`/`--no-debug-javascript` |
| `dpi`                           | `int`              | `--dpi` |
| `encoding`                      | `string`           | `--encoding` |
| `external_links`                | `bool`             | `--disable-external-links`/`--enable-external-links` |
| `forms`                         | `bool`             | `--disable-forms`/`--enable-forms` |
| `grayscale`                     | `bool`             | `--grayscale` |
| `images`                        | `bool`             | `--images`/`--no-images` |
| `image_dpi`                     | `int`              | `--image-dpi` |
| `image_quality`                 | `int`              | `--image-quality` |
| `internal_links`                | `bool`             | `--disable-internal-links`/`--enable-internal-links` |
| `javascript`                    | `bool`             | `--disable-javascript`/`--enable-javascript` |
| `javascript_delay`              | `int`              | `--javascript-delay` |
| `load_error_handling`           | `string`           | `--load-error-handling` |
| `load_media_error_handling`     | `string`           | `--load-media-error-handling` |
| `lowquality`                    | `bool`             | `--lowquality` |
| `margin_bottom`                 | `int`              | `--margin-bottom` |
| `margin_left`                   | `int`              | `--margin-left` |
| `margin_right`                  | `int`              | `--margin-right` |
| `margin_top`                    | `int`              | `--margin-top` |
| `minimum_font_size`             | `int`              | `--minimum-font-size` |
| `orientation`                   | `bool`             | `--orientation` |
| `page_height`                   | `string`           | `--page-height` |
| `page_size`                     | `int`              | `--page-size` |
| `page_width`                    | `string`           | `--page-width` |
| `no_pdf_compression`            | `int`              | `--no-pdf-compression` |
| `password`                      | `string`           | `--password` |
| `smart_width`                   | `bool`             | `--disable-smart-shrinking`/`--enable-smart-shrinking` |
| `stop_slow_scripts`             | `bool`             | `--stop-slow-scripts`/`--no-stop-slow-scripts` |
| `title`                         | `string`           | `--title` |
| `use_xserver`                   | `bool`             | `--use-xserver` |
| `username`                      | `string`           | `--username` |
| `zoom`                          | `float`            | `--zoom` |

[1,2] - The `[]object` type means that it's an array of object. In this case,
object has name and value attributes which are both strings i.e. is `{ name:
string, value: string}`.
