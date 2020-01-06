[![License](https://img.shields.io/badge/License-WTFPL-blue.svg)](https://github.com/tamada/lioss/blob/master/LICENSE)
[![Version](https://img.shields.io/badge/Version-1.0.0-yellowgreen.svg)](https://github.com/tamada/lioss/releases/tag/v1.0.0)

# lioss

License Identification tool for OSS project.

## Usage

### `lioss`

```
$ lioss [OPTIONS] <PROJECTS...>
OPTIONS
        --dbpath <DBPATH>          specifying database path.
                                   the default value is $LIOSS_HOME/licenses.db
                                   The database is build by 'mkliossdb' command.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is tfidf.
                                   Available values are: tfidf, kgram, ...
    -t, --threshold <THRESHOLD>    specifies threshold for the algorithm.
                                   Each algorithm has default value.
    -h, --help                     print this message.
PROJECTS
    project directories, and/or archive files contains LICENSE file.
```

### `mkliossdb`

Creates the database of lioss.
If the database was exists, this command updates the database.

```
$ mkliossdb [OPTIONS] <DIRs...>
OPTIONS
    -d, --dest <PATH>             specifies the destination path of licenses.db.
                                  the default value is $LIOSS_HOME/licenses.db
DIRs...
    directories containing the license files.
    The license files must name the license name.
```

### Algorithm

#### tfidf

Calculate *tfidf* from each license files, then filter by *threshold*, and obtain important term frequencies.
Calculate cosine similarities of important term frequencies.

#### k-grams

Construct frequencies based on k-grams of characters.
Then calculate cosine similarities of the frequencies.


## References

* [dmgerman/ninka](https://github.com/dmgerman/ninka)
    * Daniel M. German, Yuki Manabe and Katsuro Inoue. A sentence-matching method for automatic license identification of source code files. In 25nd IEEE/ACM International Conference on Automated Software Engineering (ASE 2010).
