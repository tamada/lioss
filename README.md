[![GitHub Action Build](https://github.com/tamada/lioss/workflows/build/badge.svg?branch=master)](https://github.com/tamada/lioss/actions?workflow=build)
[![Coverage Status](https://coveralls.io/repos/github/tamada/lioss/badge.svg?branch=master)](https://coveralls.io/github/tamada/lioss?branch=master)
[![codebeat badge](https://codebeat.co/badges/dc3481f5-852b-4537-a5f5-150e2bfa998c)](https://codebeat.co/projects/github-com-tamada-lioss-master)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamada/lioss)](https://goreportcard.com/report/github.com/tamada/lioss)
[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/lioss/blob/master/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/lioss/releases/tag/v1.0.0)

# lioss

License Identification tool for OSS project.

## Description

OSSのプロジェクトにはライセンスが設定されている．
そのライセンスには，行って良いこと，行っては行けないことが規定されている．
そのライセンスに従って，開発者は開発を進めていく必要がある．
一方で，OSSのライセンスのコンフリクトが問題になる場合がある．
しかし，このコンフリクトを発見することは容易ではない．
ライセンスごとのコンフリクトの検出が難しいこともあるものの，機械的にライセンスを特定することも困難であるためである．

SPDX が機械的なライセンスの特定に向けての整備を行っているが，全てのプロジェクトに浸透しているわけではない．
そのため本プロジェクトでは，ライセンスファイルから，どのようなOSSのライセンスを定めているかを検出する．
そして，依存ライブラリのライセンスファイルからもOSSライセンスを特定することで，コンフリクトの検出を目指す．

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
    project directories, archive files contains LICENSE file, and/or LICENSE file.
```

### `mkliossdb`

Creates the database of lioss.
If the database was exists, this command updates the database.

```
mkliossdb [OPTIONS] <LICENSE...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
    -f, --format <FORMAT>    specifies format. Default is 'json'
    -h, --help               print this message.
LICENSE
    specifies license files.
```

## Install

### Go-lang

```
$ go get github.com/tamada/lioss
```

### Homebrew

```
$ brew tap tamada/brew
$ brew install lioss
```

## References

* [dmgerman/ninka](https://github.com/dmgerman/ninka)
    * Daniel M. German, Yuki Manabe and Katsuro Inoue. A sentence-matching method for automatic license identification of source code files. In 25nd IEEE/ACM International Conference on Automated Software Engineering (ASE 2010).
