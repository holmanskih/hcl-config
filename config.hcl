api {
  host = "0.0.0.0"
  port = 8000
}

enable_auth = true

cache "local" "master" {
  type = "redis"

  redis {
    dev_mode = true
    password = "bitnami"
    host = "0.0.0.0"
  }

  nutsdb {
    path = "nuts"
    segment_size = 1024
  }
}