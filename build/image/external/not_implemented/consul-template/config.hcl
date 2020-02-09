consul {
    address = "localhost:8500"
    retry {
        enabled = true
        attempts = 12
        backoff = "250ms"
    }
    token = "w94RIMKUtQH1a4VJGN+t+vn1Y0nErc/ch93E1F1ZcHU="
}
reload_signal = "SIGHUP"
kill_signal = "SIGINT"
max_stale = "10m"
log_level = "warn"
# pid_file = "/consul-template/consul-template.pid"
wait {
    min = "5s"
    max = "10s"
}
vault {
    address = "http://localhost:8200"
    token = "R/Uf0tYa5YkhPLpNLL807KWJ4ZiJi3clyQEfaMoRSJg"
    renew_token = false
}
deduplicate {
    enabled = true
    # prefix = "consul-template/dedup/"
}

template {
    source      = "some.tpl"
    destination = "some.txt"
    left_delimiter  = "{{"
    right_delimiter = "}}"
    wait {
          min = "2s"
          max = "10s"
    }
}
// template {
//     source      = "./vault/templates/pki/ca.ctmpl"
//     destination = "./vault/output/pki/mpatel.yourdomain.com.ca.crt"
// }
// template {
//     source      = "./vault/templates/pki/key.ctmpl"
//     destination = "./vault/output/pki/mpatel.yourdomain.com.key"
// }