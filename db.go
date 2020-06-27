package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNoRows = sql.ErrNoRows
)

const (
	DEFAULT_CONNECTION_MAX_LIFE_TIME = -1
	DEFAULT_MAX_OPEN_CONNECTIONS = 80
	DEFAULT_MAX_IDLE_CONNECTIONS = 20
)

type DbConnection struct {
	Host                  string
	Port                  int
	User                  string
	Password              string
	DbName                string
	ConnectionMaxLifeTime int
	MaxOpenConnections    int
	MaxIdleConnections    int
}

type Db struct {
	connections  map[string]DbConnection
	links        map[string]*sqlx.DB
	Rows         *sqlx.Rows
	RowsAffected uint64
	LastInsertId uint64
}

type SQLParams map[string]interface{}

func NewDB() *Db {
	db := new(Db)
	db.init()
	return db
}

func (db *Db) init() {
	db.connections = make(map[string]DbConnection, 10)
	db.links = make(map[string]*sqlx.DB, 10)
}

func (db *Db) AddConnection(Name string, conn DbConnection) (err error) {
	if conn.ConnectionMaxLifeTime == 0 {
		conn.ConnectionMaxLifeTime = DEFAULT_CONNECTION_MAX_LIFE_TIME
	}
	if conn.MaxOpenConnections == 0 {
		conn.MaxOpenConnections = DEFAULT_MAX_OPEN_CONNECTIONS
	}
	if conn.MaxIdleConnections == 0 {
		conn.MaxIdleConnections = DEFAULT_MAX_IDLE_CONNECTIONS
	}
	db.connections[Name] = conn
	return nil
}

func (db *Db) Close() {
	for _, link := range db.links {
		link.Close()
	}
}

func (db *Db) Free() error {
	if db.Rows != nil {
		return db.Rows.Close()
	}
	return nil
}

func (db *Db) Query(link_name string, query_sql string, params SQLParams) (err error) {
	link, err := db.getLink(link_name)
	if err != nil {
		return err
	}

	if len(params) > 0 {
		var query_params = make(map[string]interface{}, len(params))
		for key, value := range params {
			query_params[key] = value
		}
		db.Rows, err = link.NamedQuery(query_sql, query_params)
		if strings.Index(query_sql, "RETURNING") != -1 {
			if db.Rows != nil && db.Rows.Next() {
				db.Rows.Scan(&db.LastInsertId)
			} else {
				return fmt.Errorf("Error getting last insert ID")
			}
		}
	} else {
		db.Rows, err = link.Queryx(query_sql)
	}
	if err != nil {
		return fmt.Errorf("%v (%s)", err, query_sql)
	}

	return nil
}

func (db *Db) QueryRow(link_name string, query_sql string, result interface{}) (err error) {
	link, err := db.getLink(link_name)
	if err != nil {
		return err
	}

	err = link.Get(result, query_sql)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRows
		} else {
			return fmt.Errorf("%v (%s)", err, query_sql)
		}
	}

	return nil
}

func (db *Db) Exec(link_name string, query_sql string, params SQLParams) (err error) {
	link, err := db.getLink(link_name)
	if err != nil {
		return err
	}

	var res sql.Result

	if len(params) > 0 {
		var query_params = make(map[string]interface{}, len(params))
		for key, value := range params {
			query_params[key] = value
		}
		res, err = link.NamedExec(query_sql, query_params)
	} else {
		res, err = link.Exec(query_sql)
	}
	if err != nil {
		return fmt.Errorf("%v (%s)", err, query_sql)
	}

	id, err := res.LastInsertId()
	if err == nil {
		db.LastInsertId = uint64(id)
	}

	affected, err := res.RowsAffected()
	if err == nil {
		db.RowsAffected = uint64(affected)
	}

	return err
}

func (db *Db) getLink(name string) (link *sqlx.DB, err error) {
	if link, found := db.links[name]; found {
		return link, nil
	}

	connection, found := db.connections[name]
	if !found {
		return link, nil
	}

	if len(connection.Host) == 0 {
		connection.Host = "127.0.0.1"
	}
	if connection.Port == 0 {
		connection.Port = 5432
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", connection.Host, connection.Port, connection.User, connection.Password, connection.DbName)
	link, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return link, err
	}

	link.SetConnMaxLifetime(time.Minute * time.Duration(connection.ConnectionMaxLifeTime))
	link.SetMaxOpenConns(connection.MaxOpenConnections)
	link.SetMaxIdleConns(connection.MaxIdleConnections)

	db.links[name] = link

	return link, nil
}

func (db *Db) GetInsertKeys(params SQLParams) (fields []string, values []string) {
	for key := range params {
		fields = append(fields, key)
		values = append(values, ":" + key)
	}
	return
}

func (db *Db) GetUpdateKeys(params SQLParams) (update []string) {
	for key := range params {
		update = append(update, key + "=:" + key)
	}
	return
}

func (db *Db) UInt32Slice2String(slice []uint32) (result string) {
	for idx, value := range slice {
		if idx > 0 {
			result += ","
		}
		result += strconv.FormatUint(uint64(value), 10)
	}
	return
}
