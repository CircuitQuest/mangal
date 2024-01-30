package config

import (
	"reflect"

	"github.com/spf13/viper"
)

// Fields are all the registered fields.
var Fields = make(map[string]entry)

type entry struct {
	Key         string
	Description string
	Default     any
	Marshal     func(any) (any, error)
	Unmarshal   func(any) (any, error)
	Validate    func(any) error
	SetValue    func(any) error
}

type field[Raw, Value any] struct {
	Key         string
	Description string
	Default     Value
	Unmarshal   func(Raw) (Value, error)
	Marshal     func(Value) (Raw, error)
	Validate    func(Value) error
}

type registered[Raw, Value any] struct {
	field[Raw, Value]
	value Value
}

func (r *registered[Raw, Value]) Get() Value {
	return r.value
}

// Set the value for the registered entry and viper.
func (r *registered[Raw, Value]) Set(value Value) error {
	marshalled, err := r.Marshal(value)
	if err != nil {
		return err
	}
	r.value = value
	viper.Set(r.Key, marshalled)

	return nil
}

// Register a new config entry.
func reg[Raw, Value any](f field[Raw, Value]) *registered[Raw, Value] {
	if f.Marshal == nil {
		f.Marshal = func(value Value) (raw Raw, err error) {
			return reflect.
				ValueOf(value).
				Convert(reflect.ValueOf(raw).Type()).
				Interface().(Raw), nil
		}
	}
	if f.Unmarshal == nil {
		f.Unmarshal = func(raw Raw) (value Value, err error) {
			return reflect.
				ValueOf(raw).
				Convert(reflect.ValueOf(value).Type()).
				Interface().(Value), nil
		}
	}
	if f.Validate == nil {
		f.Validate = func(Value) error {
			return nil
		}
	}

	r := &registered[Raw, Value]{
		field: f,
		value: f.Default,
	}

	// Add the entry to the exported Fields
	Fields[f.Key] = entry{
		Key:         f.Key,
		Description: f.Description,
		Default:     f.Default,
		Marshal: func(a any) (any, error) {
			return f.Marshal(a.(Value))
		},
		Unmarshal: func(a any) (any, error) {
			return f.Unmarshal(a.(Raw))
		},
		Validate: func(a any) error {
			return f.Validate(a.(Value))
		},
		SetValue: func(a any) error {
			return r.Set(a.(Value))
		},
	}

	return r
}
