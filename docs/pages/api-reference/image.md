---
title: Image Render Options
layout: page
permalink: /api-reference/image/
---

# {{ page.title }}

## Target

| Key                           | Type              | Mapped Flag             |
|-------------------------------|-------------------|-------------------------|
| `cache_dir`                   | `string`          | `--cache-dir` |
| `cookie`                      | `[]object` [1]    | `--cookie` |
| `crop_h`                      | `int`             | `--crop-h` |
| `crop_w`                      | `int`             | `--crop-w` |
| `crop_x`                      | `int`             | `--crop-x` |
| `crop_y`                      | `int`             | `--crop-y` |
| `custom_header`               | `[]object` [2]    | `--custom-header` |
| `custom_header_propagation`   | `bool`            | `--custom-header-propagation`/`--no-custom-header-propagation` |
| `debug_javascript`            | `bool`            | `--debug-javascript`/`--no-debug-javascript` |
| `encoding`                    | `string`          | `--encoding` |
| `format`                      | `string`          | `--format` |
| `height`                      | `int`             | `--height` |
| `images`                      | `bool`            | `--images`/`--no-images` |
| `javascript`                  | `bool`            | `--disable-javascript`/`--enable-javascript` |
| `javascript_delay`            | `int`             | `--javascript-delay` |
| `load_error_handling`         | `string`          | `--load-error-handling` |
| `load_media_error_handling`   | `string`          | `--load-media-error-handling` |
| `minimum_font_size`           | `int`             | `--minimum-font-size` |
| `password`                    | `string`          | `--password` |
| `quality`                     | `int`             | `--quality` |
| `smart_width`                 | `bool`            | `--disable-smart-width`/`--enable-smart-width` |
| `stop_slow_scripts`           | `bool`            | `--stop-slow-scripts`/`--no-stop-slow-scripts` |
| `transparent`                 | `bool`            | `--transparent` |
| `use_xserver`                 | `bool`            | `--use-xserver` |
| `username`                    | `string`          | `--username` |
| `width`                       | `int`             | `--width` |
| `zoom`                        | `float`           | `--zoom` |

[1,2] - The `[]object` type means that it's an array of object. In this case,
object has name and value attributes which are both strings i.e. is `{ name:
string, value: string}`.
