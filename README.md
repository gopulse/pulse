<p align="center">
    <img src="https://user-images.githubusercontent.com/32125808/231166226-c636f344-fc3b-4c71-9181-49fe97491127.png" width="320px">
</p>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/gopulse/pulse)](https://goreportcard.com/report/github.com/gopulse/pulse)
[![GitHub license](https://img.shields.io/github/license/gopulse/pulse)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/gopulse/pulse.svg)](https://pkg.go.dev/github.com/gopulse/pulse)
[![Go Doc](https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-00ACD7.svg?color=00ACD7&style=flat-square)](https://pkg.go.dev/github.com/gopulse/pulse#pkg-overview)
[![Discord Online](https://img.shields.io/discord/1095400462477426748)](https://discord.gg/JKcTwZYJ)
[![codecov](https://img.shields.io/codecov/c/github/gopulse/pulse?token=RBXPY1WN2I)](https://codecov.io/github/gopulse/pulse)
[![CircleCI](https://img.shields.io/circleci/build/github/gopulse/pulse/master?token=7eda4a74e26b544956b8333b372592ee09cd7f8b)](https://dl.circleci.com/status-badge/redirect/gh/gopulse/pulse/tree/master)

</div>

A **Golang** framework for web development that keeps your web applications and services **alive** and responsive with its fast and lightweight design.
## Features

- Routing
- Route groups
- Static files
- Simple and elegant API
- Middleware support
- Validation
- Routes grouping

### Installation

Make sure you have Go installed on your machine. Then run the following command:

Initialize your project ([Learn](https://go.dev/blog/using-go-modules)). Then install **Pulse** with the go get command:

```bash
go get github.com/gopulse/pulse
```

### Getting Started

```go
package main

import (
    "github.com/gopulse/pulse"
)

func main() {
    app := pulse.New()
	router := pulse.NewRouter()

	app.Router = router

	router.Get("/", func(c *pulse.Context) error {
        c.String("Hello, World!")
		return nil
    })

    app.Run(":3000")
}
```

### Examples

- Routing

Supports `GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE`

```go
package main

import (
	"github.com/gopulse/pulse"
)

func main() {
    app := pulse.New()
    router := pulse.NewRouter()
    
    // GET /hello
    router.Get("/", func(c *pulse.Context) error {
        c.String("Hello, World!")
        return nil
    })
    
    // GET /hello/:name
    router.Get("/profile/:id", func(c *pulse.Context) error {
        c.String("Profile: " + c.Param("id"))
        return nil
    })
    
	// GET /user/
    router.Get("/user/*", func(c *pulse.Context) error {
        c.String("Hello, World!")
        return nil
    })
    
    app.Router = router
    
    app.Run(":3000")
}
```

- Route groups

Supports `GET, POST, PUT, PATCH, DELETE, OPTIONS, HEAD, CONNECT, TRACE`

```go
package main

import (
	"github.com/gopulse/pulse"
)

func main() {
    app := pulse.New()
    router := pulse.NewRouter()
	api := &pulse.Group{
		prefix: "/api",
		router: router,
	}

	v1 := api.Group("/v1")
	v1.GET("/users", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
    
    app.Router = router
    
    app.Run(":3000")
}
```

* Static files

```go
package main

import (
	"github.com/gopulse/pulse"
	"time"
)

func main() {
	app := pulse.New()
	router := pulse.NewRouter()

	// Static files (./static) with cache duration 24 hours
	router.Static("/", "./static", &pulse.Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})

	app.Router = router

	app.Run(":3000")
}
```

* Middleware

```go
package main

import (
	"github.com/gopulse/pulse"
)

func main() {
	app := pulse.New()
	router := pulse.NewRouter()
	
	router.Get("/profile/:name", func(ctx *pulse.Context) error {
		if ctx.Param("name") != "test" {
			ctx.Abort()
			ctx.Status(404)
			return nil
		}
		ctx.String("hello")
		ctx.Next()
		return nil
	})

	app.Router = router

	app.Run(":3000")
}
```

## Available Middleware

- [x] CORS Middleware: Enable cross-origin resource sharing (CORS) with various options.
```go
package main

import (
	"github.com/gopulse/pulse"
)

func main() {
	app := pulse.New()
	router := pulse.NewRouter()

	router.Get("/", func(ctx *pulse.Context) error {
		return nil
	})

	router.Use("GET", pulse.CORSMiddleware())

	app.Router = router

	app.Run(":3000")
}
```

- [ ] Logger Middleware: Log every request with configurable options. **(Coming soon)**
- [ ] Encrypt Cookie Middleware: Encrypt and decrypt cookie values. **(Coming soon)**
- [ ] Timeout Middleware: Set a timeout for requests. **(Coming soon)**

## License

Pulse is licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.

## Contributing

Contributions are welcome! Please read the [contribution guidelines](CONTRIBUTING.md) first.

## Support

If you want to say thank you and/or support the active development of Pulse:
1. Add a [GitHub Star](star) to the project.
2. Tweet about the project [on your Twitter](https://twitter.com/intent/tweet?text=Pulse%20is%20a%20%23web%20%23framework%20for%20the%20%23Go%20programming%20language.%20It%20is%20a%20lightweight%20framework%20that%20is%20%23easy%20to%20use%20and%20easy%20to%20learn.%20It%20is%20designed%20to%20be%20a%20simple%20and%20elegant%20solution%20for%20building%20web%20applications%20and%20%23APIs%20%F0%9F%9A%80%20https%3A%2F%2Fgithub.com%2Fgopulse%2Fpulse)
3. Write a review or tutorial on [Medium](https://medium.com/), [dev.to](https://dev.to/), [Reddit](https://www.reddit.com/) or personal blog.
4. [Buy Me a Coffee](https://www.buymeacoffee.com/gopulse)

## Contributors
<!-- CONTRIBUTORS-START -->
<a href="https://github.com/gopulse/pulse/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=gopulse/pulse&columns=18" />
</a>
<!-- CONTRIBUTORS-END -->

# Stargarazers over time

[![Stargazers over time](https://starchart.cc/gopulse/pulse.svg)](https://starchart.cc/gopulse/pulse)