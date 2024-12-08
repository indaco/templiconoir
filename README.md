<h1 align="center" style="font-size: 2.5rem;">
  templiconoir
</h1>
<p align="center">
    <a href="https://pkg.go.dev/github.com/indaco/templiconoir/" target="_blank">
        <img src="https://pkg.go.dev/badge/github.com/indaco/templiconoir/.svg" alt="go reference" />
    </a>
    &nbsp;
    <a href="https://goreportcard.com/report/github.com/indaco/templiconoir" target="_blank">
        <img src="https://goreportcard.com/badge/github.com/indaco/templiconoir" alt="go report card" />
    </a>
    &nbsp;
    <a href="https://coveralls.io/github/indaco/templiconoir?branch=main">
        <img
            src="https://coveralls.io/repos/github/indaco/templiconoir/badge.svg?branch=main"
            alt="Coverage Status"
        />
    </a>
     &nbsp;
     <a href="https://github.com/indaco/templiconoir/blob/main/LICENSE" target="_blank">
        <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square&logo=none" alt="license" />
    </a>
    &nbsp;
    <a href="https://www.jetify.com/devbox/docs/contributor-quickstart/">
      <img
          src="https://www.jetify.com/img/devbox/shield_moon.svg"
          alt="Built with Devbox"
      />
    </a>
</p>

This package provides the [Iconoir](https://iconoir.com) set (_v7.10.1_) as reusable, type-safe go [templ](https://github.com/a-h/templ) components.

## Features

- **Lazy Loading**: Icons are loaded on demand at runtime, reducing memory usage and improving performance.
- **Customizable**: Easily adjust size, color, stroke-width, and add attributes with a simple, chainable API.
- **Memory Efficient**: Avoids preloading large datasets, reducing memory overhead.
- **Local Caching**: Speeds up icon with efficient local caching.

## Installation

Install the package using `go get`:

```bash
go get github.com/indaco/templiconoir@latest
```

## Icon Naming Convention

We categorize Iconoir based on their style (_Outline_, _Solid_). This is reflected in the naming convention for the components:

**1. Outline Icons**

- Default style with a size of _24px_.
- Example: `iconoir.Chromecast`, `iconoir.CheckCircle`.

**2. Solid Icons**

- Style is explicitly "solid" with a size of _24px_.
- Example: `iconoir.CheckCircleSolid`, `iconoir.ChatMinusInSolid`.

Icons are named in _PascalCase_ for consistency and ease of use. Size and style are embedded in the names to differentiate icons visually and programmatically.

## Usage

### Rendering Icons

To use the icons in your templ project, call the `Render()` method on the desired icon component:

```templ
package pages

import iconoir "github.com/indaco/templiconoir"

templ DemoPage() {
    @iconoir.Chromecast.Render()            // Outline 24px
    @iconoir.CheckCircleSolid.Render() // Solid 24px
}
```

### Customizing Icons

With the `ConfigureIcon` builder pattern, you can customize the icons more fluently and efficiently. This approach allows to chain multiple customization methods together. Once you've set your desired properties, call `Build()` to get the final icon, which you can render directly as templ component.

#### 1. SetSize()

Use the `SetSize()` method to set a custom size for the icon in pixels:

```templ
package pages

import iconoir "github.com/indaco/templiconoir"

templ CustomSizePage() {
    // Set custom size
    @iconoir.ConfigureIcon(iconoir.CheckCircleSolid).
        SetSize(32).
        Build().
        Render()
}
```

#### 2. SetColor()

Use the `SetColor()` method to modify the fill color for the icons:

```templ
package pages

import iconoir "github.com/indaco/templiconoir"

templ CustomFillColor() {
    // Customize fill color
   @iconoir.ConfigureIcon(iconoir.Chromecast).
       SetColor("#2dd4bf").
       Build().
       Render()
}
```

#### 3. SetStrokeWidth()

Use the `SetStrokeWidth()` method to modify the stroke-width for the icons:

```templ
package pages

import iconoir "github.com/indaco/templiconoir"

templ CustomStrokeWidthColor() {
    // Customize stroke.width
   @iconoir.ConfigureIcon(iconoir.Swimming).
       SetStrokeWidth("2").
       Build().
       Render()
}
```

#### 4. SetAttrs()

You can also use the `SetAttrs()` method to add custom attributes to the icons, such as _aria-hidden_, _focusable_, or custom CSS classes:

```templ
package pages

import iconoir "github.com/indaco/templiconoir"

templ CustomIconPage() {
    // Add attributes to an icon
    @iconoir.ConfigureIcon(iconoir.Swimming).
        SetAttrs(templ.Attributes{
            "aria-hidden": "true",
            "class":       "custom-icon",
        }).
        Build().
        Render()
}
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

### Development Environment Setup

To set up a development environment for this repository, you can use [devbox](https://www.jetify.com/devbox) along with the provided `devbox.json` configuration file.

1. Install devbox by following the instructions in the [devbox documentation](https://www.jetify.com/devbox/docs/installing_devbox/).
2. Clone this repository to your local machine.
3. Navigate to the root directory of the cloned repository.
4. Run `devbox install` to install all packages mentioned in the `devbox.json` file.
5. Run `devbox shell --pure` to start a new shell with access to the environment.
6. Once the devbox environment is set up, you can start developing, testing, and contributing to the repository.

### Running Tasks

This project provides both a `Makefile` and a `Taskfile` for running various tasks. You can use either `make` or `task` to execute the tasks, depending on your preference.

To view all available tasks, run:

- **Makefile**: `make help`
- **Taskfile**: `task --list-all`

Available tasks:

```bash
build                   # Generate the Go icon definitions based on parsed data/heroicons_cache.json file.
demo:                   # Run the demo server.
test                    # Run go tests.
test/coverage:          # Run go tests and use go tool cover.
```

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.
