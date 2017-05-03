package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type Model interface {
	id() *int64
	setId(int64)

	table() string
}

func save(m Model, db *sql.DB) error {
	v := reflect.ValueOf(m).Elem()
	t := reflect.TypeOf(m).Elem()
	keys := make([]string, 0, t.NumField()-1)
	qs := make([]string, 0, t.NumField()-1)
	values := make([]interface{}, 0, t.NumField())
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if f.Name == "_id" {
			continue
		}
		keys = append(keys, f.Tag.Get("db"))
		values = append(values, fv.Interface())
		qs = append(qs, "?")
	}

	if m.id() == nil {
		insertQuery := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);", m.table(), strings.Join(keys, ", "), strings.Join(qs, ", "))
		res, err := db.Exec(insertQuery, values...)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		m.setId(id)
	} else {
		updateQuery := fmt.Sprintf("UPDATE %s SET %s=? WHERE id=?;", m.table(), strings.Join(keys, "=?, "))
		values = append(values, m.id())
		_, err := db.Exec(updateQuery, values...)
		if err != nil {
			return err
		}
	}
	return nil
}

func load(m Model, db *sql.DB) error {
	v := reflect.ValueOf(m).Elem()
	t := reflect.TypeOf(m).Elem()
	keys := make([]string, 0, t.NumField()-1)
	fields := make([]*reflect.Value, 0, t.NumField()-1)
	values := make([]interface{}, 0, t.NumField()-1)
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if f.Name == "_id" {
			continue
		}
		keys = append(keys, f.Tag.Get("db"))
		fields = append(fields, &fv)
		switch f.Type.Kind() {
		case reflect.Int:
			values = append(values, new(int))
		case reflect.Bool:
			values = append(values, new(bool))
		case reflect.String:
			values = append(values, new(string))
		}

	}

	rows, err := db.Query(fmt.Sprintf("SELECT %s FROM %s WHERE id=?", strings.Join(keys, ", "), m.table()), m.id())
	if err != nil {
		return err
	}
	defer rows.Close()
	if !rows.Next() {
		return fmt.Errorf("id not found: %s", m.id())
	}
	err = rows.Scan(values...)
	if err != nil {
		return err
	}
	for i, field := range fields {
		field.Set(reflect.ValueOf(values[i]).Elem())
	}

	return nil
}

func loadAll(m Model, db *sql.DB) ([]interface{}, error) {
	t := reflect.TypeOf(m).Elem()
	rows, err := db.Query(fmt.Sprintf("SELECT id FROM %s", m.table()))
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	ifaces := make([]interface{}, 0)
	var id int64
	for rows.Next() {
		model := reflect.New(t).Interface().(Model)
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		model.setId(id)
		if err := load(model, db); err != nil {
			return nil, err
		}
		ifaces = append(ifaces, model)

	}
	return ifaces, nil
}
