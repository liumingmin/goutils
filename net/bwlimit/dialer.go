package bwlimit

import (
	"context"
	"net"
	"time"
)

// Dialer structure contains default dialer and timeout
type Dialer struct {
	net.Dialer
	rxBwLimit *BwLimit
	txBwLimit *BwLimit

	Timeout time.Duration
}

// NewDialer creates a Dialer structure with Timeout, Keepalive
// bandwidth limit rx: read, tx write
func NewDialer() *Dialer {
	dialer := &Dialer{
		Dialer:    net.Dialer{},
		rxBwLimit: &BwLimit{},
		txBwLimit: &BwLimit{},

		Timeout: time.Minute * 10,
	}
	return dialer
}

// Dial connects to the network address.
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext connects to the network address using the provided context.
func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		return c, err
	}

	con := &timeoutConn{
		Conn:      c,
		timeout:   d.Timeout,
		rxBwLimit: d.rxBwLimit,
		txBwLimit: d.txBwLimit,
	}
	return con, con.nudgeDeadline()
}

func (d *Dialer) RxBwLimit() *BwLimit {
	return d.rxBwLimit
}

func (d *Dialer) TxBwLimit() *BwLimit {
	return d.txBwLimit
}

// A net.Conn that sets deadline for every Read/Write operation
type timeoutConn struct {
	net.Conn
	timeout   time.Duration
	rxBwLimit *BwLimit
	txBwLimit *BwLimit
}

// Nudge the deadline for an idle timeout on by c.timeout if non-zero
func (c *timeoutConn) nudgeDeadline() error {
	if c.timeout > 0 {
		return c.SetDeadline(time.Now().Add(c.timeout))
	}
	return nil
}

// Read bytes with rate limiting and idle timeouts
func (c *timeoutConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.rxBwLimit.LimitBandwidth(n)
	if err == nil && n > 0 && c.timeout > 0 {
		err = c.nudgeDeadline()
	}
	return n, err
}

// Write bytes with rate limiting and idle timeouts
func (c *timeoutConn) Write(b []byte) (n int, err error) {
	c.txBwLimit.LimitBandwidth(len(b))
	n, err = c.Conn.Write(b)
	if err == nil && n > 0 && c.timeout > 0 {
		err = c.nudgeDeadline()
	}
	return n, err
}
