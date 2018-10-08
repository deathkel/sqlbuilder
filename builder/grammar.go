package builder

import (
    "strings"
)

type Grammar struct {
    selectComponents []string
    operators        []string
}

func wrapValue(value string) (string) {
    if value == "*" {
        return value
    }
    if strings.ContainsAny(value, ".()") {
        return value
    }
    
    value = strings.Replace(value, "`", "", -1)
    return "`" + value + "`"
}

func CompileSelect(query *Builder) (sql string, bindings []string) {
    
    //select
    sql = addSelect(query)
    
    //from
    sql += addFrom(query)
    
    //join
    sql += addJoin(query)
    
    //where
    sql += addWhere(query)
    
    //group by
    sql += addGroupBy(query)
    
    //having
    sql += addHaving(query)
    
    //order by
    sql += addOrderBy(query)
    
    //offset
    sql += addOffset(query)
    
    //limit
    sql += addLimit(query)
    
    bindings = append(query.bindings.where, query.bindings.having...)
    if query.offset != "" {
        bindings = append(bindings, query.offset)
    }
    if query.limit != "" {
        bindings = append(bindings, query.limit)
    }
    
    return sql, bindings
}

func CompileUpdate(query *Builder) (sql string, bindings []string) {
    //update
    sql = addUpdate(query)
    
    //where
    sql += addWhere(query)
    
    //offset
    sql += addOffset(query)
    
    //limit
    sql += addLimit(query)
    
    bindings = append(query.bindings.update, query.bindings.where...)
    if query.offset != "" {
        bindings = append(bindings, query.offset)
    }
    if query.limit != "" {
        bindings = append(bindings, query.limit)
    }
    
    return sql, bindings
}

func CompileInsert(query *Builder) (sql string, bindings []string) {
    //insert
    sql = addInsert(query)
    
    //offset
    sql += addOffset(query)
    
    //limit
    sql += addLimit(query)
    
    bindings = append(query.bindings.insert, query.bindings.where...)
    if query.offset != "" {
        bindings = append(bindings, query.offset)
    }
    if query.limit != "" {
        bindings = append(bindings, query.limit)
    }
    
    return sql, bindings
}

func CompileDelete(query *Builder) (sql string, bindings []string) {
    //delete
    sql = addDelete(query)
    
    //where
    sql += addWhere(query)
    
    //offset
    sql += addOffset(query)
    
    //limit
    sql += addLimit(query)
    
    bindings = append(bindings, query.bindings.where...)
    if query.offset != "" {
        bindings = append(bindings, query.offset)
    }
    if query.limit != "" {
        bindings = append(bindings, query.limit)
    }
    
    return sql, bindings
}

func addInsert(query *Builder) (sql string) {
    sql = "insert into " + wrapValue(query.table) + " ("
    lenInsert := len(query.insert)
    idx := 0
    for _, column := range query.insert {
        idx ++
        sql += wrapValue(column)
        if idx < lenInsert {
            sql += ", "
        }
    }
    sql += ") values ("
    for i := 1; i <= lenInsert; i++ {
        sql += "?"
        if i < lenInsert {
            sql += ", "
        }
    }
    sql += ")"
    return
}

func addUpdate(query *Builder) (sql string) {
    sql = "update " + wrapValue(query.table) + " set "
    lenUpdate := len(query.update)
    idx := 0
    for _, column := range query.update {
        idx ++
        sql += wrapValue(column) + " = ?"
        if idx < lenUpdate {
            sql += ", "
        }
    }
    return
}

func addDelete(query *Builder) (sql string) {
    sql = "delete " + wrapValue(query.table)
    return
}

func addSelect(query *Builder) (string) {
    sql := "select"
    lenColumns := len(query.columns)
    if lenColumns == 0 {
        sql += " *"
    } else {
        for key, column := range query.columns {
            sql += " " + wrapValue(column)
            if key < lenColumns-1 {
                sql += ","
            }
        }
    }
    return sql
}

func addFrom(query *Builder) (sql string) {
    sql = " from " + wrapValue(query.table)
    return
}

func addJoin(query *Builder) (sql string) {
    lenJoin := len(query.joins)
    if lenJoin == 0 {
        return sql
    }
    for _, val := range query.joins {
        switch val.joinType {
        case "left":
            sql += " left join "
        case "right":
            sql += " right join "
        default:
            sql += " join "
        }
        sql += wrapValue(val.table) + " on " + val.condition
    }
    return sql
}

func addWhere(query *Builder) (sql string) {
    lenWhere := len(query.wheres)
    if lenWhere > 0 {
        sql += " where ("
    }
    for key, where := range query.wheres {
        sql += wrapValue(where.column) + " " + where.operator + " ?"
        if key < lenWhere-1 {
            sql += " " + where.boolean + " "
        }
    }
    sql += ")"
    return sql
}

func addGroupBy(query *Builder) (sql string) {
    lenGroups := len(query.groups)
    if lenGroups > 0 {
        sql += " group by"
    }
    for key, column := range query.groups {
        sql += " " + wrapValue(column)
        if key < lenGroups-1 {
            sql += ","
        }
    }
    return sql
}

func addHaving(query *Builder) (sql string) {
    lenHavings := len(query.havings)
    if lenHavings > 0 {
        sql += " having"
    }
    for key, where := range query.havings {
        sql += " " + wrapValue(where.column) + " " + where.operator + " ?"
        if key < lenHavings-1 {
            sql += " " + where.boolean
        }
    }
    return sql
}

func addOrderBy(query *Builder) (sql string) {
    lenOrder := len(query.orders)
    if lenOrder > 0 {
        sql += " order by"
    }
    for key, order := range query.orders {
        sql += " " + order.column + " " + order.direction
        if key < lenOrder-1 {
            sql += ","
        }
    }
    return sql
}

func addOffset(query *Builder) (sql string) {
    if query.offset != "" {
        sql += " offset ?"
    }
    return sql
}

func addLimit(query *Builder) (sql string) {
    if query.limit != "" {
        sql += " limit ?"
    }
    return sql
}
