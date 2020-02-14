package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DBModel interface{
	Save()
}

type DownloadShared struct {
	Id 				int64 	`jkdb:"primary key"`
	Create_date 	int64
	File_path 		string
	Url_query 		string
	Expired_date 	int64
	Auth_code 		string
}

func (d *DownloadShared) Save()  {

	lasId,err := Save(d,d.Id<=0)
	d.Id = lasId
	fmt.Println("Save Error : ",err)
}

var database *sql.DB

func InitSqliteDB(dbName string,tables []DBModel) error {

	var err error
	database,err = sql.Open("sqlite3","./"+dbName)

	if err != nil {
		fmt.Println("Database init Error: ",err)
		return err
	}

	for _,v := range tables {
		err = createTable(v)
		if err != nil {
			break;
		}
	}

	if err != nil {
		return  err
	}

	return nil
}

func createDownloadSharedTable() error {

	err:= createTable(new(DownloadShared))
	return err
}

func createTable(model DBModel) error  {

	query := queryCreateTable(model)

	fmt.Println("Table Create: ",query)

	dbExc,err := database.Prepare(query)

	if err != nil {
		return err
	}

	_,err = dbExc.Exec()

	return err
}

func Save(m DBModel,force bool) (int64, error) {
	var query string
	if force {
		query = queryInsertData(m)
	}else{
		query = queryReplaceData(m)
	}

	fmt.Println("Save Query : ",query)

	dbExc,err := database.Prepare(query)

	if err != nil {
		return -1,err
	}

	res,err := dbExc.Exec()

	if err != nil {
		return -1,err
	}

	resId,err := res.LastInsertId()

	if err != nil {
		return -1,err
	}

	return resId,err
}

func underScoreCase(ori string) string {
	reg := regexp.MustCompile("([A-Z])")

	resStr := reg.ReplaceAllString(ori,"_$1")

	resName := strings.ToLower(resStr)

	if strings.HasPrefix(resName,"_") {
		resName = resName[1:]
	}

	if len(resName) <= 0 {
		return ori
	}

	return resName
}

func dbTypeFromType(p reflect.Type) string {
	k := p.Kind()
	switch k {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Bool:
		return "integer"
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	default:
		return "null"
	}
}

func dbValueFromModel(value reflect.Value) string {

	k := value.Kind()

	switch k {
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		return strconv.FormatInt(value.Int(),10)
	case reflect.Bool:
		if value.Bool() {
			return "1"
		}else {
			return "0"
		}
	case reflect.Float32:
		return strconv.FormatFloat(value.Float(),'f',-1,32)
	case reflect.Float64:
		return strconv.FormatFloat(value.Float(),'f',-1,64)
	case reflect.String:
		if len(value.String()) > 0 {
			return "'" + value.String() + "'"
		}else{
			return "null"
		}
	default:
		return "null"
	}
}

func queryCreateTable(model DBModel) string {

	v := reflect.TypeOf(model)

	names := strings.Split(v.String(),".")

	modelName := names[len(names) - 1]

	tableName := underScoreCase(modelName)
	fmt.Println("table : ",tableName);

	ve := v.Elem()

	query := "create table if not exists " + tableName + "("

	query = query + queryLiteralType(ve, func(sf reflect.StructField,idx int) string {
		return strings.ToLower(sf.Name) + " " + dbTypeFromType(sf.Type) + " " + sf.Tag.Get("jkdb")
	})

	query = query + ");"

	return query
}

func queryLiteralType(p reflect.Type,appendFunc func(sf reflect.StructField,idx int) string) string {

	query := ""

	for i:=0;i<p.NumField();i++ {
		f := p.Field(i)

		if f.Type.Kind() == reflect.Struct {
			query = query + queryLiteralType(f.Type,appendFunc)
		}else {
			str := appendFunc(f,i)
			if len(str) <= 0 {
				continue
			}
			query = query + str
		}

		if i == p.NumField() - 1 {
		}else {
			query = query + ","
		}
		//fmt.Printf("%d: %s %s = %s\n", i,
		//	ve.Field(i).Name, f.Type, f.Tag.Get("test"))
	}

	return query
}

func queryInsertData(m DBModel) string  {
	v := reflect.TypeOf(m)

	names := strings.Split(v.String(),".")

	modelName := names[len(names) - 1]

	tableName := underScoreCase(modelName)

	fmt.Println("Insert ",m)
	fmt.Printf("Insert : %v %s\n",v,tableName)

	query := "insert into " + tableName + " ("

	primaryKeyIdx := -1

	query = query + queryLiteralType(v.Elem(), func(sf reflect.StructField,idx int) string {
		dbInfo := sf.Tag.Get("jkdb")
		if strings.Contains(strings.ToLower(dbInfo),"primary key")  {
			primaryKeyIdx = idx
			return ""
		}else {
			return sf.Name + " "
		}
	})

	query = query + ") values ("

	f := reflect.ValueOf(m).Elem()

	query = query + queryLiteralValue(f, func(sf reflect.Value,idx int) string {
		if idx == primaryKeyIdx {
			return ""
		}
		return dbValueFromModel(sf)
	})

	query = query + ");"

	return query
}

