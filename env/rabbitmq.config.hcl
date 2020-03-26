rabbitmq {
  consumer_tag = "consumer"

  common {
    exchange = "service.direct"
    exchange_type = "direct"
  }
}

rabbitmq "local" {
  host = "0.0.0.0:5672"
  user = "bitnami"
  password = "password"
}

rabbitmq "master" {
  host = "rabbitmq:5672"
  user = "master"
  password = "password"
}