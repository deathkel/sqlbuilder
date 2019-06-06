package sqlbuilder

import (
    "reflect"
    "strconv"
)

type Builder struct {
    //The query value bindings
    bindings

    //The Columns that should be return
    columns []string

    //The table which the query is targeting
    table string

    //Update columns
    update []interface{}

    //Insert columns
    insert []string

    //Indicates is delete sql.
    delete bool

    //Indicates if the query return district results.
    district bool

    //The where constraints for the query.
    wheres []*where

    //The groupings for the query.
    groups []string

    //The havings constraints for the query.
    havings []*where

    //The table joins for the query.
    joins []*join

    //The orderings for the query.
    orders []*order

    //The maximum number of records to return.
    limit string

    //The number of records to skip.
    offset string
}

type bindings struct {
    selected []string
    insert   []string
    update   []string
    from     []string
    join     []string
    where    []string
    having   []string
    order    []string
}

//Expression for raw sql
type Expression struct {
    Value string
}

type where struct {
    column   string
    operator string
    value    interface{}
    boolean  string
}

type join struct {
    table     string
    condition string
    joinType  string // inner left right
}

type order struct {
    column    string
    direction string
}

//All of the available clause operators.
var operators = []string{"=", "<", ">", "<=", ">=", "<>", "!=", "<=>", "like", "like binary", "not like", "ilike",
    "&", "|", "^", "<<", ">>", "rlike", "regexp", "not regexp", "~", "~*", "!~", "!~*", "similar to", "not similar to",
    "not ilike", "~~*", "!~~*", "in", "not in",
}

//Is val in array
func in_array(val interface{}, array interface{}) (exists bool, index int) {
    exists = false
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }

    return
}

func (b *Builder) invalidOperator(operator string) bool {
    exists, _ := in_array(operator, operators)
    return exists
}

func (b *Builder) Select(columns []string) *Builder {
    b.columns = columns
    return b
}

func (b *Builder) From(table string) *Builder {
    b.table = table
    return b
}

/*
Where("column", "1")
Where("column", "=","1")
Where("column", "in", []string{"1","2","3"})
Where("column", "=", "1", "or")
Where(map[string]string{column1:"1", column2: "2"})
*/
func (b *Builder) Where(column interface{}, args ...interface{}) *Builder {
    switch column.(type) {
    case string:
        operator := "="
        value := ""
        booll := "and"
        lenArgs := len(args)
        if lenArgs == 1 {
            value = args[0].(string)
        } else if lenArgs >= 2 {
            operator = args[0].(string)

            if lenArgs >= 3 {
                booll = args[2].(string)
            }

            switch args[1].(type) {
            case string:
                value = args[1].(string)

            case []string:
                //支持 where in 操作
                if !b.invalidOperator(operator) {
                    operator = "="
                }
                condition := &where{column.(string), operator, args[1], booll}
                b.wheres = append(b.wheres, condition)
                b.bindings.where = append(b.bindings.where, args[1].([]string)...)
                goto end
            }
        }

        if !b.invalidOperator(operator) {
            operator = "="
        }

        condition := &where{column.(string), operator, value, booll}
        b.wheres = append(b.wheres, condition)
        b.bindings.where = append(b.bindings.where, value)
    case map[string]string:
        return b.addArrayOfWheres(column.(map[string]string), "and")
    }

end:
    return b
}

func (b *Builder) addArrayOfWheres(wheres map[string]string, boolean string) *Builder {
    for k, v := range wheres {
        condition := &where{k, "=", v, boolean}
        b.wheres = append(b.wheres, condition)
        b.bindings.where = append(b.bindings.where, v)
    }
    return b
}

/*
Join("tableB", "tableB.id = tableA.bId")
*/
func (b *Builder) Join(table string, condition string, args ...string) *Builder {
    joinType := "inner"
    lenArgs := len(args)
    if lenArgs > 0 {
        joinType = args[0]
    }
    j := &join{table, condition, joinType}
    b.joins = append(b.joins, j)
    return b
}

func (b *Builder) LeftJoin(table string, condition string) *Builder {

    return b.Join(table, condition, "left")
}

func (b *Builder) RightJoin(table string, condition string) *Builder {
    return b.Join(table, condition, "right")
}

func (b *Builder) GroupBy(group interface{}) *Builder {
    switch group.(type) {
    case []string:
        b.groups = append(b.groups, group.([]string)...)
    case string:
        b.groups = append(b.groups, group.(string))
    }
    return b
}

/*
Having("column", "1")
Having("column", "=","1")
Having("column", "=", "1", "or")
Having(map[string]string{column1:"1", column2: "2"})
*/
func (b *Builder) Having(column interface{}, operator string, args ...string) *Builder {
    switch column.(type) {
    case string:
        value := ""
        booll := "and"
        lenArgs := len(args)
        if lenArgs == 0 {
            value = operator
        } else if lenArgs >= 1 {
            value = args[0]
        } else if lenArgs >= 2 {
            booll = args[1]
        }

        if !b.invalidOperator(operator) {
            operator = "="
        }

        condition := &where{column.(string), operator, value, booll}
        b.havings = append(b.havings, condition)
        b.bindings.having = append(b.bindings.having, value)
    case map[string]string:
        return b.addArrayOfHavings(column.(map[string]string), "and")
    }

    return b
}

func (b *Builder) addArrayOfHavings(wheres map[string]string, boolean string) *Builder {
    for k, v := range wheres {
        condition := &where{k, "=", v, boolean}
        b.havings = append(b.wheres, condition)
        b.bindings.having = append(b.bindings.having, v)
    }
    return b
}

func (b *Builder) OrderBy(column string, direction string) *Builder {
    b.orders = append(b.orders, &order{column, direction})
    return b
}

func (b *Builder) OrderByAsc(column string) *Builder {
    return b.OrderBy(column, "asc")
}

func (b *Builder) OrderByDesc(column string) *Builder {
    return b.OrderBy(column, "desc")
}

func (b *Builder) Offset(offset interface{}) *Builder {
    switch offset.(type) {
    case string:
        b.offset = offset.(string)
    case int:
        b.offset = strconv.Itoa(offset.(int))
    }
    return b
}

func (b *Builder) Limit(limit interface{}) *Builder {
    switch limit.(type) {
    case string:
        b.limit = limit.(string)
    case int:
        b.limit = strconv.Itoa(limit.(int))
    }
    return b
}

func (b *Builder) Insert(table string, info map[string]string) *Builder {
    b.table = table

    for column, value := range info {
        b.insert = append(b.insert, column)
        b.bindings.insert = append(b.bindings.insert, value)
    }
    return b
}

func (b *Builder) Update(table string, info map[string]interface{}) *Builder {
    b.table = table
    for column, value := range info {

        switch value.(type) {
        case string:
            b.update = append(b.update, column)
            b.bindings.update = append(b.bindings.update, value.(string))
        case *Expression:
            //表达式： 如 a = a + 1
            b.update = append(b.update, value)
        }
    }
    return b
}

func (b *Builder) Delete(table string) *Builder {
    b.table = table
    b.delete = true
    return b
}

func (b *Builder) ToSql() (sql string, bindings []string) {
    if b.delete {
        return CompileDelete(b)
    } else if len(b.insert) > 0 {
        return CompileInsert(b)
    } else if len(b.update) > 0 {
        return CompileUpdate(b)
    } else {
        return CompileSelect(b)
    }
}
