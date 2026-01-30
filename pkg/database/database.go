package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type DB struct {
	conn *gorm.DB
}

func NewDB(host, user, password, dbname string, port int, sslmode string, maxOpenConns, maxIdleConns, connMaxLifetime int) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	db.LogMode(true)

	log.Println("Database connected successfully")

	return &DB{conn: db}, nil
}

func (db *DB) Close() error {
	return db.conn.DB().Close()
}

func (db *DB) Transaction(fn func(tx *DB) error) error {
	if db == nil || db.conn == nil {
		return fmt.Errorf("database not initialized")
	}
	transaction := db.conn.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	wrapped := &DB{conn: transaction}
	if err := fn(wrapped); err != nil {
		transaction.Rollback()
		return err
	}
	if err := transaction.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (db *DB) Ping() error {
	return db.conn.DB().Ping()
}

func (db *DB) Preload(query string, args ...interface{}) *gorm.DB {
	return db.conn.Preload(query, args...)
}

func (db *DB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return db.conn.Where(query, args...)
}

func (db *DB) First(value interface{}, conds ...interface{}) *gorm.DB {
	return db.conn.First(value, conds...)
}

func (db *DB) Find(value interface{}, conds ...interface{}) *gorm.DB {
	return db.conn.Find(value, conds...)
}

func (db *DB) Create(value interface{}) *gorm.DB {
	return db.conn.Create(value)
}

func (db *DB) Save(value interface{}) *gorm.DB {
	return db.conn.Save(value)
}

func (db *DB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return db.conn.Delete(value, conds...)
}

func (db *DB) Error() error {
	return db.conn.Error
}

func (db *DB) AutoMigrate(values ...interface{}) *gorm.DB {
	return db.conn.AutoMigrate(values...)
}

func (db *DB) Order(value interface{}) *gorm.DB {
	return db.conn.Order(value)
}

func (db *DB) Model(value interface{}) *gorm.DB {
	return db.conn.Model(value)
}

func (db *DB) Offset(offset interface{}) *gorm.DB {
	return db.conn.Offset(offset)
}

func (db *DB) Limit(limit interface{}) *gorm.DB {
	return db.conn.Limit(limit)
}

func (db *DB) Count(value interface{}) *gorm.DB {
	return db.conn.Count(value)
}

func (db *DB) Exec(sql string, values ...interface{}) *gorm.DB {
	return db.conn.Exec(sql, values...)
}

func (db *DB) Update(attrs ...interface{}) *gorm.DB {
	return db.conn.Update(attrs...)
}

func (db *DB) Updates(attrs interface{}) *gorm.DB {
	return db.conn.Updates(attrs)
}
