FROM swaggerapi/swagger-ui:v5.0.0

COPY docs/openapi.yaml /docs/openapi.yaml

ENV SWAGGER_JSON=/docs/openapi.yaml