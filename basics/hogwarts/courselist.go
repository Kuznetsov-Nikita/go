//go:build !solution

package hogwarts

import "slices"

func dfs(vertice string, graph map[string][]string, vertices map[string]bool, stack *[]string, visited *map[string]bool) {
	vertices[vertice] = true
	(*visited)[vertice] = true

	if _, ok := graph[vertice]; ok == true {
		for _, value := range graph[vertice] {
			if (*visited)[value] == true {
				panic("error: cyclic dependency")
			}
			if vertices[value] == false {
				dfs(value, graph, vertices, stack, visited)
			}
		}
	}

	*stack = append(*stack, vertice)
}

func GetCourseList(prereqs map[string][]string) []string {
	var vertices = make(map[string]bool)
	var graph = make(map[string][]string)

	for key, value := range prereqs {
		vertices[key] = false
		for _, elem := range value {
			vertices[elem] = false
			graph[elem] = append(graph[elem], key)
		}
	}

	var stack []string

	for key := range vertices {
		if vertices[key] == false {
			var visited = make(map[string]bool)
			dfs(key, graph, vertices, &stack, &visited)
		}
	}

	slices.Reverse(stack)
	return stack
}
