// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"
	"github.com/chinahdkj/core"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func (engine *Engine) HdDBMetas(targetTables []string) ([]*core.Table, error) {

	tables, err := engine.dialect.GetTables()
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		if !Include(targetTables,table.Name){
			continue
		}

		colSeq, cols, err := engine.dialect.GetColumns(table.Name)
		if err != nil {
			return nil, err
		}
		for _, name := range colSeq {
			table.AddColumn(cols[name])
		}
		indexes, err := engine.dialect.GetIndexes(table.Name)
		if err != nil {
			return nil, err
		}
		table.Indexes = indexes

		for _, index := range indexes {
			for _, name := range index.Cols {
				if col := table.GetColumn(name); col != nil {
					col.Indexes[index.Name] = index.Type
				} else {
					return nil, fmt.Errorf("Unknown col %s in index %v of table %v, columns %v", name, index.Name, table.Name, table.ColumnsSeq())
				}
			}
		}
	}
	return tables, nil
}



func (engine *Engine) HdSync(beans ...interface{}) error {
	s := engine.NewSession()
	defer s.Close()
	return s.HdSync(beans...)
}
