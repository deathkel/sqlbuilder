package test

import (
    "github.com/deathkel/sqlbuilder/builder"
    "testing"
    "reflect"
)

func Test_Select(t *testing.T) {
    b := new(builder.Builder)
    sql, bindings := b.Select([]string{"*", "sex", "a.name", "count(1) as count"}).
        From("user").
        Where("a", "1").
        GroupBy("a").
        Having("a", ">", "2").
        Limit("3").
        Offset("4").
        ToSql()
    if sql != "select *, `sex`, a.name, count(1) as count from `user` where (`a` = ?) group by `a` having `a` > ? offset ? limit ?" {
        t.Error(sql)
    }
    if !reflect.DeepEqual(bindings, []string{"1", "2", "4", "3"}) {
        t.Error(bindings)
    }
}

func Test_Join(t *testing.T) {
    b := new(builder.Builder)
    
    sql, bindings := b.Select([]string{"*"}).From("ta").
        Join("tb", "`tb`.`aid` = `ta`.`id`").
        Where("ta.id", ">", "1").
        Where("tb.name", "=", "jack").
        ToSql()
    
    if sql != "select * from `ta` join `tb` on `tb`.`aid` = `ta`.`id` where (ta.id > ? and tb.name = ?)" {
        t.Error(sql)
    }
    if !reflect.DeepEqual(bindings, []string{"1", "jack"}) {
        t.Error(bindings)
    }
}

func Test_Insert(t *testing.T) {
    b := new(builder.Builder)
    
    info := map[string]string{"name": "john"}
    sql, bindings := b.Insert("ta", info).ToSql()
    if sql != "insert into `ta` (`name`) values (?)" {
        t.Error(sql)
    }
    
    if !reflect.DeepEqual(bindings, []string{"john"}) {
        t.Error(bindings)
    }
}

func Test_Update(t *testing.T) {
    b := new(builder.Builder)
    
    info := map[string]string{"name": "john"}
    sql, bindings := b.Update("ta", info).
        Where("name", "kel").
        Where("sex", "2").
        Offset("1").
        Limit("2").
        ToSql()
    if sql != "update `ta` set `name` = ? where (`name` = ? and `sex` = ?) offset ? limit ?" {
        t.Error(sql)
    }
    
    if !reflect.DeepEqual(bindings, []string{"john", "kel", "2", "1", "2"}) {
        t.Error(bindings)
    }
}

func Test_Delete(t *testing.T) {
    b := new(builder.Builder)
    sql, bindings := b.Delete("ta").
        Where("name", "kel").
        Where("sex", "2").
        Offset("1").
        Limit("2").
        ToSql()
    if sql != "delete `ta` where (`name` = ? and `sex` = ?) offset ? limit ?" {
        t.Error(sql)
    }
    
    if !reflect.DeepEqual(bindings, []string{"kel", "2", "1", "2"}) {
        t.Error(bindings)
    }
}
