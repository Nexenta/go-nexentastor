
<a name="v2.0.0"></a>
## [v2.0.0](https://github.com/Nexenta/go-nexentastor/compare/v1.2.0...v2.0.0) (2019-07-25)

### Bug Fixes

* add unmarshal error text to rest error response

### Features

* add parameters to provider.DestroyFilesystem() method

### BREAKING CHANGE


`provider.DestroyFilesystemWithClones()` method has been removed, use `DestroyFilesystemParams.PromoteMostRecentCloneIfExists` instead.


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
