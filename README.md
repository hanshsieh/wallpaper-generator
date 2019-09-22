# Come on! Wallpaper!

## Build
```
go build ./cmd/comeonwallpaper/
```
To do cross compile, (e.g. building for Windows), use:
```
GOOS=windows GOARCH=386 go build ./cmd/comeonwallpaper/
```

## Usage
```
./comeonwallpaper <src_dir> <dest_dir> <width> <height>
```
Example:
```
./comeonwallpaper ./srcdir ./destdir 1920 768
```
