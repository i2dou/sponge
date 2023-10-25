## errcode

Error codes usually include system-level error codes and business-level error codes, consisting of a total of 5 decimal digits, e.g. 20101

| First digit                                                                          | Middle two digits | Last two digits |
|:-------------------------------------------------------------------------------------|:-------|:-------|
| For http error codes, 2 indicates a business-level error (1 is a system-level error) | Service Module Code | Specific error codes |
| For grpc error codes, 4 indicates a business-level error (3 is a system-level error) | Service Module Code | Specific error codes |

- Error levels occupy one digit: 1 (http) and 3 (grpc) indicate system-level errors, 2 (http) and 4 (grpc) indicate business-level errors, usually caused by illegal user operations.
- Double-digit service modules: A large system usually has no more than two service modules; if it exceeds that, it's time to split the system.
- Error codes take up two digits: prevents a module from being customised with too many error codes, which are not well maintained later.

<br>

### Example of use

### Example of http error code usage

```go
    import "github.com/i2dou/sponge/pkg/errcode"

    // defining error codes
    var ErrLogin = errcode.NewError(20101, "incorrect username or password")

    // return error
    response.Error(c, errcode.LoginErr)
```

<br>

### Example of grpc error code usage

```go
    import "github.com/i2dou/sponge/pkg/errcode"

    // defining error codes
    var ErrLogin = errcode.NewRPCStatus(40101, "incorrect username or password")

    // return error
    errcode.ErrLogin.Err()
    // return with error details
    errcode.ErrLogin.Err(errcode.Any("err", err))
```