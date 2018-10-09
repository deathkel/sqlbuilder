## golang SQL builder
### this is a sql builder for golang in chain style 

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Quality Score](https://scrutinizer-ci.com/g/deathkel/sqlbuilder/badges/quality-score.png?b=master)](https://scrutinizer-ci.com/g/deathkel/sqlbuilder)
#### USAGE:
#### EXAMPLE
##### SELECT
* case 1:``select * from user``
```go
b := new builder.Builder()
sql, bindings := b.Select("*").From("user").toSql()

```

* case 2: ``select *, `user`.`name`, count(1) as count from `user` where (id = ?) ``
```go
...
sql, bindings := b.Select("*", "user.name", "count(1) as count").From("user").Where("id", "1").toSql()
...
```
##### INSERT
* case 1: ``insert into `user` (`name`, `sex`, `age`) values (?, ?, ?)``
```go
...
info := map[string]string{"name":"john", "sex":"2", "age":"22"}
sql, bindings := b.Insert("user", info).toSql()
...
```

##### UPDATE
* case 1: ``update `user` set `name` = ?, `sex` = ?, `age` = ? where (id = ?)``
```go
...
info := map[string]string{"name":"john", "sex":"2", "age":"22"}
sql, bindings := b.Update("user", info).Where("id", "1").toSql()
...
```

##### DELETE
* case 1: ``delete `user` where (id = ?)``
```go
...
sql, bindings := b.Delete("user").Where("id", "1").toSql()
...
```


### support database
* mysql