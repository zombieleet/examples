package main

import (
	"context"
	"fmt"
	"os"
	execlib "os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/efficientgo/e2e"
	e2edb "github.com/efficientgo/e2e/db"
	e2einteractive "github.com/efficientgo/e2e/interactive"
	e2emonitoring "github.com/efficientgo/e2e/monitoring"
	"github.com/efficientgo/examples/pkg/parquet-export/export1"
	"github.com/efficientgo/examples/pkg/parquet-export/ref"
	"github.com/efficientgo/tools/core/pkg/testutil"
	"github.com/pkg/errors"
)

var (
	generateDataPath = func() string { a, _ := filepath.Abs("generated"); return a }()
	maxTime          = `2021-07-20T00:00:00Z`
)

// Testing export1 for now. Change it to other packages for better performance.
var exportFunction = export1.Export5mAggregations

func TestParquetExport(t *testing.T) {

}

// Test args: -test.timeout 9999m for interactive mode experience.
func TestParquetExportIntegration(t *testing.T) {
	t.Parallel()

	testParquetExportIntegration(testutil.NewTB(t))
}

func BenchmarkParquetExportIntegration(b *testing.B) {
	testParquetExportIntegration(testutil.NewTB(b))
}

func testParquetExportIntegration(tb testutil.TB) {
	ctx := context.Background()

	// Create 10k series for 1w of TSDB blocks. Cache them to 'generated' dir so we don't need to re-create on every run (it takes ~2m).
	_, err := os.Stat(generateDataPath)
	if os.IsNotExist(err) {
		err = exec(
			"sh", "-c",
			fmt.Sprintf("mkdir -p %s && "+
				"docker run -i quay.io/thanos/thanosbench:v0.2.0-rc.1 block plan -p continuous-1w-small --labels 'cluster=\"eu-1\"' --labels 'replica=\"0\"' --max-time=%s | "+
				"docker run -v %s/:/shared -i quay.io/thanos/thanosbench:v0.2.0-rc.1 block gen --output.dir /shared", generateDataPath, maxTime, generateDataPath),
		)
		if err != nil {
			_ = os.RemoveAll(generateDataPath)
		}
	}
	testutil.Ok(tb, err)

	// Start isolated environment with given reference.
	e, err := e2e.NewDockerEnvironment("parquet_bench")
	testutil.Ok(tb, err)
	// Make sure resources (e.g docker containers, network, dir) are cleaned after test.
	tb.Cleanup(e.Close)

	var mon *e2emonitoring.Service
	var p e2e.Runnable
	if !tb.IsBenchmark() {
		// Start monitoring if you want to have interactive look on resources.
		mon, err = e2emonitoring.Start(e, e2emonitoring.WithCurrentProcessAsContainer())
		testutil.Ok(tb, err)

		// Schedule parquet tool, so we can check export produced parquet files.
		// See https://github.com/NathanHowell/parquet-tools for details.
		p = e.Runnable("parquet-tools").Init(
			e2e.StartOptions{
				Image:   "nathanhowell/parquet-tools",
				Command: e2e.NewCommandWithoutEntrypoint("tail", "-f", "/dev/null"),
			},
		)
		testutil.Ok(tb, e2e.StartAndWaitReady(p))
	}

	// Schedule StoreAPI gateway, pointing to local directory with generated dataset.
	testutil.Ok(tb, exec("cp", "-r", generateDataPath+"/.", filepath.Join(e.SharedDir(), "tsdb-data")))
	store := e2edb.NewThanosStore(e, "store", []byte(`type: FILESYSTEM
config:
  directory: "/shared/tsdb-data"
`))
	testutil.Ok(tb, e2e.StartAndWaitReady(store))

	parsedMaxTime, err := time.Parse(time.RFC3339, maxTime)
	testutil.Ok(tb, err)

	minTime := ref.TimestampFromTime(parsedMaxTime.Add(-7 * 24 * time.Hour))
	maxTime := ref.TimestampFromTime(parsedMaxTime)

	for _, tcase := range []struct {
		matchers []*export1.LabelMatcher
	}{
		//{matchers: []*export1.LabelMatcher{{Name: "__name__", Value: "continuous_app_metric9.{1}", Type: export1.LabelMatcher_RE}}}, // 1k series.
		{matchers: []*export1.LabelMatcher{{Name: "__name__", Value: "", Type: export1.LabelMatcher_NEQ}}}, // All, 10k series.
	} {
		tb.Run(fmt.Sprintf("%v", tcase.matchers), func(tb testutil.TB) {
			tb.ResetTimer()

			// Perform export.
			for i := 0; i < tb.N(); i++ {
				start := time.Now()

				f, err := os.OpenFile(filepath.Join(e.SharedDir(), "output.parquet"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
				testutil.Ok(tb, err)
				defer func() {
					if f != nil {
						testutil.Ok(tb, f.Close())
					}
				}()

				seriesNum, samplesNum, err := exportFunction(ctx, store.Endpoint("grpc"), tcase.matchers, minTime, maxTime, f)
				testutil.Ok(tb, err)
				testutil.Ok(tb, f.Close())
				f = nil

				fmt.Println("Export done in ", time.Since(start).String(), "exported", seriesNum, "series,", samplesNum, "samples")

				if !tb.IsBenchmark() {
					// Validate if file is usable, by parquet tooling.
					stdout, stderr, err := p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "rowcount", "-d", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)

					stdout, stderr, err = p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "size", "-d", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)

					// Print 5 records.
					stdout, stderr, err = p.Exec(e2e.NewCommand("java", "-XX:-UsePerfData", "-jar", "/parquet-tools.jar", "head", "/shared/output.parquet"))
					fmt.Println(stdout, stderr)
					testutil.Ok(tb, err)
				}
			}
		})
	}

	if !tb.IsBenchmark() {
		// Uncomment for extra interactive resources.
		testutil.Ok(tb, mon.OpenUserInterfaceInBrowser())
		testutil.Ok(tb, e2einteractive.RunUntilEndpointHit())
	}
}

func exec(cmd string, args ...string) error {
	if o, err := execlib.Command(cmd, args...).CombinedOutput(); err != nil {
		return errors.Wrap(err, string(o))
	}
	return nil
}
