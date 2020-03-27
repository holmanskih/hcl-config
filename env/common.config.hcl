api {
  host = "0.0.0.0"
  port = 89000
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

rabbitmq "local" "master" {
  consumer_tag = "consumer"

  common {
    exchange = "service.direct"
    exchange_type = "direct"
  }
}