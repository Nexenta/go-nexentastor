# ns
--
    import "github.com/Nexenta/go-nexentastor/pkg/ns"


## Usage

#### func  GetNefErrorCode

```go
func GetNefErrorCode(err error) string
```
GetNefErrorCode - treats an error as NefError and returns its code in case of
success

#### func  IsAlreadyExistNefError

```go
func IsAlreadyExistNefError(err error) bool
```
IsAlreadyExistNefError treats an error as NefError and returns true if its code
is "EEXIST"

#### func  IsAuthNefError

```go
func IsAuthNefError(err error) bool
```
IsAuthNefError treats an error as NefError and returns true if its code is
"EAUTH"

#### func  IsBadArgNefError

```go
func IsBadArgNefError(err error) bool
```
IsBadArgNefError treats an error as NefError and returns true if its code is
"EBADARG"

#### func  IsBusyNefError

```go
func IsBusyNefError(err error) bool
```
IsBusyNefError treats an error as NefError and returns true if its code is
"EBUSY" Example: filesystem cannot be deleted because it has snapshots

#### func  IsNefError

```go
func IsNefError(err error) bool
```
IsNefError - checks if an error is an NefError

#### func  IsNotExistNefError

```go
func IsNotExistNefError(err error) bool
```
IsNotExistNefError treats an error as NefError and returns true if its code is
"ENOENT"

#### type ACLRuleSet

```go
type ACLRuleSet int64
```

ACLRuleSet - filesystem ACL rule set

```go
const (
	// ACLReadOnly - apply read only set of rules to filesystem
	ACLReadOnly ACLRuleSet = iota

	// ACLReadWrite - apply full access set of rules to filesystem
	ACLReadWrite
)
```

#### type CloneSnapshotParams

```go
type CloneSnapshotParams struct {
	// filesystem path w/o leading slash
	TargetPath          string `json:"targetPath"`
	ReferencedQuotaSize int64  `json:"referencedQuotaSize,omitempty"`
}
```

CloneSnapshotParams - params to clone snapshot to filesystem

#### type CreateFilesystemParams

```go
type CreateFilesystemParams struct {
	// filesystem path w/o leading slash
	Path string `json:"path"`
	// filesystem referenced quota size in bytes
	ReferencedQuotaSize int64 `json:"referencedQuotaSize,omitempty"`
}
```

CreateFilesystemParams - params to create filesystem

#### type CreateISCSITargetParams

```go
type CreateISCSITargetParams struct {
	Name    string   `json:"name"`
	Portals []Portal `json:"portals"`
}
```

CreateISCSITargetParamas - params to create new iSCSI target

#### type CreateLunMappingParams

```go
type CreateLunMappingParams struct {
	HostGroup   string `json:"hostGroup"`
	Volume      string `json:"volume"`
	TargetGroup string `json:"targetGroup"`
}
```

CreateLunMappingParams - params to create new lun

#### type CreateNfsShareParams

```go
type CreateNfsShareParams struct {
	// filesystem path w/o leading slash
	Filesystem    string        `json:"filesystem"`
	ReadWriteList []NfsRuleList `json:"readWriteList"`
	ReadOnlyList  []NfsRuleList `json:"readOnlyList"`
}
```

CreateNfsShareParams - params to create NFS share

#### type CreateSmbShareParams

```go
type CreateSmbShareParams struct {
	// filesystem path w/o leading slash
	Filesystem string `json:"filesystem"`
	// share name, used in mount command
	ShareName string `json:"shareName,omitempty"`
}
```

CreateSmbShareParams - params to create SMB share

#### type CreateSnapshotParams

```go
type CreateSnapshotParams struct {
	// snapshot path w/o leading slash
	Path string `json:"path"`
}
```

CreateSnapshotParams - params to create snapshot

#### type CreateTargetGroupParams

```go
type CreateTargetGroupParams struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
```

CreateTargetGroupParams - params to create target group

#### type CreateVolumeParams

```go
type CreateVolumeParams struct {
	// volume path w/o leading slash
	Path       string `json:"path"`
	VolumeSize int64  `json:"volumeSize"`
}
```

CreateVolumeParams - params to create a volume

#### type DestroyFilesystemParams

