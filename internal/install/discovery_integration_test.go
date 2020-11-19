// +build integration

package install

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiscovery(t *testing.T) {
	cmd := exec.Command("java", "-classpath", "internal/install/mockProcesses", "JavaDaemonTest")
	if err := cmd.Start(); err != nil {
		t.Skipf("error starting java process, skipping TestDiscovery")
	}

	pd := psUtilDiscoverer{}
	manifest, err := pd.discover()
	require.NoError(t, err)
	require.NotNil(t, manifest)

	require.GreaterOrEqual(t, len(manifest.processes), 1)

	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("error sending interrupt to java process")
	}

	cmd.Wait()
}
