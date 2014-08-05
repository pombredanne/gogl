package al

import (
	. "github.com/sdboyer/gogl"
)

// Contains behaviors shared across adjacency list implementations.

type al_graph interface {
	Graph
	ensureVertex(...Vertex)
	hasVertex(Vertex) bool
}

type al_digraph interface {
	al_graph
	IncidentArcEnumerator
	DirectedDegreeChecker
	Transposer
}

type al_ea interface {
	al_graph
	addEdges(...Edge)
}

type al_wea interface {
	al_graph
	addEdges(...WeightedEdge)
}

type al_lea interface {
	al_graph
	addEdges(...LabeledEdge)
}

type al_pea interface {
	al_graph
	addEdges(...DataEdge)
}

// Copies an incoming graph into any of the implemented adjacency list types.
//
// This encapsulates the full matrix of conversion possibilities between
// different graph edge types.
func functorToAdjacencyList(from GraphSource, to interface{}) Graph {
	vf := func(from GraphSource, to al_graph) {
		if Order(to) != Order(from) {
			from.EachVertex(func(vertex Vertex) (terminate bool) {
				to.ensureVertex(vertex)
				return
			})
		}
	}

	if g, ok := to.(al_ea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			g.addEdges(edge)
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_wea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(WeightedEdge); ok {
				g.addEdges(e)
			} else {
				g.addEdges(NewWeightedEdge(edge.Source(), edge.Target(), 0))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_lea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(LabeledEdge); ok {
				g.addEdges(e)
			} else {
				g.addEdges(NewLabeledEdge(edge.Source(), edge.Target(), ""))
			}
			return
		})
		vf(from, g)
	} else if g, ok := to.(al_pea); ok {
		from.EachEdge(func(edge Edge) (terminate bool) {
			if e, ok := edge.(DataEdge); ok {
				g.addEdges(e)
			} else {
				g.addEdges(NewDataEdge(edge.Source(), edge.Target(), nil))
			}
			return
		})
		vf(from, g)
	} else {
		panic("Target graph did not implement a recognized adjacency list internal type")
	}

	return to.(Graph)
}

func eachVertexInAdjacencyList(list interface{}, vertex Vertex, vs VertexStep) {
	switch l := list.(type) {
	case map[Vertex]map[Vertex]struct{}:
		if _, exists := l[vertex]; exists {
			for adjacent, _ := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]float64:
		if _, exists := l[vertex]; exists {
			for adjacent, _ := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]string:
		if _, exists := l[vertex]; exists {
			for adjacent, _ := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	case map[Vertex]map[Vertex]interface{}:
		if _, exists := l[vertex]; exists {
			for adjacent, _ := range l[vertex] {
				if vs(adjacent) {
					return
				}
			}
		}
	default:
		panic("Unrecognized adjacency list map type.")
	}
}

func eachPredecessorOf(list interface{}, vertex Vertex, vs VertexStep) {
	switch l := list.(type) {
	case map[Vertex]map[Vertex]struct{}:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target, _ := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]float64:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target, _ := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]string:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target, _ := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	case map[Vertex]map[Vertex]interface{}:
		if _, exists := l[vertex]; exists {
			for candidate, adjacent := range l {
				for target, _ := range adjacent {
					if target == vertex {
						if vs(candidate) {
							return
						}
					}
				}
			}
		}
	default:
		panic("Unrecognized adjacency list map type.")
	}

}

func inDegreeOf(g al_graph, v Vertex) (degree int, exists bool) {
	if exists = g.hasVertex(v); exists {
		g.EachEdge(func(e Edge) (terminate bool) {
			if v == e.Target() {
				degree++
			}
			return
		})
	}
	return
}

func eachEdgeIncidentToDirected(g al_digraph, v Vertex, f EdgeStep) {
	if !g.hasVertex(v) {
		return
	}

	var terminate bool
	interloper := func(e Edge) bool {
		terminate = terminate || f(e)
		return terminate
	}

	g.EachArcFrom(v, interloper)
	g.EachArcTo(v, interloper)
}