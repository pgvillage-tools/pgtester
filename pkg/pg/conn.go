package pg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

type Conn struct {
	connParams Dsn
	conn       *pgx.Conn
	retries    uint
	delay      time.Duration
}

func NewConn(connParams Dsn, retries uint, delay time.Duration) (c *Conn) {
	return &Conn{
		connParams: connParams,
		retries:    retries,
		delay:      delay,
	}
}

func (c *Conn) DSN() (dsn string) {
	var pairs []string
	for key, value := range c.connParams {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, connectStringValue(value)))
	}
	return strings.Join(pairs[:], " ")
}

func (c *Conn) Connect() (err error) {
	if c.conn != nil {
		if c.conn.IsClosed() {
			c.conn = nil
		} else {
			log.Debugf("already connected")
			return
		}
	}
	dsn := c.DSN()
	log.Debugf("connecting to %s", dsn)
	for i := 0; i <= int(c.retries); i++ {
		c.conn, err = pgx.Connect(context.Background(), dsn)
		if err == nil {
			log.Debugf("succesfully connected")
			return
		}
		c.conn = nil
		log.Debugf("error while connecting: %s", err.Error())
		log.Debugf("waiting %s seconds before trying again", c.delay)
		time.Sleep(c.delay)
	}
	return fmt.Errorf("number of connection retries (%d) exceeded", c.retries)
}

func (c *Conn) RunQueryGetOneField(query string, args ...interface{}) (result Results, err error) {
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
		ofr, err := NewResultFromByteArrayArray(fieldDescriptions, values)
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
