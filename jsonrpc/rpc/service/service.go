package service

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Service struct {
	methods map[string]*Method
}

func NewService(s any) (*Service, error) {
	service := &Service{
		methods: make(map[string]*Method),
	}

	err := service.register(s)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s *Service) register(rcvr any) error {
	v := reflect.ValueOf(rcvr)
	t := v.Type()

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.PkgPath != "" {
			continue // method not exported
		}

		me, err := NewMethod(m.Name, v.Method(i).Interface())
		if err != nil {
			log.Printf("can't register method %s: %v", m.Name, err)
			continue
		}

		s.methods[m.Name] = me
	}

	return nil
}

func (s *Service) Call(methodName string, raw *json.RawMessage) (results []any, err error) {
	method, exists := s.methods[methodName]
	if !exists {
		return nil, fmt.Errorf("method not found")
	}

	args, err := method.UnmarshalArgs(raw)
	if err != nil {
		return nil, err
	}

	return method.Call(args...)
}
