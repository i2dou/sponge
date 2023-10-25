## group

group libraries.

### Example of use

```go
    import "github.com/i2dou/sponge/pkg/container/group"

    type foo struct {
        bar string
    }
    
    gr := group.NewGroup(func () interface{} {
        return &foo{"hello"}
    })

	fmt.Println(gr.Get(*foo).bar)
```
