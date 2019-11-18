package container

import (
	"fmt"
)

type Container struct {
	registeredComponents map[string]func() interface{}

	instances map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		registeredComponents: map[string]func() interface{}{},
		instances: map[string]interface{}{},
	}
}

func (con *Container) Register(componentId string, factory func() interface{}) {
	con.registeredComponents[componentId] = factory
}

func (con *Container) Create(componentId string) (interface{}, error) {
	component, ok := con.registeredComponents[componentId]
	if !ok {
		return nil, fmt.Errorf("Unregistered component \"%s\".", componentId)
	}

	return component(), nil
}

func (con *Container) Get(componentId string) (interface{}, error) {
	if _, ok := con.instances[componentId]; !ok {
		instance, err := con.Create(componentId)
		if (err != nil) {
			return nil, err
		}

		con.instances[componentId] = instance
	}

	return con.instances[componentId], nil
}
