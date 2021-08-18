package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strings"
)

type Conn struct {
	connParams Dsn
	conn       *pgx.Conn
}

func NewConn(connParams Dsn) (c *Conn) {
	return &Conn{
		connParams: connParams,
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
			return nil
		}
	}
	dsn := c.DSN()
	log.Debugf("connecting to %s", dsn)
	c.conn, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		c.conn = nil
		log.Debugf("error while connecting: %e", err)
		return err
	}
	log.Debugf("succesfully connected")
	return nil
}

func (c *Conn) RunQueryGetOneField(query string, args ...interface{}) (result OneFieldResults, err error) {
	var fieldDescriptions []string
	err = c.Connect()
	if err != nil {
		return result, err
	}

	log.Debugf("running query %s with arguments %e", query, args)
	rows, err := c.conn.Query(context.Background(), query, args...)
	if err != nil {
		return result, err
	}
	for _, fd := range rows.FieldDescriptions() {
		fieldDescriptions = append(fieldDescriptions, string(fd.Name))
	}
	for {
		if ! rows.Next() {
			break
		}
		if rows.Err() != nil {
			return result, err
		}
		values, err := rows.Values()
		if err != nil {
			return result, err
		}
		ofr, err := NewOneFieldResultFromByteArrayArray(fieldDescriptions, values)
		if err != nil {
			return result, err
		}
		result = append(result, ofr)
	}
	return result, nil
}