func queryReplaceData(m DBModel) string{

	v := reflect.TypeOf(m)

	names := strings.Split(v.String(),".")

	modelName := names[len(names) - 1]

	tableName := underScoreCase(modelName)

	fmt.Println("Replace ",m)
	fmt.Printf("Replace : %v %s\n",v,tableName)

	query := "replace into " + tableName + " values ("

	f := reflect.ValueOf(m).Elem()

	query = query + queryLiteralValue(f, func(sf reflect.Value,idx int) string {
		return dbValueFromModel(sf)
	})

	query = query + ");"

	return query
}

func queryLiteralValue(p reflect.Value,appendFunc func(sf reflect.Value, idx int) string) string {

	query := ""

	for i:=0;i<p.NumField();i++ {
		f := p.Field(i)

		if f.Kind() == reflect.Struct {
			query = query + queryLiteralValue(f,appendFunc)
		}else {
			str := appendFunc(f,i)
			if len(str) <= 0 {
				continue
			}
			query = query + str
		}

		if i == p.NumField() - 1 {
		}else {
			query = query + ","
		}
		//fmt.Printf("%d: %s %s = %s\n", i,
		//	f.Field(i).Type(), f.Type(), f.Type())
	}

	return query
}

func Get(m interface{},where string)  {

	v := reflect.TypeOf(m)

	isSlice := false

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() == reflect.Slice {
			isSlice = true
		}
	}

	names := strings.Split(v.String(),".")

	modelName := names[len(names) - 1]

	tableName := underScoreCase(modelName)

	fmt.Println("Select ",m)
	fmt.Printf("Select : %v %s\n",v,tableName)

	query := "select * from " + tableName + " where " + where

	fmt.Println("Query : ",query)

	rows, err := database.Query(query)

	if err != nil {
		fmt.Println("Query Error : ",err)
		return
	}

	if isSlice {

		p := v.Elem()
		valuePtr := reflect.ValueOf(m)
		valueElm := valuePtr.Elem()
		for rows.Next() {
			vp := reflect.New(p.Elem())
			vv := reflect.ValueOf(vp.Interface()).Elem()
			values := queryLiteralValueInterface(vv, func(sf reflect.Value, idx int) interface{} {
				//fmt.Println("inn can set: ",sf.CanSet())
				r := dbInterfaceValueFromModel(sf)
				return r
			})
			rows.Scan(values...)
			valueElm.Set(reflect.Append(valueElm,vp))
		}

	}else {
		p := reflect.ValueOf(m).Elem()

		fmt.Println("what ",p)

		values := queryLiteralValueInterface(p, func(sf reflect.Value, idx int) interface{} {
			//fmt.Println("inn can set: ",sf.CanSet())
			r := dbInterfaceValueFromModel(sf)
			return r
		})
		for rows.Next() {
			rows.Scan(values...)
			break
		}
	}

	//r := new(DownloadShared)
	//inface := []interface{}{}
	//inface = append(inface,&r.Id,&r.Create_date,&r.File_path,&r.Url_query,&r.Expired_date,&r.Auth_code)
	//fmt.Println("inface : ",inface)
	//for rows.Next() {
	//	fmt.Println("next :",&r,&r.Id,&r.Create_date,&r.File_path,&r.Url_query,&r.Expired_date,&r.Auth_code)
	//	fmt.Printf("next %v: %v %v %v %v %v %v \n",&r,&r.Id,&r.Create_date,&r.File_path,&r.Url_query,&r.Expired_date,&r.Auth_code)
	//	//rows.Scan(&r.id,&r.create_date,&r.file_path,&r.url_query,&r.expired_date,&r.auth_code)
	//	rows.Scan(inface...)
	//}
	//fmt.Println("inner test : ",r)
}

func queryLiteralValueInterface(p reflect.Value ,appendFunc func(sf reflect.Value, idx int) interface{}) []interface{} {

	var valuesInterface []interface{}

	for i:=0;i<p.NumField();i++ {
		f := p.Field(i)

		if f.Kind() == reflect.Struct {
			valuesInterface = append(valuesInterface,queryLiteralValueInterface(f,appendFunc))
		}else {
			valuesInterface = append(valuesInterface,appendFunc(f,i))
		}
		//fmt.Println("interface : ",valuesInterface)
		//fmt.Println("f : ",f,f.UnsafeAddr(),f.Addr()," p : ",p,p.UnsafeAddr(),p.Addr())
		//fmt.Printf("%d: %s %s = %s\n", i,
		//	p.Field(i).Type(), f.Type(), f.Type())
	}

	return valuesInterface
}


func dbInterfaceValueFromModel(value reflect.Value) interface{} {

	k := value.Kind()
	switch k {
	case reflect.Int:
		return value.Addr().Interface().(*int)
	case reflect.Int8:
		return value.Addr().Interface().(*int8)
	case reflect.Int16:
		return value.Addr().Interface().(*int16)
	case reflect.Int32:
		return value.Addr().Interface().(*int32)
	case reflect.Int64:
		return value.Addr().Interface().(*int64)
	case reflect.Bool:
		return value.Addr().Interface().(*bool)
	case reflect.Float32:
		return value.Addr().Interface().(*float32)
	case reflect.Float64:
		return value.Addr().Interface().(*float64)
	case reflect.String:
		return value.Addr().Interface().(*string)
	default:
		return value.Addr().Interface()
	}

}