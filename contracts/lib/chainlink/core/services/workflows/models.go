package workflows

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/dominikbraun/graph"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

// LegacyWorkflowName is used for YAML workflows only.
// It has to be exactly 10 bytes long and Hex() encodes it without hashing.
type legacyWorkflowName struct {
	name string
}

func (d legacyWorkflowName) Hex() string {
	return hex.EncodeToString([]byte(d.name))
}

func (d legacyWorkflowName) String() string {
	return d.name
}

func NewLegacyWorkflowName(name string) types.WorkflowName {
	return legacyWorkflowName{name: name}
}

// workflow is a directed graph of nodes, where each node is a step.
//
// triggers are special steps that are stored separately, they're
// treated differently due to their nature of being the starting
// point of a workflow.
type workflow struct {
	id    string
	owner string
	name  types.WorkflowName
	graph.Graph[string, *step]

	triggers []*triggerCapability
}

func (w *workflow) walkDo(start string, do func(s *step) error) error {
	var outerErr error
	err := graph.BFS(w.Graph, start, func(ref string) bool {
		n, err := w.Graph.Vertex(ref)
		if err != nil {
			outerErr = err
			return true
		}

		err = do(n)
		if err != nil {
			outerErr = err
			return true
		}

		return false
	})
	if err != nil {
		return err
	}

	return outerErr
}

// dependents returns all steps that directly depend on the step with the given ref
func (w *workflow) dependents(start string) ([]*step, error) {
	var steps []*step
	m, err := w.Graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	adj, ok := m[start]
	if !ok {
		return nil, fmt.Errorf("could not find step with ref %s", start)
	}

	for adjacentRef := range adj {
		n, err := w.Graph.Vertex(adjacentRef)
		if err != nil {
			return nil, err
		}

		steps = append(steps, n)
	}

	return steps, nil
}

// step wraps a Vertex with additional context for execution that is mutated by the engine
type step struct {
	workflows.Vertex
	capability capabilities.ExecutableCapability
	info       capabilities.CapabilityInfo
	config     *values.Map
}

type triggerCapability struct {
	sdk.StepDefinition
	trigger capabilities.TriggerCapability

	// flag to track registration of the trigger and avoid removal of non registered triggers
	registered bool

	config *values.Map

	mu sync.Mutex
}

func Parse(sdkSpec sdk.WorkflowSpec) (*workflow, error) {
	wf2, err := workflows.BuildDependencyGraph(sdkSpec)
	if err != nil {
		return nil, err
	}

	wfs, err := createWorkflow(wf2)
	return wfs, err
}

// createWorkflow converts a StaticWorkflow to an executable workflow
// by adding metadata to the vertices that is owned by the workflow runtime.
func createWorkflow(wf2 *workflows.DependencyGraph) (*workflow, error) {
	out := &workflow{
		id:       wf2.ID,
		triggers: []*triggerCapability{},
	}

	for _, t := range wf2.Triggers {
		out.triggers = append(out.triggers, &triggerCapability{
			StepDefinition: *t,
		})
	}

	stepHash := func(s *step) string {
		// must use the same hash function as the DependencyGraph.
		// this ensures that the intermediate representation (DependencyGraph) and the workflow
		// representation label vertices with the same identifier, which in turn allows us to
		// to copy the edges from the intermediate representation to the executable representation.
		return s.Vertex.VID()
	}
	g := graph.New(
		stepHash,
		graph.PreventCycles(),
		graph.Directed(),
	)
	adjMap, err := wf2.Graph.AdjacencyMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert intermediate representation to adjacency map: %w", err)
	}

	// copy the all the vertices from the intermediate graph to the executable workflow graph
	for vertexRef := range adjMap {
		v, innerErr := wf2.Graph.Vertex(vertexRef)
		if innerErr != nil {
			return nil, fmt.Errorf("failed to retrieve vertex for %s: %w", vertexRef, innerErr)
		}
		innerErr = g.AddVertex(&step{Vertex: *v})
		if innerErr != nil {
			return nil, fmt.Errorf("failed to add vertex to executable workflow %s: %w", vertexRef, innerErr)
		}
	}
	// now we can add all the edges. this works because we are using vertex hash function is the same in both graphs.
	// see comment on `stepHash` function.
	for vertexRef, edgeRefs := range adjMap {
		for edgeRef := range edgeRefs {
			innerErr := g.AddEdge(vertexRef, edgeRef)
			// If we fail to add the edge, we'll bail out unless we encountered an ErrEdgeAlreadyExists, in which case
			// we'll continue. This is because inputs can contain multiple references to the parent node.
			if innerErr != nil && !errors.Is(innerErr, graph.ErrEdgeAlreadyExists) {
				return nil, fmt.Errorf("failed to add edge from '%s' to '%s': %w", vertexRef, edgeRef, innerErr)
			}
		}
	}
	out.Graph = g
	return out, nil
}
