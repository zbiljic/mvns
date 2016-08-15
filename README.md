# mvns

CLI tool to search artifacts in maven central repository.

Install
-------

    $ go get github.com/zbiljic/mvns

Update
------

When new code is merged to master, you can use

    $ go get -u github.com/zbiljic/mvns

To retrieve the latest version of `mvns`.

Usage
-----

### Search by query

```sh
> mvns netty
compile 'com.orange.redis-protocol:netty:0.12'                     2014-11-18
compile 'org.jboss.netty:netty:3.2.10.Final'                       2013-11-28
compile 'com.github.spullara.redis:netty:0.7'                      2013-06-22
compile 'io.netty:netty:4.0.0.Alpha8'                              2012-12-03
compile 'org.jboss.errai.io.netty:netty:4.0.0.Alpha1.errai.r1'     2012-02-23
compile 'io.netty:netty-parent:5.0.0.Alpha2'                       2015-03-03
...
```

### Search by groupId, artifactId and version

- `-g` Specify groupId.
- `-a` Specify artifactId.
- `-v` Specify version.

```sh
mvns -g io.netty -a netty-all -v 4.1.4.Final
compile 'io.netty:netty-all:4.1.4.Final'    2016-07-27
```

### Select all versions

- `-A` Show all versions.

```sh
mvns -g io.netty -a netty-all -A
compile 'io.netty:netty-all:4.1.4.Final'     2016-07-27
compile 'io.netty:netty-all:4.0.40.Final'    2016-07-27
compile 'io.netty:netty-all:4.1.3.Final'     2016-07-15
compile 'io.netty:netty-all:4.0.39.Final'    2016-07-15
compile 'io.netty:netty-all:4.1.2.Final'     2016-07-01
compile 'io.netty:netty-all:4.0.38.Final'    2016-07-01
compile 'io.netty:netty-all:4.0.37.Final'    2016-06-07
compile 'io.netty:netty-all:4.1.1.Final'     2016-06-07
compile 'io.netty:netty-all:4.1.0.Final'     2016-05-25
...
```

### Limit number of result.

- `-m` Limit number of result. Default is 20.

```sh
mvns -m 1 io.netty
compile 'io.netty:netty-parent:5.0.0.Alpha2'   2015-03-03
```

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).

---

Copyright © 2016 Nemanja Zbiljić
