api {
  host = "0.0.0.0"
  port = 8000
}

enable_auth = true

cache {
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

// common-different configs with labels
rabbitmq {
  consumer_tag = "consumer"

  common {
    exchange = "service.direct"
    exchange_type = "direct"
  }
}

rabbitmq "master" {
  host = "rabbitmq:5672"
  user = "master"
  password = "password"
}

rabbitmq "local" {
  host = "0.0.0.0:5672"
  user = "bitnami"
  password = "password"
}