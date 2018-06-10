# sweet-junit #
A plugin for the Sweet testing framework to output JUnit files for test results

## Usage ##

Using this plugin with [sweet](https://www.github.com/aphistic/sweet) is pretty
straightforward:

```go
func Test(t *testing.T) {
    sweet.T(func(s *sweet.S) {
        s.RegisterPlugin(junit.NewPlugin())

        s.RunSuite(t, &mySuite{})
    })
}

```

Once the plugin is registered with sweet, you can specify the file to write the
output to by passing the `-sweet.opt` when running `go test` and providing the
`junit.output` key with the path you'd like to write the junit file to, such as:

```bash
$ go test -sweet.opt "junit.output=junit.xml"
```
