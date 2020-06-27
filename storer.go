package main

import (
	"context"
	"fmt"
	"github.com/volatiletech/authboss"
)

type User struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Role     int    `db:"role"`
}

type Storer struct {
	config *Config
	db     *Db
	users  map[string]authboss.User
}

var (
	assertUser                 = &User{}
	//assertStorer               = &Storer{}
	_            authboss.User = assertUser
)

// PutPID into user
func (u *User) PutPID(pid string) {}

// PutPassword into user
func (u *User) PutPassword(password string) { u.Password = password }

// PutEmail into user
func (u *User) PutEmail(email string) {}

// GetPID from user
func (u User) GetPID() string { return u.Name }

// GetPassword from user
func (u User) GetPassword() string { return u.Password }

// GetEmail from user
func (u User) GetEmail() string { return "" }

func NewStorer(config *Config) *Storer {
	db := NewDB()
	db.AddConnection(ABSERVER_DB_NAME, DbConnection{
		Host:                  config.DbHost,
		Port:                  config.DbPort,
		User:                  config.DbUser,
		Password:              config.DbPassword,
		DbName:                config.DbName,
		ConnectionMaxLifeTime: config.DbConnectionMaxLifeTime,
		MaxOpenConnections:    config.DbMaxOpenConnections,
		MaxIdleConnections:    config.DbMaxIdleConnections,
	})

	return &Storer{
		config: config,
		db:     db,
		users:  make(map[string]authboss.User, 10),
	}
}

func (s *Storer) Save(_ context.Context, user authboss.User) error {
	return nil
}

func (s *Storer) Load(_ context.Context, key string) (user authboss.User, err error) {
	u, ok := s.users[key]
	if ok {
		return u, nil
	}

	sql := "SELECT id, name, password, role FROM users WHERE name = :name"
	err = s.db.Query(ABSERVER_DB_NAME, sql, SQLParams{"name": key})
	if err != nil {
		return nil, fmt.Errorf("Internal error")
	}
	defer s.db.Free()

	user = &User{}
	if s.db.Rows.Next() {
		err = s.db.Rows.StructScan(user)
		if err == nil {
			s.users[key] = user
		}
	}
	return user, err
}

func (s *Storer) New(_ context.Context) authboss.User {
	return &User{}
}

func (s *Storer) Create(_ context.Context, user authboss.User) error {
	return nil
}
