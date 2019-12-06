# lioss

License Identification tool for OSS project.

## Usage

```
$ lioss [OPTIONS] <PROJECTS...>
OPTIONS
        --dbpath <DBPATH>          specifying database path.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is tfidf.
                                   Available values are: tfidf, kgram, ...
    -t, --threshold <THRESHOLD>    specifies threshold for the algorithm.
                                   Each algorithm has default value.
    -h, --help                     print this message.
PROJECTS
    project directories, and/or archive files contains LICENSE file.
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
