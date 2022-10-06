# Extra Credit

1. ## Convert an empty interface{} to a typed var

I normally don't use interfaces. However,

1. ## Show an example of a Go routine and how to stop it

```go
func hello() {  
    fmt.Println("Hello...")
}
func main() {  
    go foo()
    time.Sleep(1 * time.Second)
    fmt.Println("world")
}
```

1. ## Create a simple in-memory cache that can safely be accessed concurrently

1. ## What is a slice and how is it related to an array?

- Slices are references to arrays. This concept is similar to static arrays and dynamic arrays in many other compiled languages such as cpp. Some differences that are probably unique to golang is that slices aren't garbage collected as long as the underlying array is still in use(in this case use `copy` to allow the reference to be deleted). Another is that you can mix types in a slice. Things that are common to slices and dynamic arrays are the ability to grow and shrink in size.

1. ## What is the syntax for zero allocation, swapping of two integers variables?

I'm not sure. I assume that swapping using the python style still allocates extra memory behind the scenes. Example of what I mean:

```go
    a,b:=1,2
    a,b=b,a //a=2 b=1
```
