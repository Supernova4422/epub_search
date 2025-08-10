FROM golang:1.24-alpine

ARG binary_filename="main"
ARG project_path="/project"

COPY . ${project_path}
WORKDIR ${project_path}

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o ${project_path}/${binary_filename} .

ENV binary_filename=${project_path}/${binary_filename}
ENV OTEL_EXPORTER_OTLP_INSECURE="true"
ENV OTEL_RESOURCE_ATTRIBUTES="service.name=epubSearch,service.namespace=epubSearch,deployment.environment=prod"
ENV OTEL_EXPORTER_OTLP_ENDPOINT=http://pi-server:4317
ENV OTEL_EXPORTER_OTLP_PROTOCOL=grpc

EXPOSE ${port}

CMD $binary_filename "-p" "${port}" "-f" "${config_path}"
