package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

// assembleSQLInsertStatement traverse the the object
// returns a SQL insert instruction and a string array containing the exact
// parameters order
func (handler *Handler) assembleSQLUpdateStatement() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	j := 1
	sqlInstruction := "update " + tableName + " set "
	var fieldMap []saveField
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if fieldName == "ID" {
			continue
		}

		fieldMap = append(fieldMap, saveField{name: typeOfTable.Field(i).Name, fieldType: typeOfTable.Field(i).Type.Name()})
		sqlInstruction = sqlInstruction + fieldName + " = $" + strconv.Itoa(j) + ", "
		j = j + 1
	}
	fieldMap = append(fieldMap, saveField{name: "ID", fieldType: "int"})

	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " where id = $" + strconv.Itoa(j) + ";"

	handler.sqlUpdate = sqlInstruction
	handler.mapUpdate = fieldMap
}

func (handler Handler) update(objectPtr interface{}) error {
	object := reflect.ValueOf(objectPtr).Elem()
	tableName := reflect.TypeOf(objectPtr).Elem().Name()
	if tableName != handler.tableName {
		return fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

	//build the arguments array
	var args []interface{}
	for _, field := range handler.mapUpdate {
		if field.fieldType == "int" {
			args = append(args, int(object.FieldByName(field.name).Int()))
		}
		if field.fieldType == "string" {
			args = append(args, string(object.FieldByName(field.name).String()))
		}
	}

	//run Update
	_, err := handler.db.Exec(handler.sqlUpdate, args...)
	if err != nil {
		return err
	}

	return nil
}
