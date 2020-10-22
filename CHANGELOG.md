
<a name="v2.5.64"></a>
## [v2.5.64](https://github.com/Nexenta/go-nexentastor/compare/v2.5.63...v2.5.64) (0001-01-01)


<a name="v2.5.63"></a>
## [v2.5.63](https://github.com/Nexenta/go-nexentastor/compare/v2.5.62...v2.5.63) (2020-10-22)

### Pull Requests

* Merge pull request [#2](https://github.com/Nexenta/go-nexentastor/issues/2) from Nexenta/host_groups


<a name="v2.5.62"></a>
## [v2.5.62](https://github.com/Nexenta/go-nexentastor/compare/v2.5.6.1...v2.5.62) (2020-09-22)

### Bug Fixes

* ignore EEXIST on hostGroup creation


<a name="v2.5.6.1"></a>
## [v2.5.6.1](https://github.com/Nexenta/go-nexentastor/compare/v2.5.61...v2.5.6.1) (2020-09-22)


<a name="v2.5.61"></a>
## [v2.5.61](https://github.com/Nexenta/go-nexentastor/compare/v2.5.5...v2.5.61) (2020-09-22)


<a name="v2.5.5"></a>
## [v2.5.5](https://github.com/Nexenta/go-nexentastor/compare/v2.5.6...v2.5.5) (2020-07-17)


<a name="v2.5.6"></a>
## [v2.5.6](https://github.com/Nexenta/go-nexentastor/compare/v2.5.4...v2.5.6) (2020-07-17)

### Bug Fixes

* fixed GetVolumesWithStartingToken


<a name="v2.5.4"></a>
## [v2.5.4](https://github.com/Nexenta/go-nexentastor/compare/v2.5.3...v2.5.4) (2020-07-17)

### Features

* implemented GetVolumesWithStartingToken


<a name="v2.5.3"></a>
## [v2.5.3](https://github.com/Nexenta/go-nexentastor/compare/v2.5.2...v2.5.3) (2020-07-17)

### Bug Fixes

* changed base offset to 0 for getVolumes


<a name="v2.5.2"></a>
## [v2.5.2](https://github.com/Nexenta/go-nexentastor/compare/v2.5.1...v2.5.2) (2020-07-16)

### Features

* added GetVolumes interface for listing volumes


<a name="v2.5.1"></a>
## [v2.5.1](https://github.com/Nexenta/go-nexentastor/compare/v2.5.0...v2.5.1) (2020-07-16)

### Features

* added getVolume interface, changed getVolumeGroup to return vg props


<a name="v2.5.0"></a>
## [v2.5.0](https://github.com/Nexenta/go-nexentastor/compare/v2.4.0...v2.5.0) (2020-06-30)


<a name="v2.4.0"></a>
## [v2.4.0](https://github.com/Nexenta/go-nexentastor/compare/v2.3.0...v2.4.0) (2020-06-30)


<a name="v2.3.0"></a>
## [v2.3.0](https://github.com/Nexenta/go-nexentastor/compare/v2.1.0...v2.3.0) (2020-03-11)


<a name="v2.1.0"></a>
## [v2.1.0](https://github.com/Nexenta/go-nexentastor/compare/v2.0.0...v2.1.0) (2020-03-11)

### Bug Fixes

* fixed changelog

### Features

* added UpdateFilesystem to be able to expand volume


<a name="v2.0.0"></a>
## [v2.0.0](https://github.com/Nexenta/go-nexentastor/compare/v1.2.0...v2.0.0) (2020-03-11)

### Bug Fixes

* add unmarshal error text to rest error response

### Features

* added UpdateFilesystem to be able to expand volume
* NEX-22162 ability to pass nfs list on volume creation
* add parameters to provider.DestroyFilesystem() method

### BREAKING CHANGE


`provider.DestroyFilesystemWithClones()` method was removed, use `DestroyFilesystemParams.PromoteMostRecentCloneIfExists` instead.


<a name="v1.2.0"></a>
## [v1.2.0](https://github.com/Nexenta/go-nexentastor/compare/v1.1.0...v1.2.0) (2019-07-18)

### Features

* add provider.GetFilesystemsWithStartingToken() method


<a name="v1.1.0"></a>
## [v1.1.0](https://github.com/Nexenta/go-nexentastor/compare/v1.0.1...v1.1.0) (2019-07-10)

### Features

* add nefError.IsBadArgNefError() method


<a name="v1.0.1"></a>
## [v1.0.1](https://github.com/Nexenta/go-nexentastor/compare/v1.0.0...v1.0.1) (2019-05-07)

### Bug Fixes

* NEX-20603 - filesystem list returns only first 100 items


<a name="v1.0.0"></a>
## v1.0.0 (2019-04-19)

### Features

* add vendors to versioning
* initial NS API Go library commit

