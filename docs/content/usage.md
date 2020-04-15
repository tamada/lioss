---
title: "Usage"
---

## `lioss`

`lioss` identifies license names of specified project directories, zip file and/or LICENSE files.

```sh
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
    project directories, and/or archive files contains LICENSE file.
```

### Example

```sh
$ lioss LICENSE    # run on the lioss directory.
LICENSE
	WTFPL (0.9285)
$ lioss --algorithm 9gram testdata   # run lioss for identifying project licenses in testdata directory.
testdata/project1/LICENSE
	WTFPL (0.9403)
testdata/project3/license
	Apache-License-2.0 (1.0000)
testdata/project2/license.txt
	GPLv3.0 (1.0000)
	AGPLv3.0 (0.9694)
testdata/project3/subproject/license
	BSD (1.0000)
```

## `mkliossdb`

`mkliossdb` creates database for `lioss` from given LICENSE data.
The resultant database is written to `liossdb.json` in json format as default.

Supported algorithm is `kgram` (k=1, ..., 9), `wordfreq` and `tfidf`.

```sh
mkliossdb [OPTIONS] <LICENSE...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
    -h, --help               print this message.
LICENSE
    specifies license files.
```
