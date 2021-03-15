/*
------------------------------------------------------------------------------------------------------------------------
####### hapgsql ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package hapgsql

import (
	"context"
	"time"

	"github.com/mls-361/logger"
	"github.com/mls-361/pgsql"
)

type (
	// Node AFAIRE.
	Node struct {
		host    string
		client  *pgsql.Client
		alive   bool
		primary bool
		latency time.Duration
	}
)

// NewNode AFAIRE.
func NewNode(host string, client *pgsql.Client) *Node {
	return &Node{
		host:   host,
		client: client,
		alive:  true,
	}
}

// Host AFAIRE.
func (n *Node) Host() string {
	return n.host
}

// Client AFAIRE.
func (n *Node) Client() *pgsql.Client {
	return n.client
}

func (n *Node) check(ctx context.Context, logger logger.Logger) error {
	t0 := time.Now()
	err := n.client.QueryRow(ctx, "SELECT NOT pg_is_in_recovery()").Scan(&n.primary)
	n.latency = time.Since(t0)

	if err != nil {
		if n.alive {
			if logger != nil {
				logger.Error( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
					"Problem when checking this node",
					"node", n.host,
					"database", n.client.Database(),
					"reason", err.Error(),
				)
			}

			n.alive = false
		}

		return err
	}

	n.alive = true

	if logger != nil {
		logger.Trace( //::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
			"CheckNode",
			"node", n.host,
			"database", n.client.Database(),
			"primary", n.primary,
			"latency", n.latency,
		)
	}

	return nil
}

// Close AFAIRE.
func (n *Node) Close() {
	n.client.Close()
}

/*
######################################################################################################## @(°_°)@ #######
*/
