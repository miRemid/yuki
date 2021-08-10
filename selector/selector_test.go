package selector_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/miRemid/yuki/selector"
)

const (
	replicas = 9
	nodeNums = 5
)

func generateNode() []*selector.Node {
	nodes := make([]*selector.Node, 0)
	for i := 1; i <= nodeNums; i++ {
		nodes = append(nodes, &selector.Node{
			RemoteAddr: fmt.Sprintf("%d", i),
		})
	}
	return nodes
}

func testSelector(s selector.Selector) {
	for i := 0; i < replicas; i++ {
		node, _ := s.Peek(fmt.Sprintf("%d", i+rand.Intn(40)))
		log.Println(node.RemoteAddr)
	}
}

func TestRandomSelector(t *testing.T) {
	nodes := generateNode()
	s := selector.NewRandomSelector(nodes...)
	testSelector(s)
}

func TestRoundRobinSelector(t *testing.T) {
	nodes := generateNode()
	s := selector.NewRoundRobinSelector(nodes...)
	testSelector(s)
}

func TestWeightSelector(t *testing.T) {
	nodes := generateNode()
	s := selector.NewWeightRoundRobinSelector(nodes...)
	testSelector(s)
}

func TestHashSelector(t *testing.T) {
	nodes := generateNode()
	s := selector.NewHashSelector(nodes...)
	testSelector(s)
}

func TestHashDelete(t *testing.T) {
	nodes := generateNode()
	s := selector.NewHashSelector(nodes...)
	ns, _ := s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	if err := s.Delete("5"); err != nil {
		log.Fatal(err)
	}
	log.Println("=========DELETE========")
	ns, _ = s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
}
func TestWeightDelete(t *testing.T) {
	nodes := generateNode()
	s := selector.NewWeightRoundRobinSelector(nodes...)
	ns, _ := s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	if err := s.Delete("3"); err != nil {
		log.Fatal(err)
	}
	log.Println("=========DELETE========")
	ns, _ = s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	if err := s.Delete("5"); err != nil {
		log.Fatal(err)
	}
	log.Println("=========DELETE========")
	ns, _ = s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	if err := s.Delete("1"); err != nil {
		log.Fatal(err)
	}
	log.Println("=========DELETE========")
	ns, _ = s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	s.Add(&selector.Node{
		RemoteAddr: "888",
	})
	log.Println("===========ADD==========")
	ns, _ = s.Getall()
	for _, n := range ns {
		log.Println(n.RemoteAddr)
	}
	testSelector(s)
}
