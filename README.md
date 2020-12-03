# Setup

Clone:
```
git clone git@github.com:carsonip/b2h-2020-got-swagger.git
```

Add magic comment to target repo (`~/dev/pendo-appengine/src/pendo.io/server/server.go:469`):
```
/* HACK HERE */
```

Add module path to target repo (`~/dev/pendo-appengine/src/go.mod`):
```
require github.com/carsonip/b2h v0.0.0
replace github.com/carsonip/b2h v0.0.0 => /path/to/b2h-2020-got-swagger
```

Run `go mod vendor` in appengine:
```
cd ~/dev/pendo-appengine/src
go mod vendor
```

Build the tool:
```
cd /path/to/b2h-2020-got-swagger
go build .
```

Run it:
```
./martiniExample match -m get -p /api/dashboard/share/list
```