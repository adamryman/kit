# kit

A group of small go libraries, each of them ready for you to do a little
copying, [as is the go way.](https://www.youtube.com/watch?v=PAAkCSZUG1c&feature=youtu.be&t=569)

If you import, please vendor, the api may change time to time.

For vendoring, I recommend [`govendor`](https://github.com/kardianos/govendor),
though the [`dep`](https://github.com/golang/dep) project may outshine
it one day.

## kit - listed

```
go get -d github.com/adamryman/kit
```

### [`kit/dbconn`](./dbconn)

[`https://godoc.org/github.com/adamryman/kit/dbconn`](https://godoc.org/github.com/adamryman/kit/dbconn)

```
    Package dbconn creates database connection strings for
    "database/sql.Open(driverName string, dataSourceName string)" where
    dataSourceName is the connection string
```

## (UN)-LICENCE

All software is released into the public domain. See [`LICENCE.md`](./LICENCE.md)
