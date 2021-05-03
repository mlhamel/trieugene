# Trieugene

Data aggregation pipeline properly scrape and collect interesting dataset

**Project details**: [Github](https://github.com/mlhamel/trieugene)

**Deployment**: [DigitalOcean Apps](https://cloud.digitalocean.com/apps/db2692b5-3f79-46cc-9e64-44adeb1ded01?i=0e8e19)

**Internal Dashboard**: [Faktory](https://trieugene-njod9.ondigitalocean.app/)

**Includes Services**: [Faktory](services/faktory/) •  [Rougecombien](services/rougecombien/) • [Rotondo](services/rotondo/)

## Installation

For local development you simply have to follow those steps:

Get the code:

```shell
$ git clone https://github.com/mlhamel/trieugene.git`
```

Change current working dir:

```shell
$ cd trieugene
```

Run dev setup:

```shell
$ dev up
```

Build from source:

```shell
$ dev build
```

And you're done!

## Usage

Locally you're gonna need to run multiple processes in parallele to have trieugene and it's different services up and running:

Change current working dir:

```shell
$ cd trieugene
```

Run trieugene:

```shell
$ bin/trieugene dev
```

Then you can start running the different services like the following:

```shell
$ bin/rougecombien dev
```

or simply running the main job scheduler:

```shell
$ bin/rotondo dev
```

## License

[GPL-3.0-or-later](COPYING)

## Authors


Crafted by [Mathieu Leduc-Hamel](mailto:info@mlhamel.org).