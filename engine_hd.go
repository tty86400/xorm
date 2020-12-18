// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

func (engine *Engine) HdSync(beans ...interface{}) error {
	s := engine.NewSession()
	defer s.Close()
	return s.HdSync(beans...)
}
