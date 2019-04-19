# rest
--
    import "github.com/Nexenta/go-nexentastor/pkg/rest"


## Usage

#### type Client

```go
type Client struct {
}
```

Client - request client for any REST API

#### func (*Client) BuildURI

```go
func (c *Client) BuildURI(uri string, params map[string]string) string
```
BuildURI builds request URI using [path?params...] format

#### func (*Client) Send

```go
func (c *Client) Send(method, path string, data interface{}) (int, []byte, error)
```
Send sends request to REST server data interface{} - request payload, any
interface for json.Marshal()

#### func (*Client) SetAuthToken

```go
func (c *Client) SetAuthToken(token string)
```
SetAuthToken sets Bearer auth token for all requests

#### type ClientArgs

```go
type ClientArgs struct {
	Address string
	Log     *logrus.Entry

	// InsecureSkipVerify controls whether a client verifies the server's certificate chain and host name.
	InsecureSkipVerify bool
}
```

ClientArgs - params to create Client instance

#### type ClientInterface

```go
type ClientInterface interface {
	BuildURI(uri string, params map[string]string) string
	Send(method, path string, data interface{}) (int, []byte, error)
	SetAuthToken(token string)
}
```

ClientInterface - request client interface

#### func  NewClient

```go
func NewClient(args ClientArgs) ClientInterface
```
NewClient creates new REST client
