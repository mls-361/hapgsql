/*
------------------------------------------------------------------------------------------------------------------------
####### hapgsql ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package hapgsql

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mls-361/logger"
)

const (
	_updateInterval = time.Second * 5
	_updateTimeout  = time.Second * 1
)

type (
	aliveNodes map[bool][]*Node

	// Cluster AFAIRE.
	Cluster struct {
		logger         *logger.Logger
		updateInterval time.Duration
		updateTimeout  time.Duration
		nodes          []*Node
		stop           chan struct{}
		waitGroup      sync.WaitGroup
		aliveNodes     atomic.Value
	}
)

// NewCluster AFAIRE.
func NewCluster(options ...ClusterOption) *Cluster {
	cluster := &Cluster{
		updateInterval: _updateInterval,
		updateTimeout:  _updateTimeout,
		nodes:          make([]*Node, 0),
		stop:           make(chan struct{}),
	}

	for _, o := range options {
		o(cluster)
	}

	cluster.aliveNodes.Store(make(aliveNodes))

	return cluster
}

// AddNode AFAIRE.
func (c *Cluster) AddNode(node *Node) {
	c.nodes = append(c.nodes, node)
}

func sortNodes(nodes []*Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].latency < nodes[j].latency
	})
}

func (c *Cluster) checkNodes(ctx context.Context) aliveNodes {
	aNodes := aliveNodes{
		false: make([]*Node, 0), // alive
		true:  make([]*Node, 0), // alive and primary
	}

	var (
		mutex     sync.Mutex
		waitGroup sync.WaitGroup
	)

	waitGroup.Add(len(c.nodes))

	for _, node := range c.nodes {
		go func(node *Node) { //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
			defer waitGroup.Done()

			if err := node.check(ctx, c.logger); err != nil {
				return
			}

			mutex.Lock()
			defer mutex.Unlock()

			if node.primary {
				aNodes[true] = append(aNodes[true], node)
			} else {
				aNodes[false] = append(aNodes[false], node)
			}
		}(node)
	}

	waitGroup.Wait()

	sortNodes(aNodes[false])
	sortNodes(aNodes[true])

	return aNodes
}

func (c *Cluster) updateNodes() {
	ctx, cancel := context.WithTimeout(context.Background(), c.updateTimeout)
	defer cancel()

	c.aliveNodes.Store(c.checkNodes(ctx))
}

func (c *Cluster) update() {
	defer c.waitGroup.Done()

	c.updateNodes()

	ticker := time.NewTicker(c.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			c.updateNodes()
		}
	}
}

// Update AFAIRE.
func (c *Cluster) Update() {
	c.waitGroup.Add(1)
	go c.update() //@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
}

func (c *Cluster) getNode(primary bool) *Node {
	m := c.aliveNodes.Load().(aliveNodes)
	list := m[primary]

	if len(list) == 0 {
		return nil
	}

	return list[0]
}

// Primary AFAIRE.
func (c *Cluster) Primary() *Node {
	return c.getNode(true)
}

// PrimaryPreferred AFAIRE.
func (c *Cluster) PrimaryPreferred() *Node {
	if node := c.Primary(); node != nil {
		return node
	}

	return c.getNode(false)
}

// Close AFAIRE.
func (c *Cluster) Close() {
	close(c.stop)
	c.waitGroup.Wait()

	for _, node := range c.nodes {
		node.Close()
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
