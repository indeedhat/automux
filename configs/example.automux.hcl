# the session id to use for this directory
# NOTE: this is the only required field
session = "my-session"

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
