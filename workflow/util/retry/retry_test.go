package retry

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func TestRemoveDuplicates(t *testing.T) {
	t.Run("EmptySlice", func(t *testing.T) {
		assert.Equal(t, []string{}, RemoveDuplicates([]string{}))
	})
	t.Run("RemoveDuplicates", func(t *testing.T) {
		assert.ElementsMatch(t, []string{"a", "c", "d"}, RemoveDuplicates([]string{"a", "c", "c", "d"}))
	})
}

func TestGetFailHosts(t *testing.T) {
	nodes := wfv1.Nodes{
		"retry": wfv1.NodeStatus{
			ID:           "A1",
			HostNodeName: "host1",
			Type:         wfv1.NodeTypeRetry,
			Phase:        wfv1.NodeRunning,
			Children:     []string{"n1", "stepgroup"},
		},
		"n1": wfv1.NodeStatus{
			ID:           "n1",
			HostNodeName: "hostn1",
			Type:         wfv1.NodeTypePod,
			Phase:        wfv1.NodeFailed,
			Children:     []string{},
		},
		"stepgroup": wfv1.NodeStatus{
			ID:           "stepgroup",
			HostNodeName: "host2",
			Type:         wfv1.NodeTypeStepGroup,
			Phase:        wfv1.NodeError,
			Children:     []string{"steps"},
		},
		"steps": wfv1.NodeStatus{
			ID:           "steps",
			HostNodeName: "host4",
			Type:         wfv1.NodeTypeSteps,
			Phase:        wfv1.NodeError,
			Children:     []string{"n2", "n3"},
		},
		"n2": wfv1.NodeStatus{
			ID:           "n2",
			HostNodeName: "hostn2",
			Type:         wfv1.NodeTypePod,
			Phase:        wfv1.NodeRunning,
			Children:     []string{},
		},
		"n3": wfv1.NodeStatus{
			ID:           "n3",
			HostNodeName: "hostn3",
			Type:         wfv1.NodeTypePod,
			Phase:        wfv1.NodeError,
			Children:     []string{},
		},
	}
	t.Run("NotExistParent", func(t *testing.T) {
		assert.Equal(t, []string{}, GetFailHosts(nodes, "not-exist-node"))
	})
	t.Run("ParentWithoutChildrenPodTypeError", func(t *testing.T) {
		assert.Equal(t, []string{"hostn3"}, GetFailHosts(nodes, "n3"))
	})
	t.Run("ParentWithoutChildrenPodTypeRunning", func(t *testing.T) {
		assert.Equal(t, []string{}, GetFailHosts(nodes, "n2"))
	})
	t.Run("ParentWithChildrenFromRetryNode", func(t *testing.T) {
		assert.ElementsMatch(t, GetFailHosts(nodes, "retry"), []string{"hostn1", "hostn3"})
	})
	t.Run("ParentWithChildrenFromNonRetryNode", func(t *testing.T) {
		assert.ElementsMatch(t, GetFailHosts(nodes, "steps"), []string{"hostn3"})
	})
}
