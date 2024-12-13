from locust import HttpUser, task, between

class MicroserviceLoadTest(HttpUser):
    wait_time = between(1, 5)  # Espera entre 1 y 5 segundos entre cada solicitud

    @task
    def test_broker(self):
        payload = {"content": "Empty request"}
        self.client.post("/broker", json=payload)  # Cambia la URL según tu endpoint

    @task
    def test_authentication(self):
        payload = {"email": "admin@example.com", "password": "password123"}
        self.client.post("/authentication", json=payload)  # Cambia la URL según tu endpoint

    @task
    def test_rabbitmq_auth(self):
        payload = {"email": "admin@example.com", "password": "password123"}
        self.client.post("/rabbitmq-authentication", json=payload)  # Cambia la URL según tu endpoint

    @task
    def test_logger(self):
        payload = {"name": "activity", "data": "some kind of grpc data"}
        self.client.post("/grpc-logger", json=payload)  # Cambia la URL según tu endpoint
