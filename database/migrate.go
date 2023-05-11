package database

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/ivangago06/app"
	"github.com/lib/pq"
)

func Migrate(app *app.App, anyStruc interface{}) error {

	valueOfStruct := reflect.ValueOf(anyStruc)
	typeOfStruct := valueOfStruct.Type()

	tableName := typeOfStruct.Name()

	err := createTable(app, tableName)

	if err != nil {
		return err
	}

	for i := 0; i < valueOfStruct.NumField(); i++ {
		fieldType := typeOfStruct.Field(i)
		filedName := fieldType.Name

		if filedName != "ID" && filedName != "id" {
			err := createColumn(app, tableName, filedName, fieldType.Type.Name())

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createTable(app *app.App, tableName string) error {

	var tableExist bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE c.relname ~ $1 AND pg_catalog.pg_table_is_visible(c.oid))", "^"+tableName+"$").Scan(&tableExist)

	if err != nil {
		log.Println("Error validando si existe en la tabla " + tableName)
		return err
	}

	if tableExist {
		log.Println("Ya existe la tabla")
		return nil
	} else {
		sanitazedTableQuey := fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\" (\"Id\" serial primary key)", tableName)

		_, err := app.Db.Query(sanitazedTableQuey)

		if err != nil {
			log.Println("Error creando la tabla " + tableName)
			return err
		}

		log.Println("La tabla fue creadacon éxito " + tableName)
		return nil
	}

}

func createColumn(app *app.App, tableName, columnName, columnType string) error {

	var columnExist bool
	err := app.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", tableName, columnName).Scan(&columnExist)

	if err != nil {
		log.Println("Error creando la columna " + columnName + " en la tabla " + tableName)
		return err
	}

	if columnExist {
		log.Println("La coumna fue creadacon éxito " + columnName + " en la tabla " + tableName)
		return nil
	} else {
		postgresType, err := getPostgresType(columnType)

		if err != nil {
			log.Println("Error creando la columna " + columnName + " con el tipo " + columnType)
			return err
		}

		sanitazedTableName := pq.QuoteIdentifer(tableName)
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS \"%s\" %s", sanitazedTableName, columnName, postgresType)

		_, err = app.Db.Query(query)

		if err != nil {
			log.Println("Error editando la columna " + columnName + " en la tabla " + tableName)
			return err
		}

		log.Println("La columna fue creada con éxito", columnName)
		return nil

	}

}

func getPostgresType(goType string) (string, error) {

	switch goType {
	case "int", "int32", "uint", "uint32":
		return "interger", nil
	case "int64", "uint64":
		return "biginit", nil
	case "int16", "uint16", "uint8", "int8", "byte":
		return "samllint", nil
	case "string":
		return "text", nil
	case "float64":
		return "dooble precision", nil
	case "bool":
		return "boolean", nil
	case "Time":
		return "tymestamp", nil
	case "[]byte":
		return "byt", nil
	}

	return "", errors.New("Tipo noreconocido" + goType)

}
