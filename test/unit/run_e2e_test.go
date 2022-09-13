package unit_test

import "testing"

func TestRunE2eTests(t *testing.T) {
	t.Parallel()
	sc := newShellScript(
		loadFile("exec-test-e2e-tests.bash"),
		mockGo(),
		mockKubectl(
			response{startsWith{"config current-context"}, simply("gke_deadbeef_1.24")},
			response{
				startsWith{"get pods --no-headers -n"},
				simply("beef-e3c1 1/1 Running 0 2s\nceed-45b3 1/1 Running 0 1s"),
			},
		),
		mockBinary("gcloud"),
	)
	tcs := []testCase{{
		name: `exec_e2e_tests --run-tests`,
		stdout: []check{
			contains("SETTING UP TEST CLUSTER"),
			contains("Cluster is gke_deadbeef_1.24"),
			contains("kubectl wait job --for=condition=Complete --all -n istio-system --timeout=5m"),
			contains("STARTING KNATIVE SERVING"),
			contains("Waiting until all pods in namespace knative-serving are up"),
			contains("E2E TESTS PASSED"),
		},
	}}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, tc.test(sc))
	}
}