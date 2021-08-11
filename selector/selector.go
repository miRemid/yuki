package selector

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	SELECTOR_BUCKET      = "selector_bucket"
	SELECTOR_KEY         = "selector_function"
	RANDOM_SELECTOR      = "random"
	ROUND_ROBIN_SELECTOR = "round_robin"
	HASH_SELECTOR        = "hash"
	WEIGHT_SELECTOR      = "weight"
)

var (
	ErrNodesEmpty = errors.New("empty nodes")
)

func ErrNodeNotFound(remote string) error {
	return fmt.Errorf("%s node not found", remote)
}

type Selector interface {
	Add(node *Node) error
	Peek(seed string) (node *Node, err error)
	Getall() (nodes []*Node, err error)
	Check(remote string) error
	Delete(remote string) error
	Name() string
}

type RandomSelector struct {
	mu  sync.RWMutex
	cur int

	nodes     []*Node
	remoteMap map[string]int
}

func NewRandomSelector(nodes ...*Node) *RandomSelector {
	var s = new(RandomSelector)
	s.mu = sync.RWMutex{}
	s.nodes = make([]*Node, 0)
	s.remoteMap = make(map[string]int)
	rand.Seed(time.Now().Unix())
	for i, n := range nodes {
		s.add(n)
		s.remoteMap[n.RemoteAddr] = i
	}
	return s
}

func (r *RandomSelector) Name() string {
	return RANDOM_SELECTOR
}

func (r *RandomSelector) check(remote string) error {
	if _, ok := r.remoteMap[remote]; ok {
		return nil
	}
	return ErrNodeNotFound(remote)
}

func (r *RandomSelector) add(node *Node) error {
	r.nodes = append(r.nodes, node)
	return nil
}

func (r *RandomSelector) Add(node *Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.remoteMap[node.RemoteAddr] = len(r.nodes)
	return r.add(node)
}

func (r *RandomSelector) Delete(remote string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.check(remote); err != nil {
		return err
	}
	index := r.remoteMap[remote]
	r.nodes = append(r.nodes[:index], r.nodes[index+1:]...)
	// modify index -> len(r.nodes) map index data
	for i := index; i < len(r.nodes); i++ {
		r.remoteMap[r.nodes[i].RemoteAddr] = i
	}
	// delete map data
	delete(r.remoteMap, remote)
	return nil
}

func (r *RandomSelector) Peek(seed string) (*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return nil, ErrNodesEmpty
	}
	r.cur = rand.Intn(len(r.nodes))
	return r.nodes[r.cur], nil
}

func (r *RandomSelector) Check(remote string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.check(remote)
}

func (r *RandomSelector) Getall() ([]*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes, nil
}

type RoundRobinSelector struct {
	mu        sync.RWMutex
	cur       int
	nodes     []*Node
	remoteMap map[string]int
}

func NewRoundRobinSelector(nodes ...*Node) *RoundRobinSelector {
	var s = new(RoundRobinSelector)
	s.mu = sync.RWMutex{}
	s.nodes = make([]*Node, 0)
	s.remoteMap = make(map[string]int)
	for i, n := range nodes {
		s.remoteMap[n.RemoteAddr] = i
		s.add(n)
	}
	return s
}

func (r *RoundRobinSelector) Name() string {
	return ROUND_ROBIN_SELECTOR
}

func (r *RoundRobinSelector) add(node *Node) error {
	r.nodes = append(r.nodes, node)
	return nil
}

func (r *RoundRobinSelector) check(remote string) error {
	if _, ok := r.remoteMap[remote]; ok {
		return nil
	}
	return ErrNodeNotFound(remote)
}

func (r *RoundRobinSelector) Add(node *Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.remoteMap[node.RemoteAddr] = len(r.nodes)
	return r.add(node)
}

func (r *RoundRobinSelector) Check(remote string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.check(remote)
}

func (r *RoundRobinSelector) Delete(remote string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.check(remote); err != nil {
		return err
	}
	index := r.remoteMap[remote]
	r.nodes = append(r.nodes[:index], r.nodes[index+1:]...)
	// modify index -> len(r.nodes) map index data
	for i := index; i < len(r.nodes); i++ {
		r.remoteMap[r.nodes[i].RemoteAddr] = i
	}
	// delete map data
	delete(r.remoteMap, remote)
	return nil
}

func (r *RoundRobinSelector) Peek(seed string) (*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return nil, ErrNodesEmpty
	}
	lens := len(r.nodes)
	if r.cur >= lens {
		r.cur = 0
	}
	node := r.nodes[r.cur]
	r.cur = (r.cur + 1) % lens
	return node, nil
}

