//go:build wasm
// +build wasm

package main

import (
	"context"
	"encoding/json"
	"syscall/js"

	"github.com/Permify/permify/internal/commands"
	"github.com/Permify/permify/pkg/development"
	"github.com/Permify/permify/pkg/dsl/schema"
	"github.com/Permify/permify/pkg/errors"
	"github.com/Permify/permify/pkg/graph"
	base "github.com/Permify/permify/pkg/pb/base/v1"
	"github.com/Permify/permify/pkg/tuple"
)

var dev *development.Development

// check -
func check() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &development.CheckQuery{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), params)
		if mErr != nil {
			return js.ValueOf([]interface{}{false, mErr.Error()})
		}
		var err errors.Error
		var result commands.CheckResponse
		result, err = development.Check(context.Background(), dev.P, params.Subject, params.Action, params.Entity, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{false, err.Error()})
		}
		return js.ValueOf([]interface{}{result.Can, nil})
	})
}

// lookupQuery -
func lookupQuery() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &development.LookupQueryQuery{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), params)
		if mErr != nil {
			return js.ValueOf([]interface{}{"", []interface{}{}, mErr.Error()})
		}
		var err errors.Error
		var result commands.LookupQueryResponse
		result, err = development.LookupQuery(context.Background(), dev.P, params.EntityType, params.Action, params.Subject, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{"", []interface{}{}, err.Error()})
		}
		return js.ValueOf([]interface{}{result.Query, result.Args, nil})
	})
}

// writeSchema -
func writeSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var err errors.Error
		var version string
		version, err = development.WriteSchema(context.Background(), dev.M, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{"", err.Error()})
		}
		return js.ValueOf([]interface{}{version, nil})
	})
}

// writeTuple -
func writeTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		t := &base.Tuple{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), t)
		if mErr != nil {
			return js.ValueOf([]interface{}{mErr.Error()})
		}
		var err errors.Error
		err = development.WriteTuple(context.Background(), dev.R, t, string(args[1].String()))
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		return js.ValueOf([]interface{}{nil})
	})
}

// deleteTuple -
func deleteTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		t := &base.Tuple{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), t)
		if mErr != nil {
			return js.ValueOf([]interface{}{mErr.Error()})
		}
		var err errors.Error
		err = development.DeleteTuple(context.Background(), dev.R, t)
		if err != nil {
			return js.ValueOf([]interface{}{err.Error()})
		}
		return js.ValueOf([]interface{}{nil})
	})
}

// readSchema -
func readSchema() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var sch schema.Schema
		var err errors.Error
		sch, err = development.ReadSchema(context.Background(), dev.M, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		result, mErr := json.Marshal(sch)
		if mErr != nil {
			return js.ValueOf([]interface{}{nil, mErr.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

// readTuple -
func readTuple() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := &base.TupleFilter{}
		mErr := json.Unmarshal([]byte(string(args[0].String())), params)
		if mErr != nil {
			return js.ValueOf([]interface{}{false, mErr.Error()})
		}
		var tuples tuple.ITupleCollection
		var err errors.Error
		tuples, err = development.ReadTuple(context.Background(), dev.R, params)
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		var result []byte
		result, mErr = json.Marshal(tuples)
		if mErr != nil {
			return js.ValueOf([]interface{}{nil, mErr.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

// readSchemaGraph -
func readSchemaGraph() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var sch schema.Schema
		var err errors.Error
		sch, err = development.ReadSchema(context.Background(), dev.M, string(args[0].String()))
		if err != nil {
			return js.ValueOf([]interface{}{nil, err.Error()})
		}
		r, gErr := schema.GraphSchema(sch)
		if gErr != nil {
			return js.ValueOf([]interface{}{nil, gErr.Error()})
		}
		result, mErr := json.Marshal(struct {
			Nodes []*graph.Node `json:"nodes"`
			Edges []*graph.Edge `json:"edges"`
		}{Nodes: r.Nodes(), Edges: r.Edges()})
		if mErr != nil {
			return js.ValueOf([]interface{}{nil, mErr.Error()})
		}
		return js.ValueOf([]interface{}{string(result), nil})
	})
}

func main() {
	ch := make(chan struct{}, 0)
	dev = development.NewDevelopment()
	js.Global().Set("check", check())
	js.Global().Set("writeSchema", writeSchema())
	js.Global().Set("writeTuple", writeTuple())
	js.Global().Set("readSchema", readSchema())
	js.Global().Set("readTuple", readTuple())
	js.Global().Set("deleteTuple", deleteTuple())
	js.Global().Set("readSchemaGraph", readSchemaGraph())
	js.Global().Set("lookupQuery", lookupQuery())
	<-ch
}