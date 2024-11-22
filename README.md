[![Go Report Card](https://goreportcard.com/badge/github.com/indeedhat/automux)](https://goreportcard.com/report/github.com/indeedhat/automux)
[![codecov](https://codecov.io/gh/indeedhat/automux/graph/badge.svg?token=5M4H16EDWM)](https://codecov.io/gh/indeedhat/automux)

# Automux
Auto mux checks for a .automux.hcl on each cd and if one is found it will automatically run tmux with the layout provided by the config file.

This is a simple project, there are other projects out there that take the same concept a lot further.  
I personally like simpler tools so i created automux.

![demo.gif](_examples/demo-24-11-22.gif)

## Insall
```sh
go install github.com/indeedhat/automux@latest

# if you would like automux to run on cd without having to be manually called:
# add the following function to your .bashrc (or whatever your shells rc file is)
cd() {
    builtin cd "$@" && automux
}

# Another function i have found useful is to use :qa to close kill the entire tmux session
:qa() {
    if [ -n "$TMUX" ]; then
        tmux kill-session
    else
        exit
    fi
}
```

## Features
### Windows
- can open one or more tabs
- can give names to tabs
- can auto run command on open
- specific windows can be focused on open
- can be given a sub directory to open in

### Splits
- windows can have one or more splits
- each split can be set to desired size
- splits can each run a command on open
- specific splits can be focused on open
- can be given a sub directory to open in, this will always be relative to the base session, not the window the split is defined within

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
session_id = "my-session"

# config lets you set a custom tmux config for this directory
config = "./tmux.conf"

# when set automux will open tmux and attach to the existing session for the directory (if one exists)
# when not set automux will do nothing if a session exists
attach_existing = false # default true

# the first window block will setup the original window/tab
# each additional block will add a new window/tab
window "window/tab title" {
    exec = "cmd_to_run_in_window"
    # focus can be set for any window/split
    # once the setup is done focus will be set to the last window/split in the config file that 
    # has focus = true
    focus = true

    # The sub directory to open the window in
    dir = "sub_dir/"

    split {
        vertical = true
        exec = "cmd_to_run_in_split"

        # set the size of the split in % of the total screen size
        # vertical splits will set the height, horizontal the width
        #
        # NOTE: the size is set at create time so will work fine for simple layouts but may cause issues
        #       trying to create more complex ones
        size = 30

        # The sub directory to open the split in (relative to the containing window)
        dir = "sub_dir/"
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

    # session_id = "my-session"
    # config = "./tmux.conf"
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

## Upgrade
If you are coming from an older version of automux it was configured with a `.automux.hcl` file,  
This has been updated to use an icl file called `.automux`.

It is recommended that you rename your `.automux.hcl` file to `.automux`, the former will still be picked up
but it is deprecated and support will be dropped in a future version.

### File contents changes
- add `version = 1` as the very first non config line in your file
- convert any `session = "..."` lines to `session_id = "..."`

### tmux-sessionizer.sh
if you are using the tmux-sessionizer script provided in the repo it will need to be updated to the latest version
