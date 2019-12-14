# img - Command-line image viewer

A command line tool to view images (PNG, JPEG, GIF) right on the terminal. `img` comes in handy in the following scenarios:
- to view images over SSH and VPN connections (where it's cumbersome to grab images and view them on the host machine)
- can be used to generate splash screens for Linux logins (e.g. motd)
- you never have to leave the terminal if you are working with image generation code
- just for fun!

#### GIF
<img src="resources/readme/animated.gif" alt="GIF demo" width="65%" height="65%"/>

#### JPEG/PNG
<img src="resources/readme/static.gif" alt="Static image demo" width="65%" height="65%"/>

## Installation

#### macOS
```
brew tap codeliveroil/apps
brew install img
```

#### Other
Download the [latest release](../../releases/latest) for your operating system and machine architecture. If one is not available, you can very easily [compile from source](#compile-from-source).

## Usage
```
img -help
```

#### Examples
```
img car.png
img -w logo.sh logo.gif
img -l 2 wheel.gif
```

## Library API for Go

```golang
img := viz.Image{
	Filename:  "input.gif",
	LoopCount: 2,
}

// Read the image
if err := img.Init(); err != nil {
	//handle error
}

// Render the image
if err := img.Draw(&viz.StdoutCanvas{}); err != nil {
	//handle error
}
```


## Compile from source

### Setup
1. Install [Go](https://golang.org/)
1. Clone this repository

### Build for your current platform
```
make
make install
```

### Cross compile for a different platform
1. Build
	```
	make platform
	```
1. Follow the prompts and select a platform
1. The binary will be available in the `build` folder
