package processes

import (
  "fmt"
  "github.com/MarconiProtocol/cli/core/configs"
  "os"
)

type Graph struct {
  SortedNodes []*Node
  NodesMap    map[string]*Node
}

type Node struct {
  ProcessConfig configs.ProcessConfig
  visited       bool
  done          bool
}

func (g *Graph) getOrderedProcessConfigs() []configs.ProcessConfig {
  procConfigs := []configs.ProcessConfig{}
  for _, node := range g.SortedNodes {
    procConfigs = append(procConfigs, node.ProcessConfig)
  }
  return procConfigs
}

/*
  Builds a dependency graph from config which is topologically sorted to find the correct order for process execution
*/
func buildDependencyGraph(processes []configs.ProcessConfig) *Graph {
  graph := Graph{}
  graph.buildDependencyNodesMap(processes)

  // DFS topological sort
  for _, node := range graph.NodesMap {
    if !node.visited && !node.done {
      graph.visit(node)
    }
  }
  return &graph
}

// dfs visit
func (graph *Graph) visit(n *Node) {
  if n.done {
    return
  }
  if n.visited {
    fmt.Println("Cyclic graph found")
    os.Exit(1)
  }

  // Look to see if we have any other dependencies
  n.visited = true
  for _, dependencyId := range n.ProcessConfig.Dependencies {
    if d, exists := graph.NodesMap[dependencyId]; exists {
      // Visit dependencies
      graph.visit(d)
    }
  }
  n.done = true

  // Add node to sorted dependency nodes list
  graph.SortedNodes = append(graph.SortedNodes, n)
}

// parses ProcessConfigs into Node objects for use in a dependency graph
func (graph *Graph) buildDependencyNodesMap(processes []configs.ProcessConfig) {
  graph.NodesMap = make(map[string]*Node)
  for _, procConfig := range processes {
    graph.NodesMap[procConfig.Id] = &Node{
      procConfig,
      false,
      false,
    }
  }
}
