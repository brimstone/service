service
=======

This is a little toy service manager for my Go projects.

## Usage

Create a implementation of the Run interface
```go
type svc struct {
// This could be anything your service needs
// Nothing needs to be exported
}

func (s *svc) Run (ctx context.Context) error {

}
```

Then add it to the manager singleton
```go
service.Manager.Add(&svc)
```

Once all of the services have been added, start them
```go
ctx, allDone := context.WithCancel(context.Background())
service.Manager.RunAll(ctx)
```

When it's time to stop all of the services, cancel the context or call `service.manager.StopAll()`
```go
alldone()
```
