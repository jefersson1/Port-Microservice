import { useState } from "react";

// Función principal del componente TestArea
export default function TestArea() {
  const [sent, setSent] = useState<string>("Nothing sent yet...");
  const [received, setReceived] = useState<string>("Nothing received yet...");
  const [outputs, setOutputs] = useState<string[]>([]);

  // Función para hacer solicitudes POST a la URL del microservicio
  const makeRequest = async (url: string, payload: object, service: string) => {
    try {
      const headers = new Headers();
      headers.append("Content-Type", "application/json");

      const res = await fetch(url, {
        method: "POST",
        body: JSON.stringify(payload),
        headers: headers,
      });

      if (!res.ok) {
        throw new Error(`HTTP error! Status: ${res.status}`);
      }

      const data = await res.json();

      setSent(JSON.stringify(payload, undefined, 4));
      setReceived(JSON.stringify(data, undefined, 4));
      setOutputs([`Response from ${service}`, data.message, new Date().toString()]);
    } catch (error) {
      console.error("Request failed", error);
      setOutputs([`Error with ${service}`, error.message, new Date().toString()]);
    }
  };

  // Funciones específicas para cada microservicio
  const getBroker = async () => {
    const payload = { content: "Empty request" };
    const apiUrl = import.meta.env.VITE_API_URL_LOCAL || "http://localhost:8080"; // Usa la URL correcta del .env
    await makeRequest(apiUrl, payload, "Broker");
  };

  const getAuthentication = async () => {
    const payload = { email: "admin@example.com", password: "password123" };
    const apiUrl = import.meta.env.VITE_API_URL_LOCAL || "http://localhost:8080"; // Usa la URL correcta
    await makeRequest(apiUrl + "/authentication", payload, "Authenticator");
  };

  const getRabbitMQAuthentication = async () => {
    const payload = { email: "admin@example.com", password: "password123" };
    const apiUrl = import.meta.env.VITE_API_URL_LOCAL || "http://localhost:8080"; // Usa la URL correcta
    await makeRequest(apiUrl + "/rabbitmq-authentication", payload, "RabbitMQ Authenticator");
  };

  const getLogger = async () => {
    const payload = { name: "activity", data: "some kind of grpc data" };
    const apiUrl = import.meta.env.VITE_API_URL_LOCAL || "http://localhost:8080"; // Usa la URL correcta
    await makeRequest(apiUrl + "/grpc-logger", payload, "Logger");
  };

  return (
    <div className="container">
      <div className="row">
        <div className="col">
          <h1 className="mt-5 text-light">Go Distributed System</h1>
          <hr />
          {/* Botones para cada servicio */}
          <button className="btn btn-outline-secondary text-light" onClick={getBroker}>Test Broker</button>
          <button className="btn btn-outline-secondary text-light" onClick={getAuthentication}>Test Auth</button>
          <button className="btn btn-outline-secondary text-light" onClick={getRabbitMQAuthentication}>Test RabbitMQ Auth</button>
          <button className="btn btn-outline-secondary text-light" onClick={getLogger}>Test gRPC Logger</button>

          <div className="mt-5" style={{ outline: "1px solid silver", padding: "2em" }}>
            {outputs.length === 0 ? (
              <span className="text-muted">Output shows here...</span>
            ) : (
              <>
                <strong className="text-success">Started</strong><br />
                <i className="text-light">Sending request...</i><br />
                <strong className="text-light">{outputs[0]}: </strong>
                <span className="text-light">{outputs[1]}</span><br />
                <strong className="text-danger">Ended: </strong>
                <span className="text-light">{outputs[2]}</span>
              </>
            )}
          </div>
        </div>
      </div>

      <div className="row">
        <div className="col">
          <h4 className="mt-5 text-light">Sent</h4>
          <div className="mt-1" style={{ outline: "1px solid silver", padding: "2em" }}>
            <pre>
              <span className="text-light" style={{ fontWeight: "bold" }}>
                {sent}
              </span>
            </pre>
          </div>
        </div>
        <div className="col">
          <h4 className="mt-5 text-light">Received</h4>
          <div className="mt-1" style={{ outline: "1px solid silver", padding: "2em" }}>
            <pre>
              <span className="text-light" style={{ fontWeight: "bold" }}>
                {received}
              </span>
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}
