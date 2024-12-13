package main

func (s *server) routes() {
	s.Router.Post("/", s.broker)
	s.Router.Post("/authentication", s.authentication)
	s.Router.Post("/grpc-logger", s.gRPCLogger)
	s.Router.Post("/rabbitmq-authentication", s.rabbitMQAuthentication)
}
