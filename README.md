# Automux
Auto mux checks for a .automux.hcl on each cd and if one is found it will automatically run tmux with the layout provided by the config file.

Its a pretty silly project and there are many others that can do a lot more with the concept but i had the thought and had to make it

![demo.gif](_examples/demo.gif)

## insall
```sh
go install github.com/indeedhat/automux@latest

# if you would like automux to run on cd without having to be manually called:
# add the following function to your .bashrc (or whatever your shells rc file is)
cd() {
    builtin cd "$@" && automux 
}
```

## Usage
```
automux -h
Usage of automux:
  -debug
        print tmux commands rather than running them
  -init
        Init the automux config template in the current directory
```

## Configure
[See here](configs/example.automux.hcl)

## TODO
- [ ] find a way to detach the session from the go app so it can shutdown without killing the sessien
- [x] have a way to auto focus a given split/window on load
