FROM golang:1.15-alpine

ARG binary_filename="main"
ARG project_path="/project"

COPY . ${project_path}
WORKDIR ${project_path}

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o ${project_path}/${binary_filename} .

ENV binary_filename=${project_path}/${binary_filename}

EXPOSE ${port}

CMD $binary_filename "-p" "${port}" "-f" "${config_path}"
