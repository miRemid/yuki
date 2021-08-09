package main

func main() {
	g, err := NewGateway(":8080")
	if err != nil {
		panic(err)
	}
	if err := g.ListenAndServe(); err != nil {
		panic(err)
	}
}
