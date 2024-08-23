[![Go Report Card](https://goreportcard.com/badge/github.com/indeedhat/automux)](https://goreportcard.com/report/github.com/indeedhat/automux)
[![codecov](https://codecov.io/gh/indeedhat/automux/graph/badge.svg?token=5M4H16EDWM)](https://codecov.io/gh/indeedhat/automux)

# Automux
Auto mux checks for a .automux.hcl on each cd and if one is found it will automatically run tmux with the layout provided by the config file.

This is a simple project, there are other projects out there that take the same concept a lot further.  
I personally like simpler tools so i created automux.

![demo.gif](_examples/demo.gif)

## Insall
```sh
go install github.com/indeedhat/automux@latest

# if you would like automux to run on cd without having to be manually called:
# add the following function to your .bashrc (or whatever your shells rc file is)
cd() {
    builtin cd "$@" && automux 
}
```

## Features
### Windows
- can open one or more tabs
- can give names to tabs
- can auto run command on open
- specific windows can be focused on open

### Splits
- windows can have one or more splits
- each split can be set to desired size
- splits can each run a command on open
- specific splits can be focused on open

### Background Sessions
When opening the main automux session you can optionally open one or more
background sessions, allowing you to open multiple related projects at once ready to
be focused from a single terminal window at will.

## Usage
```
automux -h
automux [flags] [path]
Usage of automux:
  -d    Run the automux session detached
        This will allow you to start an automux session from another session
  -debug
        print tmux commands rather than running them
  -init
        Init the automux config template in the current directory
  -print-name
        Print the session name if the target directory is a automux directory
```

## Configure
```hcl
# the session id to use for this directory
# NOTE: this is the only required field
session = "my-session"

# config lets you set a custom tmux config for this directory
config = "./tmux.conf"

# when set automux will not run if there is already a tmux session with the provided {session}
single_session = false 

# the first window block will setup the original window/tab
# each additional block will add a new window/tab
window "window/tab title" {
    exec = "cmd_to_run_in_window"
    # focus can be set for any window/split
    # once the setup is done focus will be set to the last window/split in the config file that 
    # has focus = true
    focus = true

    split {
        vertical = true
        exec = "cmd_to_run_in_split"

        # set the size of the split in % of the total screen size
        # vertical splits will set the height, horizontal the width
        #
        # NOTE: the size is set at create time so will work fine for simple layouts but may cause issues
        #       trying to create more complex ones
        size = 30
    }
}

# can have multiple windows/tabs
window "vim" {
    exec = "nvim"

    # can have multiple splits
    split {}
    split {
        exec = "nload"
        vertical = true
    }
}

# sub sessions will be opened in the background
session "path/to/session_dir" {
    # if a .automux.hcl file is found in the session dir then it will be loaded
    # any config put in the session block will overwrite config found there

    # session = "my-session"
    # config = "./tmux.conf"
    # single_session = false 

    window "window_name" {
        # if a window with the same name is found in the .autmux.hcl file then the two blocks will be
        # merged with any values set here taking presedence

        # exec = ""
        # focus = false

        spit {
            # splits will be merged by index
            # with any values set here taking presedence

            # any aditional splits will be appended to the final window config
        }
    }
}
```

## TODO
- [ ] find a way to detach the session from the go app so it can shutdown without killing the sessien
- [x] have a way to auto focus a given split/window on load
