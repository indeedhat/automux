session = "my-overrides-subsession"

window "editor" {
    exec = "nvim"
}

window "cmd" {
    exec = "uname -a"
    split {
        exec = "whoami"
    }
    split {
        exec =  "htop"
    }
}
