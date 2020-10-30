---
title: ":runner: Usage"
---

## `lioss`

`lioss` identifies license names of specified project directories, zip file and/or LICENSE files.

```sh
lioss version 1.0.0
lioss [OPTIONS] <PROJECTs...>
OPTIONS
        --database-path <PATH>     specifies the database path.
                                   If specifying this option, database-type option is ignored.
        --database-type <TYPE>     specifies the database type. Default is osi (enable multi options, separating by comma).
                                   Available values are: non-osi, osi, deprecated, osi-deprecated, and whole.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
                                   Available values are: kgram, wordfreq, and tfidf.
    -t, --threshold <THRESHOLD>    specifies threshold of the similarities of license files.
                                   Each algorithm has default value. Default value is 0.75.
    -h, --help                     prints this message.
PROJECTs
    LICENSE files, project directories, and/or archive files contains LICENSE file.
```

### Examples

```sh
$ lioss LICENSE    # run on the lioss directory.
LICENSE
$ lioss --database-type non-osi LICENSE    # run on the lioss directory.
LICENSE
	WTFPL (0.9804)
$ lioss --algorithm 9gram testdata   # run lioss for identifying project licenses in testdata directory.
osi -> OSI_APPROVED_DATABASE
testdata/project1/LICENSE
testdata/project2/license.txt
	GPL-3.0-only (0.9803)
	GPL-3.0-or-later (0.9803)
	AGPL-3.0-only (0.9654)
	AGPL-3.0-or-later (0.9654)
testdata/project3/license
	Apache-2.0 (1.0000)
	ECL-2.0 (0.9669)
testdata/project3/subproject/license
	MIT (0.9677)
	Xnet (0.8972)
	MIT-0 (0.8904)
$ lioss --algorithm 9gram --database-type osi,non-osi testdata   # run lioss for identifying project licenses in testdata directory.
testdata/project1/LICENSE
	WTFPL (0.9536)
testdata/project2/license.txt
	GPL-3.0-only (0.9803)
	GPL-3.0-or-later (0.9803)
	AGPL-3.0-only (0.9654)
	AGPL-3.0-or-later (0.9654)
	SSPL-1.0 (0.9234)
testdata/project3/license
	Apache-2.0 (1.0000)
	ECL-2.0 (0.9669)
	SHL-0.51 (0.9489)
	SHL-0.5 (0.9488)
	ImageMagick (0.8764)
testdata/project3/subproject/license
	MIT (0.9677)
	JSON (0.9554)
	Xnet (0.8972)
	MIT-feh (0.8907)
	MIT-0 (0.8904)
	MIT-advertising (0.8424)
	X11 (0.8340)
	MITNFA (0.8191)
	SGI-B-2.0 (0.7619)
```

## `mkliossdb`

`mkliossdb` creates database for `lioss` from given LICENSE data.
The resultant database is written to `default.liossdb` in json format as default.
if the extension of dest file is `.liossgz`, the resultant database is gzipped json file.

Supported algorithm is `kgram` (k=1, ..., 9), `wordfreq` and `tfidf`.

```sh
mkliossdb [OPTIONS] <LICENSE...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'default.liossdb'
    -h, --help               print this message.
LICENSE
    specifies license files.
```

