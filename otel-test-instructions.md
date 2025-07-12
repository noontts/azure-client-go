# OpenTelemetry Tracing Test Instructions

## 1. Start a Trace Backend (Jaeger or OTLP Collector)
- For local testing, you can use Jaeger all-in-one:
  ```sh
  docker run -d --name jaeger \
    -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
    -p 5775:5775/udp \
    -p 6831:6831/udp \
    -p 6832:6832/udp \
    -p 5778:5778 \
    -p 16686:16686 \
    -p 14268:14268 \
    -p 14250:14250 \
    -p 9411:9411 \
    jaegertracing/all-in-one:1.56
  ```
- Or use an OTLP-compatible collector and UI (e.g., Uptrace, Tempo, etc.).

## 2. Set Environment Variable
- Set the OTLP endpoint for your Go app (default is `http://localhost:4318`):
  ```sh
  $env:OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
  ```
  Or set in your `.env` file:
  ```
  OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
  ```

## 3. Run Your Service
- Start your Go service as usual:
  ```sh
  go run main.go
  ```

## 4. Exercise the API
- Use curl, Postman, or your frontend to hit the `/members` endpoints (CRUD).

## 5. View Traces
- Open Jaeger UI at [http://localhost:16686](http://localhost:16686)
- Search for service `azure-client-go` and view traces for HTTP, service, and DB spans.

## 6. Troubleshooting
- Ensure the OTLP endpoint is reachable from your Go app.
- Check logs for any OpenTelemetry errors.

---

You now have full distributed tracing for HTTP, service, and DB layers!
