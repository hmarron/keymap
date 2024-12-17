# Keymap

### Description
Generates images/wallpapers to display key map configuration.
![2024-12-17_16-38-04](https://github.com/user-attachments/assets/7e3a6b73-9849-48af-899a-7b331c848404)

### Usage
1. Set up styling, size of image, etc in the config file
2. Set up frames and key maps in config file (each frame is a block of key maps with a title)
3. Generate the image `go run cmd/keymap/main.go`

### Getting the key bindings
Different programs have different ways of getting their key bindings programatically.  
I've provided the following for getting key maps for tmux, nvim and aerospace respectively
- `cmd/tmux/main.go`
- `cmd/nvim/main.go`
- `cmd/aerospace.go`

Each of these can take a `--config-file config.yml` flag that will attempt to override any frame definied in your config.yml
where the frame name is matches `tmux`, `nvim`, `aerospace-main`, `aerospace-service`

It is expected that you may want to edit the generated key map config before generating the image, to tweak things just as you want


### Known issues
I'm using `github.com/InfinityTools/go-binpack2d` for rectangle packing.  
This seems to only work sometimes, and sometimes fails to pack all frames into the image.  

If this happens, it's worth re-running a few times to see if it can find a fit, otherwise change the styles (font size, padding, etc) so the generated rectangles will be smaller
