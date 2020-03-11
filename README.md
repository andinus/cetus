# Cetus

Cetus is a wallpaper setting tool written in Go. This is a work in-progress.

## Dependency

- [feh](https://feh.finalrewind.org/)

## Features

- Supports various Image of the Day services

There are several Image of the Day services on the web, cetus will pull the
latest image & set it as wallpaper.

Currently it supports:

- [Astronomy Picture of the Day](http://apod.nasa.gov/apod/astropix.html)

## Examples

### Image of the Day

#### Astronomy Picture of the Day

``` sh
# Pull image from Astronomy Picture of the Day
cetus -apod -apod-api=https://api.nasa.gov/planetary/apod \
      -apod-api-key=DEMO_KEY
```

### Set given wallpaper

``` sh
# Local image as wallpaper
cetus -img-path=/path/to/img

# Remote image as wallpaper
cetus -img-path=http://127.0.0.1

```
