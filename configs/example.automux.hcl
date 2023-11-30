session = "my-session"
# when set automux will not run if there is already a tmux session with the provided {session}
single_session = false

window "window/tab title" {
    exec = "cmd_to_run_in_window" (optional)

    split {
        vertical = true # (optional)
        exec = "cmd_to_run_in_split" # (optional)
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
