package provider_test

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/Nexenta/go-nexentastor/pkg/ns"
)

// defaults
const (
	defaultUsername       = "admin"
	defaultPassword       = "Nexenta@1"
	defaultPoolName       = "testPool"
	defaultDatasetName    = "testDataset"
	defaultFilesystemName = "testFilesystem"
)

// count of concurrent REST calls to create filesystems on NS
const concurrentProcesses = 20

type config struct {
	address      string
	username     string
	password     string
	pool         string
	dataset      string
	filesystem   string
	smbShareName string
	snapshotName string
	cluster      bool
}

var c *config
var l *logrus.Entry

func TestMain(m *testing.M) {
	var (
		address    = flag.String("address", "", "NS API [schema://host:port,...]")
		username   = flag.String("username", defaultUsername, "overwrite NS API username from config")
		password   = flag.String("password", defaultPassword, "overwrite NS API password from config")
		pool       = flag.String("pool", defaultPoolName, "pool on NS")
		dataset    = flag.String("dataset", defaultDatasetName, "dataset on NS")
		filesystem = flag.String("filesystem", defaultFilesystemName, "filesystem on NS")
		cluster    = flag.Bool("cluster", false, "this is a NS cluster")
		log        = flag.Bool("log", false, "show logs")
	)

	flag.Parse()

	l = logrus.New().WithField("ns", *address)
	l.Logger.SetLevel(logrus.PanicLevel)
	if *log {
		l.Logger.SetLevel(logrus.DebugLevel)
	}

	if *address == "" {
		l.Fatal("--address=[schema://host:port,...] flag cannot be empty")
	}

	c = &config{
		address:      *address,
		username:     *username,
		password:     *password,
		pool:         *pool,
		dataset:      fmt.Sprintf("%s/%s", *pool, *dataset),
		filesystem:   fmt.Sprintf("%s/%s/%s", *pool, *dataset, *filesystem),
		cluster:      *cluster,
		smbShareName: "testShareName",
		snapshotName: "snap-test",
	}

	os.Exit(m.Run())
}

