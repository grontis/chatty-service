package main

type Datastore interface {
	GetAll() ([]*User, error)
	GetByID(id string) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id string) error
}
