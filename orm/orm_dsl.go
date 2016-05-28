package orm

import (
	"fmt"
	"reflect"
)

// createTableSQL() must traverse the table structure, colect its fields and
// assemble the sql CREATE TABLE instruction
func (handler Handler) assembleSQLCreateTable() (string, error) {
	var existsID = false

	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	fieldsList := ""
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		fieldType := ""
		if typeOfTable.Field(i).Type.Name() == "int" {
			fieldType = "integer"
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			fieldType = "character varying"
		}

		if fieldName == "ID" {
			existsID = true
			fieldsList = fieldsList + fieldName + " serial NOT NULL, "
		} else {
			fieldsList = fieldsList + fieldName + " " + fieldType + ", "
		}
	}

	if existsID == false {
		return "", fmt.Errorf("ID field not found on struct %v", tableName)
	}

	primaryKey := "constraint " + tableName + "_pkey primary key (id)"
	sqlInstruction := "create table " + tableName + " (" + fieldsList + " " + primaryKey + ");"

	handler.sqlCreateTable = sqlInstruction

	return sqlInstruction, nil
}

// createTable() must execute the sql CREATE TABLE instruction
func (handler Handler) createTable() (err error) {
	sqlInstruction := ""

	if handler.sqlCreateTable != "" {
		sqlInstruction = handler.sqlCreateTable
	} else {
		sqlInstruction, err = handler.assembleSQLCreateTable()
		if err != nil {
			return err
		}
	}

	_, err = handler.db.Exec(sqlInstruction)
	return err
}

func (handler Handler) dropTable() error {
	tableName := reflect.TypeOf(handler.table).Name()
	sqlInstruction := "drop table " + tableName + ";"

	_, err := handler.db.Exec(sqlInstruction)

	return err
}
