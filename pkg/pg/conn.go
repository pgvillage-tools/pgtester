// Package pg is the module that can be used for communication with postgres
package pg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

// Conn is a struct that can contain information about the connection (for example the delay)
type Conn struct {
	connParams Dsn
	conn       *pgx.Conn
	retries    uint
	delay      time.Duration
}

// NewConn is a function that creates a new connection with given information about the connection
func NewConn(connParams Dsn, retries uint, delay time.Duration) (c *Conn) {
	return &Conn{
		connParams: connParams,
		retries:    retries,
		delay:      delay,
	}
}

// DSN is function that turns connection parameters into a dsn string
func (c *Conn) DSN() (dsn string) {
	var pairs []string
	for key, value := range c.connParams {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, connectStringValue(value)))
	}
	return strings.Join(pairs[:], " ")
}

// Connect is a function that uses the conn object to connect you to postgres
func (c *Conn) Connect() (err error) {
	const zero = 0
	if c.conn != nil {
		if !c.conn.IsClosed() {
			log.Debugf("already connected")
			return nil
		}
		c.conn.IsClosed()
	}
	dsn := c.DSN()
	log.Debugf("connecting to %s", dsn)
	for i := zero; i <= int(c.retries); i++ {
		c.conn, err = pgx.Connect(context.Background(), dsn)
		if err == nil {
			log.Debugf("successfully connected")
			return nil
		}
		c.conn = nil
		log.Debugf("error while connecting: %s", err.Error())
		log.Debugf("waiting %s seconds before trying again", c.delay)
		time.Sleep(c.delay)
	}
	return fmt.Errorf("number of connection retries (%d) exceeded", c.retries)
}

// RunQueryGetOneField is a function that allows you to send querys to postgres and returns the result.
func (c *Conn) RunQueryGetOneField(query string, args ...any) (result Results, err error) {
	var fieldDescriptions []string
	err = c.Connect()
	if err != nil {
		return Results{}, err
	}

	log.Debugf("running query %s with arguments %e", query, args)
	rows, err := c.conn.Query(context.Background(), query, args...)
	if err != nil {
		if closeErr := c.conn.Close(context.Background()); closeErr != nil {
			log.Fatal("Error on query, and I failed to close the connection.")
		}
		return result, err
	} else if rows.Err() != nil {
		return result, err
	}
	for _, fd := range rows.FieldDescriptions() {
		fieldDescriptions = append(fieldDescriptions, string(fd.Name))
	}
	for rows.Next() {
		if rows.Err() != nil {
			return result, err
		}
		values, err := rows.Values()
		if err != nil {
			return result, err
		}
		ofr, err := newResultFromByteArrayArray(fieldDescriptions, values)
		if err != nil {
			return result, err
		}
		result = append(result, ofr)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}