func TestProvider_NewProvider(t *testing.T) {
	t.Logf("Using config:\n---\n%+v\n---", c)

	testSnapshotPath := fmt.Sprintf("%s@%s", c.filesystem, c.snapshotName)
	testSnapshotCloneTargetPath := fmt.Sprintf("%s/csiDriverFsCloned", c.dataset)

	nsp, err := ns.NewProvider(ns.ProviderArgs{
		Address:            c.address,
		Username:           c.username,
		Password:           c.password,
		Log:                l,
		InsecureSkipVerify: true,
	})
	if err != nil {
		t.Error(err)
	}

	t.Run("GetLicense()", func(tt *testing.T) {
		license, err := nsp.GetLicense()
		if err != nil {
			t.Error(err)
		} else if !license.Valid {
			t.Errorf("License %+v is not valid, on NS %s", license, c.address)
		} else if license.Expires[0:2] != "20" {
			tt.Errorf("License expires date should starts with '20': %+v, on NS %s", license, c.address)
		}
	})

	t.Run("GetPools()", func(t *testing.T) {
		pools, err := nsp.GetPools()
		if err != nil {
			t.Error(err)
		} else if !poolArrayContains(pools, c.pool) {
			t.Errorf("Pool %s doesn't exist on NS %s", c.pool, c.address)
		}
	})

	t.Run("GetFilesystems()", func(t *testing.T) {
		filesystems, err := nsp.GetFilesystems(c.pool)
		if err != nil {
			t.Error(err)
		} else if filesystemArrayContains(filesystems, c.pool) {
			t.Errorf("Pool %s should not be in the results", c.pool)
		} else if !filesystemArrayContains(filesystems, c.dataset) {
			t.Errorf("Dataset %s doesn't exist", c.dataset)
		}
	})

	t.Run("GetFilesystem() exists", func(t *testing.T) {
		filesystem, err := nsp.GetFilesystem(c.dataset)
		if err != nil {
			t.Error(err)
		} else if filesystem.Path != c.dataset {
			t.Errorf("No %s filesystem in the result", c.dataset)
		}
	})

	t.Run("GetFilesystem() not exists", func(t *testing.T) {
		nonExistingName := "NON_EXISTING"
		filesystem, err := nsp.GetFilesystem(nonExistingName)
		if err != nil && !strings.Contains(err.Error(), "not found") {
			t.Error(err)
		} else if filesystem.Path != "" {
			t.Errorf("Filesystem %s should not exist, but found in the result: %v", nonExistingName, filesystem)
		}
	})

	t.Run("CreateFilesystem()", func(t *testing.T) {
		destroyFilesystemWithDependents(nsp, c.filesystem)

		err = nsp.CreateFilesystem(ns.CreateFilesystemParams{
			Path: c.filesystem,
		})
		if err != nil {
			t.Error(err)
			return
		}

		filesystems, err := nsp.GetFilesystems(c.dataset)
		if err != nil {
			t.Error(err)
		} else if !filesystemArrayContains(filesystems, c.filesystem) {
			t.Errorf("New filesystem %s wasn't created on NS %s", c.filesystem, c.address)
		}
	})

	t.Run("GetFilesystem() created filesystem should not be shared", func(t *testing.T) {
		filesystem, err := nsp.GetFilesystem(c.filesystem)
		if err != nil {
			t.Error(err)
		} else if filesystem.SharedOverNfs {
			t.Errorf("Created filesystem %s should not be shared over NFS (NS %s)", c.filesystem, c.address)
		} else if filesystem.SharedOverSmb {
			t.Errorf("Created filesystem %s should not be shared over SMB (NS %s)", c.filesystem, c.address)
		}
	})

	t.Run("CreateNfsShare()", func(t *testing.T) {
		nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})

		err = nsp.CreateNfsShare(ns.CreateNfsShareParams{
			Filesystem: c.filesystem,
		})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("GetFilesystem() created filesystem should be shared over NFS", func(t *testing.T) {
		filesystem, err := nsp.GetFilesystem(c.filesystem)
		if err != nil {
			t.Error(err)
		} else if !filesystem.SharedOverNfs {
			t.Errorf("Created filesystem %s should be shared (NS %s)", c.filesystem, c.address)
		}
	})

	t.Run("nfs share should appear on NS", func(t *testing.T) {
		//TODO other way to cut out host from address
		host := strings.Split(c.address, "//")[1]
		host = strings.Split(host, ":")[0]

		out, err := exec.Command("showmount", "-e", host).Output()
		if err != nil {
			t.Error(err)
		} else if !strings.Contains(fmt.Sprintf("%s", out), c.filesystem) {
			t.Errorf("Cannot find '%s' nfs in the 'showmount' output: \n---\n%s\n---\n", c.filesystem, out)
		}
	})

	t.Run("DeleteNfsShare()", func(t *testing.T) {
		filesystems, err := nsp.GetFilesystems(c.dataset)
		if err != nil {
			t.Error(err)
			return
		} else if !filesystemArrayContains(filesystems, c.filesystem) {
			t.Skipf("Filesystem %s doesn't exist on NS %s", c.filesystem, c.address)
			return
		}

		err = nsp.DeleteNfsShare(c.filesystem)
		if err != nil {
			t.Error(err)
		}
	})

	for _, smbShareName := range []string{c.smbShareName, ""} {
		smbShareName := smbShareName

		t.Run(
			fmt.Sprintf("CreateSmbShare() should create SMB share with '%s' share name", smbShareName),
			func(t *testing.T) {
				nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})

				err = nsp.CreateSmbShare(ns.CreateSmbShareParams{
					Filesystem: c.filesystem,
					ShareName:  smbShareName,
				})
				if err != nil {
					t.Error(err)
				}
			},
		)

		t.Run("GetFilesystem() created filesystem should be shared over SMB", func(t *testing.T) {
			filesystem, err := nsp.GetFilesystem(c.filesystem)
			if err != nil {
				t.Error(err)
			} else if !filesystem.SharedOverSmb {
				t.Errorf("Created filesystem %s should be shared over SMB (NS %s)", c.filesystem, c.address)
			}
		})

		t.Run("GetSmbShareName() should return SMB share name", func(t *testing.T) {
			filesystem, err := nsp.GetFilesystem(c.filesystem)
			if err != nil {
				t.Error(err)
				return
			}

			var expectedShareName string
			if smbShareName == "" {
				expectedShareName = filesystem.GetDefaultSmbShareName()
			} else {
				expectedShareName = smbShareName
			}

			shareName, err := nsp.GetSmbShareName(c.filesystem)
			if err != nil {
				t.Error(err)
			} else if shareName != expectedShareName {
				t.Errorf(
					"Expected shareName='%s' but got '%s', for filesystem '%s' on NS %s",
					expectedShareName,
					shareName,
					c.filesystem,
					c.address,
				)
			}
		})

		//TODO test SMB share, mount cifs?

		t.Run("DeleteSmbShare()", func(t *testing.T) {
			err = nsp.DeleteSmbShare(c.filesystem)
			if err != nil {
				t.Error(err)
			}
		})
	}

	t.Run("DestroyFilesystem()", func(t *testing.T) {
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})

		err = nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{DestroySnapshots: true})
		if err != nil {
			t.Error(err)
			return
		}

		filesystems, err := nsp.GetFilesystems(c.dataset)
		if err != nil {
			t.Error(err)
		} else if filesystemArrayContains(filesystems, c.filesystem) {
			t.Errorf("Filesystem %s still exists on NS %s", c.filesystem, c.address)
		}
	})

	t.Run("CreateFilesystem() with referenced quota size", func(t *testing.T) {
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})

		var referencedQuotaSize int64 = 2 * 1024 * 1024 * 1024

		err = nsp.CreateFilesystem(ns.CreateFilesystemParams{
			Path:                c.filesystem,
			ReferencedQuotaSize: referencedQuotaSize,
		})
		if err != nil {
			t.Error(err)
			return
		}

		filesystem, err := nsp.GetFilesystem(c.filesystem)
		if err != nil {
			t.Error(err)
			return
		} else if filesystem.GetReferencedQuotaSize() != referencedQuotaSize {
			t.Errorf(
				"New filesystem %s referenced quota size expected to be %d, but got %d (NS %s)",
				filesystem.Path,
				referencedQuotaSize,
				filesystem.GetReferencedQuotaSize(),
				c.address,
			)
		}
	})

	t.Run("CreateSnapshot()", func(t *testing.T) {
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})

		err = nsp.CreateSnapshot(ns.CreateSnapshotParams{
			Path: testSnapshotPath,
		})
		if err != nil {
			t.Error(err)
		}

		snapshot, err := nsp.GetSnapshot(testSnapshotPath)
		if err != nil {
			t.Error(err)
			return
		} else if snapshot.Path != testSnapshotPath {
			t.Errorf(
				"Created snapshot path expected to be '%s', but got '%s' (Snapshot: %+v, NS %s)",
				testSnapshotPath,
				snapshot.Path,
				snapshot,
				c.address,
			)
			return
		} else if snapshot.Name != c.snapshotName {
			t.Errorf(
				"Created snapshot name expected to be '%s', but got '%s' (Snapshot: %+v, NS %s)",
				c.snapshotName,
				snapshot.Name,
				snapshot,
				c.address,
			)
			return
		} else if snapshot.Parent != c.filesystem {
			t.Errorf(
				"Created snapshot parent expected to be '%s', but got '%s' (Snapshot: %+v, NS %s)",
				c.filesystem,
				snapshot.Parent,
				snapshot,
				c.address,
			)
			return
		}

		snapshots, err := nsp.GetSnapshots(c.filesystem, true)
		if err != nil {
			t.Errorf("Cannot get '%s' snapshot list: %v", c.filesystem, err)
			return
		} else if len(snapshots) == 0 {
			t.Errorf(
				"New snapshot '%s' was not found in '%s' snapshot list, the list is empty: %v",
				c.snapshotName,
				c.filesystem,
				snapshots,
			)
			return
		} else if !snapshotArrayContains(snapshots, testSnapshotPath) {
			t.Errorf(
				"New snapshot '%s' was not found in '%s' snapshot list: %v",
				c.snapshotName,
				c.filesystem,
				snapshots,
			)
			return
		}
	})

	t.Run("CloneSnapshot()", func(t *testing.T) {
		nsp.DestroySnapshot(testSnapshotPath)
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.DestroyFilesystem(testSnapshotCloneTargetPath, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})

		err := nsp.CreateSnapshot(ns.CreateSnapshotParams{Path: testSnapshotPath})
		if err != nil {
			t.Error(err)
			return
		}

		err = nsp.CloneSnapshot(testSnapshotPath, ns.CloneSnapshotParams{
			TargetPath: testSnapshotCloneTargetPath,
		})
		if err != nil {
			t.Error(err)
			return
		}

		_, err = nsp.GetFilesystem(testSnapshotCloneTargetPath)
		if err != nil {
			t.Errorf("Cannot get created filesystem '%s': %v", testSnapshotCloneTargetPath, err)
			return
		}
	})

	t.Run("PromoteFilesystem()", func(t *testing.T) {
		err := nsp.PromoteFilesystem(testSnapshotCloneTargetPath)
		if err != nil {
			t.Error(err)
		}

		//TODO check that snapshots have been actually "migrated" to the promoted filesystem
	})

	t.Run("DestroySnapshot()", func(t *testing.T) {
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})
		nsp.CreateSnapshot(ns.CreateSnapshotParams{Path: testSnapshotPath})

		err := nsp.DestroySnapshot(testSnapshotPath)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("DestroyFilesystem() for filesystem with snapshots", func(t *testing.T) {
		nsp.DestroySnapshot(testSnapshotPath)
		nsp.DestroyFilesystem(testSnapshotCloneTargetPath, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})

		err := nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})
		if err != nil {
			t.Errorf("Failed to create preconditions: Create filesystem '%s' failed: %v", c.filesystem, err)
			return
		}
		err = nsp.CreateSnapshot(ns.CreateSnapshotParams{Path: testSnapshotPath})
		if err != nil {
			t.Errorf("Failed to create preconditions: Create snapshot '%s' failed: %v", testSnapshotPath, err)
			return
		}

		err = nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{DestroySnapshots: false})
		if !ns.IsBusyNefError(err) {
			t.Errorf(
				`Filesystem delete request is supposed to return EBUSY error in case of deleting
				filesystem with snapshots, but it's not: %v`,
				err,
			)
			return
		}

		err = nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{DestroySnapshots: true})
		if err != nil {
			t.Errorf("Cannot destroy filesystem, even with snapshots=true option: %v", err)
			return
		}

		filesystem, err := nsp.GetFilesystem(c.filesystem)
		if !ns.IsNotExistNefError(err) {
			t.Errorf(
				"Get filesystem request should return ENOENT error, but it returns filesystem: %v, error: %v",
				filesystem,
				err,
			)
		}
	})

	t.Run("DestroyFilesystem() for filesystem with clones", func(t *testing.T) {
		nsp.DestroySnapshot(testSnapshotPath)
		nsp.DestroyFilesystem(testSnapshotCloneTargetPath, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})

		err := nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem})
		if err != nil {
			t.Errorf("Failed to create preconditions: Create filesystem '%s' failed: %v", c.filesystem, err)
			return
		}
		err = nsp.CreateSnapshot(ns.CreateSnapshotParams{Path: testSnapshotPath})
		if err != nil {
			t.Errorf("Failed to create preconditions: Create snapshot '%s' failed: %v", testSnapshotPath, err)
			return
		}
		err = nsp.CloneSnapshot(testSnapshotPath, ns.CloneSnapshotParams{
			TargetPath: testSnapshotCloneTargetPath,
		})
		if err != nil {
			t.Errorf(
				"Failed to create preconditions: Create clone '%s' of '%s' failed: %v",
				testSnapshotCloneTargetPath,
				testSnapshotPath,
				err,
			)
			return
		}

		err = nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: false,
		})
		if !ns.IsAlreadyExistNefError(err) {
			t.Errorf(
				`Filesystem delete request is supposed to return EEXIST error in case of deleting
				filesystem with clones, but it's not: %v`,
				err,
			)
			return
		}

		err = nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})
		if err != nil {
			t.Errorf("Cannot destroy filesystem: %v", err)
			return
		}

		filesystem, err := nsp.GetFilesystem(c.filesystem)
		if !ns.IsNotExistNefError(err) {
			t.Errorf(
				"Get filesystem request should return ENOENT error, but it returns filesystem: %v, error: %v",
				filesystem,
				err,
			)
		}

		filesystem, err = nsp.GetFilesystem(testSnapshotCloneTargetPath)
		if err != nil {
			t.Errorf(
				"Cloned filesystem '%s' should be presented, but there is an error while getting it: %v",
				testSnapshotCloneTargetPath,
				err,
			)
		}
	})

	t.Run("GetFilesystemAvailableCapacity()", func(t *testing.T) {
		nsp.DestroyFilesystem(c.filesystem, ns.DestroyFilesystemParams{
			DestroySnapshots:               true,
			PromoteMostRecentCloneIfExists: true,
		})

		var referencedQuotaSize int64 = 3 * 1024 * 1024 * 1024

		err = nsp.CreateFilesystem(ns.CreateFilesystemParams{
			Path:                c.filesystem,
			ReferencedQuotaSize: referencedQuotaSize,
		})
		if err != nil {
			t.Error(err)
			return
		}

		availableCapacity, err := nsp.GetFilesystemAvailableCapacity(c.filesystem)
		if err != nil {
			t.Error(err)
			return
		} else if availableCapacity == 0 {
			t.Errorf("New filesystem %s indicates wrong available capacity (0), on: %s", c.filesystem, c.address)
		} else if availableCapacity >= referencedQuotaSize {
			t.Errorf(
				"New filesystem %s available capacity expected to be more or equal to %d, but got %d (NS %s)",
				c.filesystem,
				referencedQuotaSize,
				availableCapacity,
				c.address,
			)
		}
	})

	t.Run("GetRSFClusters()", func(t *testing.T) {
		expectedToBeACluster := c.cluster

		clusters, err := nsp.GetRSFClusters()
		if err != nil {
			t.Error(err)
			return
		}

		if expectedToBeACluster && len(clusters) == 0 {
			t.Errorf(
				"NS %s expected to be in a cluster (--cluster=true flag) but got no clusters from the API",
				c.address,
			)
		} else if !expectedToBeACluster && len(clusters) > 0 {
			t.Errorf(
				"NS %s expected not to be in a cluster (--cluster=false flag) but got clusters from the API: %+v",
				c.address,
				clusters,
			)
		}
	})

	t.Run("GetFilesystemsSlice()", func(t *testing.T) {
		destroyFilesystemWithDependents(nsp, c.filesystem)

		err = nsp.CreateFilesystem(ns.CreateFilesystemParams{
			Path: c.filesystem,
		})
		if err != nil {
			t.Error(err)
			return
		}

		count := 5
		err = createFilesystemChildren(nsp, c.filesystem, count)
		if err != nil {
			t.Error(err)
			return
		}

		filesystems, err := nsp.GetFilesystemsSlice(c.filesystem, 0, 0)
		if err == nil {
			t.Errorf("Should return an error when limit is equal 0, but got: %v", err)
			return
		}

		filesystems, err = nsp.GetFilesystemsSlice(c.filesystem, 2, 0)
		if err != nil {
			t.Error(err)
			return
		} else if len(filesystems) != 2 {
			t.Errorf("GetFilesystems() returned %d filesystems, but expected 2", len(filesystems))
			return
		} else if filesystems[0].Path != getFilesystemChildName(c.filesystem, 1) {
			t.Errorf(
				"Limit: '2', offset: '0' - first item expected to be '%s' but got: %+v",
				getFilesystemChildName(c.filesystem, 1),
				filesystems,
			)
			return
		} else if filesystems[1].Path != getFilesystemChildName(c.filesystem, 2) {
			t.Errorf(
				"Limit: '2', offset: '0' - second item expected to be '%s' but got: %+v",
				getFilesystemChildName(c.filesystem, 2),
				filesystems,
			)
			return
		}

		filesystems, err = nsp.GetFilesystemsSlice(c.filesystem, 4, 3)
		if err != nil {
			t.Error(err)
			return
		} else if len(filesystems) != 3 {
			t.Errorf("Returned %d filesystems, but expected 3", len(filesystems))
			return
		} else if filesystems[0].Path != getFilesystemChildName(c.filesystem, 3) {
			t.Errorf(
				"Limit: '4', offset: '3' - first item expected to be '%s' but got: %+v",
				getFilesystemChildName(c.filesystem, 3),
				filesystems,
			)
			return
		} else if filesystems[1].Path != getFilesystemChildName(c.filesystem, 4) {
			t.Errorf(
				"Limit: '4', offset: '3' - second item expected to be '%s' but got: %+v",
				getFilesystemChildName(c.filesystem, 4),
				filesystems,
			)
			return
		} else if filesystems[2].Path != getFilesystemChildName(c.filesystem, 5) {
			t.Errorf(
				"Limit: '4', offset: '3' - third item expected to be '%s' but got: %+v",
				getFilesystemChildName(c.filesystem, 5),
				filesystems,
			)
			return
		}
	})

	t.Run("GetFilesystems() pagination", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping pagination test in short mode")
			return
		}

		destroyFilesystemWithDependents(nsp, c.filesystem)

		t.Log("create parent filesystem")
		err = nsp.CreateFilesystem(ns.CreateFilesystemParams{
			Path: c.filesystem,
		})
		if err != nil {
			t.Error(err)
			return
		}

		count := 101
		t.Logf("create %d children filesystems", count)
		err = createFilesystemChildren(nsp, c.filesystem, count)
		if err != nil {
			t.Error(err)
			return
		}

		t.Log("get all filesystems")
		filesystems, err := nsp.GetFilesystems(c.filesystem)
		if err != nil {
			t.Error(err)
			return
		} else if len(filesystems) != count {
			t.Errorf("GetFilesystems() returned %d filesystems, but expected %d", len(filesystems), count)
		}

		t.Log("check if all filesystems are in the list")
		for i := 1; i <= len(filesystems); i++ {
			fs := getFilesystemChildName(c.filesystem, i)
			if !filesystemArrayContains(filesystems, fs) {
				t.Errorf("Filesystem list doesn't contain '%s' filesystem", fs)
			}
		}
	})

	t.Run("GetFilesystemsWithStartingToken() pagination", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping pagination test in short mode")
			return
		}

		destroyFilesystemWithDependents(nsp, c.filesystem)

		t.Log("create parent filesystem")
		if err = nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: c.filesystem}); err != nil {
			t.Error(err)
			return
		}

		count := 101
		t.Logf("create %d children filesystems", count)
		if err = createFilesystemChildren(nsp, c.filesystem, count); err != nil {
			t.Error(err)
			return
		}

		tokenTests := []struct {
			Limit             int
			ExpectedCount     int
			StartingToken     string
			ExpectedNextToken string
		}{
			{count + 1, count, "", ""},
			{5, 5, "", getFilesystemChildName(c.filesystem, 5)},
			{5, 5, getFilesystemChildName(c.filesystem, 5), getFilesystemChildName(c.filesystem, 10)},
			{10, 10, getFilesystemChildName(c.filesystem, 3), getFilesystemChildName(c.filesystem, 13)},
			{10, 1, getFilesystemChildName(c.filesystem, 100), ""},
			{1, 0, "NOT_EXISTING", ""},
			{0, count, "", ""},
		}

		for _, v := range tokenTests {
			f := fmt.Sprintf("startingToken: '%s', limit: '%d'", v.StartingToken, v.Limit)
			t.Logf("...check %s", f)

			filesystems, nextToken, err := nsp.GetFilesystemsWithStartingToken(c.filesystem, v.StartingToken, v.Limit)
			if err != nil {
				t.Error(err)
				return
			} else if len(filesystems) != v.ExpectedCount {
				t.Errorf("%s: returned %d filesystems, but expected %d", f, len(filesystems), v.ExpectedCount)
				return
			} else if nextToken != v.ExpectedNextToken {
				t.Errorf("%s: returned '%s' next token, but expected '%s'", f, nextToken, v.ExpectedNextToken)
			}
		}

		t.Log("get all filesystems using tokens")
		nextToken := ""
		filesystems := []ns.Filesystem{}
		for {
			filesystemsSlice, nt, err := nsp.GetFilesystemsWithStartingToken(c.filesystem, nextToken, 25)
			if err != nil {
				t.Error(err)
				return
			}

			filesystems = append(filesystems, filesystemsSlice...)
			if len(filesystems) > count {
				t.Errorf(
					"Get all filesystems operation is expected to return %d filesystems, but already got %d; Cancel.",
					count,
					len(filesystems),
				)
				return
			}

			if nt == "" {
				break
			} else {
				nextToken = nt
			}
		}

		if len(filesystems) != count {
			t.Errorf("Get all filesystems operation returned %d filesystems but expected %d", len(filesystems), count)
			return
		}

		t.Log("check if all filesystems are in the list")
		for i := 1; i <= len(filesystems); i++ {
			fs := getFilesystemChildName(c.filesystem, i)
			if !filesystemArrayContains(filesystems, getFilesystemChildName(c.filesystem, i)) {
				t.Errorf("Filesystem list doesn't contain '%s' filesystem", fs)
			}
		}
	})

	// clean up
	nsp.DestroySnapshot(testSnapshotPath)
	destroyFilesystemWithDependents(nsp, testSnapshotCloneTargetPath)
	destroyFilesystemWithDependents(nsp, c.filesystem)
}

