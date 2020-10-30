[![GitHub Action Build](https://github.com/tamada/lioss/workflows/build/badge.svg?branch=master)](https://github.com/tamada/lioss/actions?workflow=build)
[![Coverage Status](https://coveralls.io/repos/github/tamada/lioss/badge.svg?branch=master)](https://coveralls.io/github/tamada/lioss?branch=master)
[![codebeat badge](https://codebeat.co/badges/dc3481f5-852b-4537-a5f5-150e2bfa998c)](https://codebeat.co/projects/github-com-tamada-lioss-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamada/lioss)](https://goreportcard.com/report/github.com/tamada/lioss)
[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/lioss/blob/master/LICENSE)
[![Version](https://img.shields.io/badge/Version-0.9.0-yellowgreen.svg)](https://github.com/tamada/lioss/releases/tag/v1.0.0)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftamada%2Flioss.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftamada%2Flioss?ref=badge_shield)

# lioss

License Identification tool for OSS project.

## :speaking_head: Overview

Generally, OSS projects have licenses.
The licenses grant permissions to users for using, modifying, and sharing the software.
The users of the software must follow the terms shown in the licenses.

On the other hand, today's software generally has some dependencies.
Additionally, dependant software has some dependencies, too.
Therefore, the dependant graph of the OSS becomes complex.

In such a situation, it is a quite tough task for checking the conflicts among licenses.
The first problem is to detect a conflict between two given licenses.
The second problem is to identify the license of a project.
`lioss` tries to solve the above second problem by identifying the license of the given project.

SPDX is trying to automatically identify licenses, however,  it is hard to say that it became common sense.
This project detects the OSS licenses from the LICENSE files of the given projects.
Then, we aim to detect conflicts by identifying OSS licenses from the license files of dependent libraries.


## Usage

### `lioss`

Identifies license name from file and/or project directories.

```
lioss version 1.0.0
lioss [OPTIONS] <PROJECTS...>
OPTIONS
        --dbpath <DBPATH>          specifying database path.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
                                   Available values are: kgram, wordfreq, and tfidf.
    -t, --threshold <THRESHOLD>    specifies threshold of the similarities of license files.
                                   Each algorithm has default value. Default value is 0.75.
    -h, --help                     print this message.
PROJECTS
    project directories, archive files (jar, and zip) contains LICENSE file, and/or LICENSE file.
```


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftamada%2Flioss.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftamada%2Flioss?ref=badge_large)

### `mkliossdb`

Creates the database of lioss from License documents.

```
mkliossdb [OPTIONS] <LICENSE...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
    -h, --help               print this message.
LICENSE
    specifies license files.
```

## Install

### Go-lang

```
$ go get github.com/tamada/lioss
```

### :beer: Homebrew

```
$ brew tap tamada/brew
$ brew install lioss
```

### :muscle: Compiling yourself

```sh
$ git clone github.com/tamada/lioss
$ cd lioss
$ make
```

## References

* [dmgerman/ninka](https://github.com/dmgerman/ninka)
    * Daniel M. German, Yuki Manabe and Katsuro Inoue. A sentence-matching method for automatic license identification of source code files. In 25nd IEEE/ACM International Conference on Automated Software Engineering (ASE 2010).
    * This product identifies the license of each source file.
      However, it does not work on my environment.
* [pivotal/LicenseFinder](https://github.com/pivotal/LicenseFinder)
    * This product finds dependencies from build file, and find license.
* [SPDX](https://spdx.org) (Software Package Data Exchange).
    *