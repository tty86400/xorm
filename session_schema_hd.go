// Copyright 2016 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"fmt"
	"strings"

	"github.com/chinahdkj/core"
)




// Sync2 synchronize structs to database tables
func (session *Session) HdSync(beans ...interface{}) error {
	engine := session.engine

	if session.isAutoClose {
		session.isAutoClose = false
		defer session.Close()
	}

	targetTables := []string{}
	for _, bean := range beans {
		v := rValue(bean)
		table, err := engine.mapType(v)
		if err != nil {
			return err
		}
		var tbName= session.tbNameNoSchema(table)
		targetTables = append(targetTables,tbName)
	}

	fmt.Println("XXXXXXXXX-targetTables:",targetTables)


	tables, err := engine.HdDBMetas(targetTables)
	if err != nil {
		return err
	}

	var structTables []*core.Table

	for _, bean := range beans {
		v := rValue(bean)
		table, err := engine.mapType(v)
		if err != nil {
			return err
		}
		structTables = append(structTables, table)
		var tbName = session.tbNameNoSchema(table)

		var oriTable *core.Table
		for _, tb := range tables {
			if strings.EqualFold(tb.Name, tbName) {
				oriTable = tb
				break
			}
		}

		if oriTable == nil {
			err = session.StoreEngine(session.statement.StoreEngine).createTable(bean)
			if err != nil {
				return err
			}

			err = session.createUniques(bean)
			if err != nil {
				return err
			}

			err = session.createIndexes(bean)
			if err != nil {
				return err
			}
		} else {
			for _, col := range table.Columns() {
				var oriCol *core.Column
				for _, col2 := range oriTable.Columns() {
					if strings.EqualFold(col.Name, col2.Name) {
						oriCol = col2
						break
					}
				}

				if oriCol != nil {
					expectedType := engine.dialect.SqlType(col)
					curType := engine.dialect.SqlType(oriCol)
					if expectedType != curType {
						if expectedType == core.Text &&
							strings.HasPrefix(curType, core.Varchar) {
							// currently only support mysql & postgres
							if engine.dialect.DBType() == core.MYSQL ||
								engine.dialect.DBType() == core.POSTGRES {
								engine.logger.Infof("Table %s column %s change type from %s to %s\n",
									tbName, col.Name, curType, expectedType)
								_, err = session.exec(engine.dialect.ModifyColumnSql(table.Name, col))
							} else {
								engine.logger.Warnf("Table %s column %s db type is %s, struct type is %s\n",
									tbName, col.Name, curType, expectedType)
							}
						} else if strings.HasPrefix(curType, core.Varchar) && strings.HasPrefix(expectedType, core.Varchar) {
							if engine.dialect.DBType() == core.MYSQL {
								if oriCol.Length < col.Length {
									engine.logger.Infof("Table %s column %s change type from varchar(%d) to varchar(%d)\n",
										tbName, col.Name, oriCol.Length, col.Length)
									_, err = session.exec(engine.dialect.ModifyColumnSql(table.Name, col))
								}
							}
						} else {
							if !(strings.HasPrefix(curType, expectedType) && curType[len(expectedType)] == '(') {
								engine.logger.Warnf("Table %s column %s db type is %s, struct type is %s",
									tbName, col.Name, curType, expectedType)
							}
						}
					} else if expectedType == core.Varchar {
						if engine.dialect.DBType() == core.MYSQL {
							if oriCol.Length < col.Length {
								engine.logger.Infof("Table %s column %s change type from varchar(%d) to varchar(%d)\n",
									tbName, col.Name, oriCol.Length, col.Length)
								_, err = session.exec(engine.dialect.ModifyColumnSql(table.Name, col))
							}
						}
					}
					if col.Default != oriCol.Default {
						engine.logger.Warnf("Table %s Column %s db default is %s, struct default is %s",
							tbName, col.Name, oriCol.Default, col.Default)
					}
					if col.Nullable != oriCol.Nullable {
						engine.logger.Warnf("Table %s Column %s db nullable is %v, struct nullable is %v",
							tbName, col.Name, oriCol.Nullable, col.Nullable)
					}
				} else {
					session.statement.RefTable = table
					session.statement.tableName = tbName
					err = session.addColumn(col.Name)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	for _, table := range tables {
		var oriTable *core.Table
		for _, structTable := range structTables {
			if strings.EqualFold(table.Name, session.tbNameNoSchema(structTable)) {
				oriTable = structTable
				break
			}
		}

		if oriTable == nil {
			//engine.LogWarnf("Table %s has no struct to mapping it", table.Name)
			continue
		}

		for _, colName := range table.ColumnsSeq() {
			if oriTable.GetColumn(colName) == nil {
				engine.logger.Warnf("Table %s has column %s but struct has not related field", table.Name, colName)
			}
		}
	}
	return nil
}
