package test

import (
    "sqlbuilder/builder"
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
    if sql != "select *, `sex`, a.name, count(1) as count from `user` where `a` = ? group by `a` having `a` > ? offset ? limit ?" {
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
    
    if sql != "select * from `ta` join `tb` on `tb`.`aid` = `ta`.`id` where ta.id > ? and tb.name = ?" {
        t.Error(sql)
    }
    if !reflect.DeepEqual(bindings, []string{"1", "jack"}) {
        t.Error(bindings)
    }
}