```go
type DestroyFilesystemParams struct {
	// If set to `true`, then tries to destroy filesystem's snapshots as well.
	// In case some snapshots have clones, the filesystem cannot be deleted
	// without deleting all dependent clones, OR promoting one of the clones
	// to take over the snapshots (see "PromoteMostRecentCloneIfExists" parameter).
	DestroySnapshots bool

	// If set to `true`, then tries to find the most recent snapshot clone and if found one,
	// that clone will be promoted to take over all the snapshots from the original filesystem,
	// then the original filesystem will be destroyed.
	//
	// Initial state:
	//    [fsSource]---+                       // source filesystem
	//                 |    [snapshot1]        // source filesystem snapshots
	//                 |    [snapshot2]
	//                 `--->[snapshot3]<---+
	//                                     |
	//    [fsClone1]-----------------------+   // filesystem clone of "snapshot3"
	//    [fsClone2]-----------------------+   // another filesystem clone of "snapshot3"
	//
	// After destroy "fsSource" filesystem call (PromoteMostRecentCloneIfExists=true and DestroySnapshots=true):
	//    [fsClone1]<----------------------+   // "fsClone1" is still linked to "snapshot3"
	//    [fsClone2]---+                   |   // "fsClone2" is got promoted to take over snapshots of "fsSource"
	//                 |    [snapshot1]    |
	//                 |    [snapshot2]    |
	//                 `--->[snapshot3]<---+
	//
	PromoteMostRecentCloneIfExists bool
}
```

DestroyFilesystemParams - filesystem deletion parameters

#### type DestroyVolumeParams

```go
type DestroyVolumeParams struct {
	DestroySnapshots               bool
	PromoteMostRecentCloneIfExists bool
}
```


#### type Filesystem

```go
type Filesystem struct {
	Path           string `json:"path"`
	MountPoint     string `json:"mountPoint"`
	SharedOverNfs  bool   `json:"sharedOverNfs"`
	SharedOverSmb  bool   `json:"sharedOverSmb"`
	BytesAvailable int64  `json:"bytesAvailable"`
	BytesUsed      int64  `json:"bytesUsed"`
}
```

Filesystem - NexentaStor filesystem

#### func (*Filesystem) GetDefaultSmbShareName

```go
func (fs *Filesystem) GetDefaultSmbShareName() string
```
GetDefaultSmbShareName - get default SMB share name (all slashes get replaced by
underscore) Converts '/pool/dataset/fs' to 'pool_dataset_fs'

#### func (*Filesystem) GetReferencedQuotaSize

```go
func (fs *Filesystem) GetReferencedQuotaSize() int64
```
GetReferencedQuotaSize - get total referenced quota size

#### func (*Filesystem) String

```go
func (fs *Filesystem) String() string
```

#### type License

```go
type License struct {
	Valid   bool   `json:"valid"`
	Expires string `json:"expires"`
}
```

License - NexentaStor license

#### type LunMapping

```go
type LunMapping struct {
	Id          string `json:"id"`
	Volume      string `json:"volume"`
	TargetGroup string `json:"targetGroup"`
	HostGroup   string `json:"hostGroup"`
	Lun         int    `json:"lun"`
}
```

LunMapping - NexentaStor lunmapping

#### type NefError

```go
type NefError struct {
	Err  error
	Code string
}
```

NefError - nef error format

#### func (*NefError) Error

```go
func (e *NefError) Error() string
```

#### type NfsRuleList

```go
type NfsRuleList struct {
	Etype  string `json:"etype"`
	Entity string `json:"entity"`
	Mask   int    `json:"mask"`
}
```


#### type Pool

```go
type Pool struct {
	Name string `json:"poolName"`
}
```

Pool - NS pool

#### type Portal

```go
type Portal struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}
```


#### type Provider

```go
type Provider struct {
	Address    string
	Username   string
	Password   string
	RestClient rest.ClientInterface
	Log        *logrus.Entry
}
```

Provider - NexentaStor API provider

#### func (*Provider) CloneSnapshot

```go
func (p *Provider) CloneSnapshot(path string, params CloneSnapshotParams) error
```
CloneSnapshot clones snapshot to FS

#### func (*Provider) CreateFilesystem

```go
func (p *Provider) CreateFilesystem(params CreateFilesystemParams) error
```
CreateFilesystem creates filesystem by path

#### func (*Provider) CreateISCSITarget

```go
func (p *Provider) CreateISCSITarget(params CreateISCSITargetParams) error
```
CreateISCSITarget - create new iSCSI target on NexentaStor

#### func (*Provider) CreateLunMapping

```go
func (p *Provider) CreateLunMapping(params CreateLunMappingParams) error
```
CreateLunMapping - creates lun for given volume

#### func (*Provider) CreateNfsShare

```go
func (p *Provider) CreateNfsShare(params CreateNfsShareParams) error
```
CreateNfsShare creates NFS share on specified filesystem CLI test:

    showmount -e HOST
    mkdir -p /mnt/test && sudo mount -v -t nfs HOST:/pool/fs /mnt/test
    findmnt /mnt/test

#### func (*Provider) CreateSmbShare

```go
func (p *Provider) CreateSmbShare(params CreateSmbShareParams) error
```
CreateSmbShare creates SMB share (cifs) on specified filesystem Leave shareName
empty to generate default value CLI test:

    mkdir -p /mnt/test && sudo mount -v -t cifs -o username=admin,password=Nexenta@1 //HOST//pool_fs /mnt/test
    findmnt /mnt/test

#### func (*Provider) CreateSnapshot

```go
func (p *Provider) CreateSnapshot(params CreateSnapshotParams) error
```
CreateSnapshot creates snapshot by filesystem path

#### func (*Provider) CreateUpdateTargetGroup

```go
func (p *Provider) CreateUpdateTargetGroup(params CreateTargetGroupParams) error
```
CreateUpdateTargetGroup - create new target group on NexentaStor

#### func (*Provider) CreateVolume

```go
func (p *Provider) CreateVolume(params CreateVolumeParams) error
```
CreateVolume creates volume by path and size

#### func (*Provider) DeleteNfsShare

```go
func (p *Provider) DeleteNfsShare(path string) error
```
DeleteNfsShare destroys NFS chare by filesystem path

#### func (*Provider) DeleteSmbShare

```go
func (p *Provider) DeleteSmbShare(path string) error
```
DeleteSmbShare destroys SMB share by filesystem path

#### func (*Provider) DestroyFilesystem

```go
func (p *Provider) DestroyFilesystem(path string, params DestroyFilesystemParams) error
```
DestroyFilesystem destroys filesystem on NS, may destroy snapshots and promote
clones (see DestroyFilesystemParams) Path format: 'pool/dataset/filesystem'

#### func (*Provider) DestroyLunMapping

```go
func (p *Provider) DestroyLunMapping(id string) error
```

#### func (*Provider) DestroySnapshot

```go
func (p *Provider) DestroySnapshot(path string) error
```
DestroySnapshot destroys snapshot by path

#### func (*Provider) DestroyVolume

```go
func (p *Provider) DestroyVolume(path string, params DestroyVolumeParams) error
```

#### func (*Provider) GetFilesystem

```go
func (p *Provider) GetFilesystem(path string) (filesystem Filesystem, err error)
```
GetFilesystem returns NexentaStor filesystem by its path

#### func (*Provider) GetFilesystemAvailableCapacity

```go
func (p *Provider) GetFilesystemAvailableCapacity(path string) (int64, error)
```
GetFilesystemAvailableCapacity returns NexentaStor filesystem available size by
its path

#### func (*Provider) GetFilesystems

```go
func (p *Provider) GetFilesystems(parent string) ([]Filesystem, error)
```
GetFilesystems returns all NexentaStor filesystems by parent filesystem

#### func (*Provider) GetFilesystemsSlice

```go
func (p *Provider) GetFilesystemsSlice(parent string, limit, offset int) ([]Filesystem, error)
```
GetFilesystemsSlice returns a slice of filesystems by parent filesystem with
specified limit and offset offset - the first record number of collection, that
would be included in result

#### func (*Provider) GetFilesystemsWithStartingToken

```go
func (p *Provider) GetFilesystemsWithStartingToken(parent string, startingToken string, limit int) (
	filesystems []Filesystem,
	nextToken string,
	err error,
)
```
GetFilesystemsWithStartingToken returns filesystems by parent filesystem after
specified starting token parent - parent filesystem's path startingToken - a
path to a specific filesystem to start AFTER this token limit - the maximum
count of filesystems to return in the list Function may return nextToken if
there is more filesystems than limit value

#### func (*Provider) GetLicense

```go
func (p *Provider) GetLicense() (license License, err error)
```
GetLicense returns NexentaStor license

#### func (*Provider) GetLunMapping

```go
func (p *Provider) GetLunMapping(path string) (lunMapping LunMapping, err error)
```
GetLunMapping returns NexentaStor lunmapping for a volume

#### func (*Provider) GetPools

```go
func (p *Provider) GetPools() ([]Pool, error)
```
GetPools returns NexentaStor pools

#### func (*Provider) GetRSFClusters

```go
func (p *Provider) GetRSFClusters() ([]RSFCluster, error)
```
GetRSFClusters returns RSF clusters from NS

#### func (*Provider) GetSmbShareName

```go
func (p *Provider) GetSmbShareName(path string) (string, error)
```
GetSmbShareName returns share name for filesystem that shared over SMB

#### func (*Provider) GetSnapshot

```go
func (p *Provider) GetSnapshot(path string) (snapshot Snapshot, err error)
```
GetSnapshot returns snapshot by its path path - full path to snapshot w/o
leading slash (e.g. "p/d/fs@s")

#### func (*Provider) GetSnapshots

```go
func (p *Provider) GetSnapshots(volumePath string, recursive bool) ([]Snapshot, error)
```
GetSnapshots returns snapshots by volume path

#### func (*Provider) GetVolume

```go
func (p *Provider) GetVolume(path string) (volume Volume, err error)
```
GetVolume - returns NexentaStor volume properties

#### func (*Provider) GetVolumeGroup

```go
func (p *Provider) GetVolumeGroup(path string) (volumeGroup VolumeGroup, err error)
```
GetVolumeGroup returns NexentaStor volumeGroup by its path

#### func (*Provider) GetVolumes

```go
func (p *Provider) GetVolumes(parent string) ([]Volume, error)
```
GetVolumes returns all NexentaStor volumes by parent volumeGroup

#### func (*Provider) GetVolumesSlice

```go
func (p *Provider) GetVolumesSlice(parent string, limit, offset int) ([]Volume, error)
```
GetVolumesSlice returns a slice of filesystems by parent filesystem with
specified limit and offset offset - the first record number of collection, that
would be included in result

#### func (*Provider) GetVolumesWithStartingToken

```go
func (p *Provider) GetVolumesWithStartingToken(parent string, startingToken string, limit int) (
	volumes []Volume,
	nextToken string,
	err error,
)
```
GetVolumesWithStartingToken returns filesystems by parent filesystem after
specified starting token parent - parent filesystem's path startingToken - a
path to a specific filesystem to start AFTER this token limit - the maximum
count of filesystems to return in the list Function may return nextToken if
there is more filesystems than limit value

#### func (*Provider) IsJobDone

```go
func (p *Provider) IsJobDone(jobID string) (bool, error)
```
IsJobDone checks if job is done by jobId

#### func (*Provider) LogIn

```go
func (p *Provider) LogIn() error
```
LogIn logs in to NexentaStor API and get auth token

#### func (*Provider) PromoteFilesystem

```go
func (p *Provider) PromoteFilesystem(path string) error
```
PromoteFilesystem promotes a cloned filesystem to be no longer dependent on its
original snapshot

#### func (*Provider) SetFilesystemACL

```go
func (p *Provider) SetFilesystemACL(path string, aclRuleSet ACLRuleSet) error
```
SetFilesystemACL sets filesystem ACL, so NFS share can allow user to write w/o
checking UNIX user uid

#### func (*Provider) String

```go
func (p *Provider) String() string
```

#### func (*Provider) UpdateFilesystem

```go
func (p *Provider) UpdateFilesystem(path string, params UpdateFilesystemParams) error
```
UpdateFilesystem updates filesystem by path

#### func (*Provider) UpdateVolume

```go
func (p *Provider) UpdateVolume(path string, params UpdateVolumeParams) error
```
UpdateVolume updates volume by path

#### type ProviderArgs

```go
type ProviderArgs struct {
	Address  string
	Username string
	Password string
	Log      *logrus.Entry

	// InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
	InsecureSkipVerify bool
}
```

ProviderArgs - params to create Provider instance

#### type ProviderInterface

```go
type ProviderInterface interface {
	// system
	LogIn() error
	IsJobDone(jobID string) (bool, error)
	GetLicense() (License, error)
	GetRSFClusters() ([]RSFCluster, error)

	// pools
	GetPools() ([]Pool, error)

	// filesystems
	CreateFilesystem(params CreateFilesystemParams) error
	UpdateFilesystem(path string, params UpdateFilesystemParams) error
	DestroyFilesystem(path string, params DestroyFilesystemParams) error
	SetFilesystemACL(path string, aclRuleSet ACLRuleSet) error
	GetFilesystem(path string) (Filesystem, error)
	GetFilesystemAvailableCapacity(path string) (int64, error)
	GetFilesystems(parent string) ([]Filesystem, error)
	GetFilesystemsWithStartingToken(parent string, startingToken string, limit int) ([]Filesystem, string, error)
	GetFilesystemsSlice(parent string, limit, offset int) ([]Filesystem, error)

	// filesystems - nfs share
	CreateNfsShare(params CreateNfsShareParams) error
	DeleteNfsShare(path string) error

	// filesystems - smb share
	CreateSmbShare(params CreateSmbShareParams) error
	DeleteSmbShare(path string) error
	GetSmbShareName(path string) (string, error)

	// snapshots
	CreateSnapshot(params CreateSnapshotParams) error
	DestroySnapshot(path string) error
	GetSnapshot(path string) (Snapshot, error)
	GetSnapshots(volumePath string, recursive bool) ([]Snapshot, error)
	CloneSnapshot(path string, params CloneSnapshotParams) error
	PromoteFilesystem(path string) error

	// volumes
	CreateVolume(params CreateVolumeParams) error
	GetVolume(path string) (Volume, error)
	GetVolumes(parent string) ([]Volume, error)
	UpdateVolume(path string, params UpdateVolumeParams) error
	DestroyVolume(path string, params DestroyVolumeParams) error
	GetVolumeGroup(path string) (VolumeGroup, error)
	GetVolumesWithStartingToken(parent string, startingToken string, limit int) ([]Volume, string, error)

	// iSCSI
	CreateLunMapping(params CreateLunMappingParams) error
	GetLunMapping(path string) (LunMapping, error)
	DestroyLunMapping(id string) error
	CreateISCSITarget(params CreateISCSITargetParams) error
	CreateUpdateTargetGroup(params CreateTargetGroupParams) error
}
```

ProviderInterface - NexentaStor provider interface

#### func  NewProvider

```go
func NewProvider(args ProviderArgs) (ProviderInterface, error)
```
NewProvider creates NexentaStor provider instance

#### type RSFCluster

```go
type RSFCluster struct {
	Name string `json:"clusterName"`
}
```

RSFCluster - RSF cluster with a name

#### type Resolver

```go
type Resolver struct {
	Nodes []ProviderInterface
	Log   *logrus.Entry
}
```

Resolver - NexentaStor cluster API provider

#### func  NewResolver

```go
func NewResolver(args ResolverArgs) (*Resolver, error)
```
NewResolver creates NexentaStor resolver instance based on configuration

#### func (*Resolver) IsCluster

```go
func (r *Resolver) IsCluster() (bool, error)
```
IsCluster checks if nodes is a NS cluster For now it simple checks if all nodes
return at least one similar cluster name

#### func (*Resolver) Resolve

```go
func (r *Resolver) Resolve(path string) (ProviderInterface, error)
```
Resolve returns one NS from the list of NSs by provided pool/dataset/fs path

#### func (*Resolver) ResolveFromVg

```go
func (r *Resolver) ResolveFromVg(path string) (ProviderInterface, error)
```
Resolve returns one NS from the list of NSs by provided pool/volumeGroup path

#### type ResolverArgs

```go
type ResolverArgs struct {
	Address  string
	Username string
	Password string
	Log      *logrus.Entry

	// InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
	InsecureSkipVerify bool
}
```

ResolverArgs - params to create resolver instance from config

#### type Snapshot

```go
type Snapshot struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Parent       string    `json:"parent"`
	Clones       []string  `json:"clones"`
	CreationTxg  string    `json:"creationTxg"`
	CreationTime time.Time `json:"creationTime"`
}
```

Snapshot - NexentaStor snapshot

#### func (*Snapshot) String

```go
func (snapshot *Snapshot) String() string
```

#### type UpdateFilesystemParams

```go
type UpdateFilesystemParams struct {
	// filesystem referenced quota size in bytes
	ReferencedQuotaSize int64 `json:"referencedQuotaSize,omitempty"`
}
```

UpdateFilesystemParams - params to update filesystem

#### type UpdateTargetGroupParams

```go
type UpdateTargetGroupParams struct {
	Members []string `json:"members"`
}
```

UpdateTargetGroupParams - params to update existing target group

#### type UpdateVolumeParams

```go
type UpdateVolumeParams struct {
	// volume referenced quota size in bytes
	VolumeSize int64 `json:"volumeSize,omitempty"`
}
```

UpdateVolumeParams - params to update volume

#### type Volume

```go
type Volume struct {
	Path           string `json:"path"`
	BytesAvailable int64  `json:"bytesAvailable"`
	BytesUsed      int64  `json:"bytesUsed"`
	VolumeSize     int64  `json:"volumeSize"`
}
```

Volume - NexentaStor volume

#### type VolumeGroup

```go
type VolumeGroup struct {
	Path           string `json:"path"`
	BytesAvailable int64  `json:"bytesAvailable"`
	BytesUsed      int64  `json:"bytesUsed"`
}
```

VolumeGroup - NexentaStor volumeGroup
