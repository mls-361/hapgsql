/*
------------------------------------------------------------------------------------------------------------------------
####### hapgsql ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package hapgsql

import (
	"time"

	"github.com/mls-361/logger"
)

// ClusterOption AFAIRE.
type ClusterOption func(*Cluster)

// WithLogger AFAIRE.
func WithLogger(l logger.Logger) ClusterOption {
	return func(c *Cluster) {
		c.logger = l
	}
}

// WithUpdateInterval AFAIRE.
func WithUpdateInterval(d time.Duration) ClusterOption {
	return func(c *Cluster) {
		c.updateInterval = d
	}
}

// WithUpdateTimeout AFAIRE.
func WithUpdateTimeout(d time.Duration) ClusterOption {
	return func(c *Cluster) {
		c.updateTimeout = d
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