// getFilesystemChildName("fs", 13) === "fs/child-013"
func getFilesystemChildName(parent string, id int) string {
	return path.Join(parent, fmt.Sprintf("child-%03d", id))
}

func createFilesystemChildren(nsp ns.ProviderInterface, parent string, count int) error {
	jobs := make([]func() error, count)
	for i := 0; i < count; i++ {
		i := i
		jobs[i] = func() error {
			return nsp.CreateFilesystem(ns.CreateFilesystemParams{Path: getFilesystemChildName(parent, i+1)})
		}
	}

	return runConcurrentJobs("create filesystem", jobs)
}

func destroyFilesystemWithDependents(nsp ns.ProviderInterface, filesystem string) error {
	children, err := nsp.GetFilesystems(filesystem)
	if err != nil {
		return fmt.Errorf("destroyFilesystemWithDependents(%s): failed to get children: %v", filesystem, err)
	}

	if len(children) > 0 {
		jobs := make([]func() error, len(children))
		for i, c := range children {
			c := c
			jobs[i] = func() error {
				return destroyFilesystemWithDependents(nsp, c.Path)
			}
		}
		err := runConcurrentJobs("delete filesystem", jobs)
		if err != nil {
			return fmt.Errorf(
				"destroyFilesystemWithDependents(%s): failed to remove child filesystems: %s",
				filesystem,
				err,
			)
		}
	}

	err = nsp.DestroyFilesystem(filesystem, ns.DestroyFilesystemParams{DestroySnapshots: true})
	if err != nil {
		return fmt.Errorf("destroyFilesystemWithDependents(%s): failed to destroy filesystem: %v", filesystem, err)
	}

	return nil
}

