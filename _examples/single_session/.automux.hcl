session = "my-single-session"

window "Text Editor" {
    exec = "nvim"
    focus = true
}

window "System Monitor" {
    exec = "htop"
    split {
        exec = "nload"
        vertical = false
    }
}