func (r *RoundRobinSelector) Getall() ([]*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes, nil
}

type WeightRoundRobinSelector struct {
	mu        sync.RWMutex
	nodes     []*Node
	remoteMap map[string]int
}

func NewWeightRoundRobinSelector(nodes ...*Node) *WeightRoundRobinSelector {
	var s = new(WeightRoundRobinSelector)
	s.mu = sync.RWMutex{}
	s.nodes = make([]*Node, 0)
	s.remoteMap = make(map[string]int)
	for i, n := range nodes {
		s.add(n)
		s.remoteMap[n.RemoteAddr] = i
	}
	return s
}

func (r *WeightRoundRobinSelector) Name() string {
	return WEIGHT_SELECTOR
}

func (r *WeightRoundRobinSelector) check(remote string) error {
	if _, ok := r.remoteMap[remote]; ok {
		return nil
	}
	return ErrNodeNotFound(remote)
}

func (r *WeightRoundRobinSelector) add(node *Node) error {
	node.weight = len(r.nodes)
	node.effectiveWeight = node.weight
	r.nodes = append(r.nodes, node)
	return nil
}

func (r *WeightRoundRobinSelector) Add(node *Node) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.remoteMap[node.RemoteAddr] = len(r.nodes)
	return r.add(node)
}

func (r *WeightRoundRobinSelector) Check(remote string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.check(remote)
}

func (r *WeightRoundRobinSelector) Delete(remote string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.check(remote); err != nil {
		return err
	}
	index := r.remoteMap[remote]
	r.nodes = append(r.nodes[:index], r.nodes[index+1:]...)
	// modify index -> len(r.nodes) map index data
	for i := index; i < len(r.nodes); i++ {
		r.remoteMap[r.nodes[i].RemoteAddr] = i
	}
	// delete map data
	delete(r.remoteMap, remote)
	return nil
}

func (r *WeightRoundRobinSelector) Peek(seed string) (*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return nil, ErrNodesEmpty
	}
	var best *Node
	total := 0
	for i := 0; i < len(r.nodes); i++ {
		w := r.nodes[i]
		total += w.effectiveWeight
		w.currentWeight += w.effectiveWeight
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	best.currentWeight -= total
	return best, nil
}

func (r *WeightRoundRobinSelector) Getall() ([]*Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.nodes, nil
}

type HashFunc func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type HashSelector struct {
	mu        sync.RWMutex
	hash      HashFunc
	replicas  int
	keys      UInt32Slice
	nodes     map[uint32]*Node
	nodesArr  []*Node
	remoteMap map[string]int // index
}

func NewHashSelector(nodes ...*Node) *HashSelector {
	var s = new(HashSelector)
	s.nodesArr = append([]*Node{}, nodes...)
	s.hash = crc32.ChecksumIEEE
	s.replicas = 5
	s.mu = sync.RWMutex{}
	s.reset()
	return s
}

func (c *HashSelector) Name() string {
	return HASH_SELECTOR
}

func (c *HashSelector) reset() {
	c.nodes = make(map[uint32]*Node)
	c.keys = make(UInt32Slice, 0)
	c.remoteMap = make(map[string]int)
	for i, n := range c.nodesArr {
		c.remoteMap[n.RemoteAddr] = i
		c.Add(n)
	}
}

func (c *HashSelector) add(node *Node) error {
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + node.RemoteAddr))
		c.keys = append(c.keys, hash)
		node.keys = append(node.keys, hash)
		c.nodes[hash] = node
	}
	sort.Sort(c.keys)
	return nil
}

func (c *HashSelector) check(remote string) error {
	if _, ok := c.remoteMap[remote]; ok {
		return nil
	}
	return ErrNodeNotFound(remote)
}

func (c *HashSelector) isEmpty() bool {
	return c.keys.Len() == 0
}

func (c *HashSelector) Add(node *Node) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.add(node)
}

func (c *HashSelector) Delete(remote string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.check(remote); err != nil {
		return err
	}
	index := c.remoteMap[remote]
	c.nodesArr = append(c.nodesArr[:index], c.nodesArr[index+1:]...)
	c.reset()
	return nil
}

func (c *HashSelector) Check(remote string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.check(remote)
}

func (c *HashSelector) Peek(seed string) (*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.isEmpty() {
		return nil, ErrNodesEmpty
	}
	hash := c.hash([]byte(seed))
	idx := sort.Search(c.keys.Len(), func(i int) bool {
		return c.keys[i] >= hash
	})
	if idx == c.keys.Len() {
		idx = 0
	}
	return c.nodes[c.keys[idx]], nil
}

func (c *HashSelector) Getall() ([]*Node, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.nodesArr, nil
}
