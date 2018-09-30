package builder

import (
    "reflect"
)

type Builder struct {
    bindings
    /*
    The Columns that should be return
     */
    columns []string
    
    table string
    
    /*
    Indicates if the query return district results
     */
    district bool
    
    wheres []*where
    
    groups []string
    
    havings []*where
    
    joins []*join
    
    orders []*order
    
    limit string
    
    offset string
    
    unions []string
    
    unionList string
    
    unionLimit string
    
    unionOrders []string
}

type bindings struct {
    selected []string
    from     []string
    join     []string
    where    []string
    having   []string
    order    []string
    union    []string
}

type where struct {
    column   string
    operator string
    value    string
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

var operators = []string{"=", "<", ">", "<=", ">=", "<>", "!=", "<=>", "like", "like binary", "not like", "ilike",
    "&", "|", "^", "<<", ">>", "rlike", "regexp", "not regexp", "~", "~*", "!~", "!~*", "similar to", "not similar to",
    "not ilike", "~~*", "!~~*",
}

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

func (b *Builder) invalidOperator(operator string) (bool) {
    exists, _ := in_array(operator, operators)
    return exists
}

func (b *Builder) Select(columns []string) (*Builder) {
    b.columns = columns
    return b
}

func (b *Builder) From(table string) (*Builder) {
    b.table = table
    return b
}

/**
Where("column", "1")
Where("column", "=","1")
Where("column", "=", "1", "or")
Where(map[string]string{column1:"1", column2: "2"})
 */
func (b *Builder) Where(column interface{}, operator string, args ...string) (*Builder) {
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
        b.wheres = append(b.wheres, condition)
        b.bindings.where = append(b.bindings.where, value)
    case map[string]string:
        return b.addArrayOfWheres(column.(map[string]string), "and")
    }
    
    return b
}

func (b *Builder) addArrayOfWheres(wheres map[string]string, boolean string) (*Builder) {
    for k, v := range wheres {
        condition := &where{k, "=", v, boolean}
        b.wheres = append(b.wheres, condition)
        b.bindings.where = append(b.bindings.where, v)
    }
    return b
}

/**
Join("tableB", "tableB.id = tableA.bId")
 */
func (b *Builder) Join(table string, condition string, args ...string) (*Builder) {
    joinType := "inner"
    lenArgs := len(args)
    if lenArgs > 0 {
        joinType = args[0]
    }
    j := &join{table, condition, joinType}
    b.joins = append(b.joins, j)
    return b
}

func (b *Builder) LeftJoin(table string, condition string) (*Builder) {
    
    return b.Join(table, condition, "left")
}

func (b *Builder) RightJoin(table string, condition string) (*Builder) {
    return b.Join(table, condition, "right")
}

func (b *Builder) GroupBy(group interface{}) (*Builder) {
    switch group.(type) {
    case []string:
        b.groups = append(b.groups, group.([]string)...)
    case string:
        b.groups = append(b.groups, group.(string))
    }
    return b
}

/**
Having("column", "1")
Having("column", "=","1")
Having("column", "=", "1", "or")
Having(map[string]string{column1:"1", column2: "2"})
 */
func (b *Builder) Having(column interface{}, operator string, args ...string) (*Builder) {
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

func (b *Builder) addArrayOfHavings(wheres map[string]string, boolean string) (*Builder) {
    for k, v := range wheres {
        condition := &where{k, "=", v, boolean}
        b.havings = append(b.wheres, condition)
        b.bindings.having = append(b.bindings.having, v)
    }
    return b
}

func (b *Builder) OrderBy(column string, direction string) (*Builder) {
    b.orders = append(b.orders, &order{column, direction})
    return b
}

func (b *Builder) OrderByAsc(column string) (*Builder) {
    return b.OrderBy(column, "asc")
}

func (b *Builder) OrderByDesc(column string) (*Builder) {
    return b.OrderBy(column, "desc")
}

func (b *Builder) Offset(offset string) (*Builder) {
    b.offset = offset
    return b
}

func (b *Builder) Limit(limit string) (*Builder) {
    b.limit = limit
    return b
}

func (b *Builder) Union() (*Builder) {
    return b
}

func (b *Builder) Count() {

}

func (b *Builder) Min() {
    
}

func (b *Builder) Max() {
    
}

func (b *Builder) Sum() {
    
}

func (b *Builder) Avg() {
    
}

func (b *Builder) Insert() {
    
}

func (b *Builder) Update() {

}

func (b *Builder) Delete() {

}

func (b *Builder) Increment() {

}

func (b *Builder) Decrement() {

}

func (b *Builder) ToSql() (sql string, bindings []string) {
    sql, bindings = CompileSelect(b)
    return
}
