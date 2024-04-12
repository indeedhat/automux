session = "my-multi-session"

window "editor" {
    exec = "nvim"
    focus = true
}

window "cmd" {
    split {}
}

# sub sessions
session "./untouched/" {}
session "./with_overrides/" {
    window "editor" {
        exec = ""
    }

    window "cmd" {
        split {
            focus = true
        }
    }
}
session "./no_config_file/" {
    session = "my-no-config-subsession"

    window "editor" {
        exec = "nvim"
    }

    window "htop" {
        exec =  "htop"
        focus = true
    }
}