func runConcurrentJobs(description string, jobs []func() error) error {
	count := len(jobs)

	worker := func(jobsPool <-chan func() error, results chan<- error) {
		for job := range jobsPool {
			err := job()
			if err != nil {
				results <- fmt.Errorf("Job failed: %s: %s", description, err)
			} else {
				results <- nil
			}
		}
	}

	jobsPool := make(chan func() error, count)
	results := make(chan error, count)

	// start workers
	for i := 0; i < concurrentProcesses; i++ {
		go worker(jobsPool, results)
	}

	// schedule jobs
	for _, job := range jobs {
		jobsPool <- job
	}
	close(jobsPool)

	// collect all results
	errors := []error{}
	for i := 0; i < count; i++ {
		err := <-results
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		err := ""
		for _, e := range errors {
			err += fmt.Sprintf("\n%s;", e)
		}
		return fmt.Errorf("%d of %d jobs failed: %s: %s", len(errors), count, description, err)
	}

	return nil
}

func poolArrayContains(array []ns.Pool, value string) bool {
	for _, v := range array {
		if v.Name == value {
			return true
		}
	}
	return false
}

func filesystemArrayContains(array []ns.Filesystem, value string) bool {
	for _, v := range array {
		if v.Path == value {
			return true
		}
	}
	return false
}

func snapshotArrayContains(array []ns.Snapshot, value string) bool {
	for _, v := range array {
		if v.Path == value {
			return true
		}
	}
	return false
}
