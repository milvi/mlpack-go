FROM ubuntu:26.04
ARG GO_VERSION=1.27.3
ARG GO_SHA256=65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    build-essential libarmadillo-dev libensmallen-dev libcereal-dev ca-certificates curl python3 && \
    rm -rf /var/lib/apt/lists/*
RUN curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz && \
    echo "${GO_SHA256}  /tmp/go.tgz" | sha256sum -c - && \
    tar -C /usr/local -xzf /tmp/go.tgz && rm /tmp/go.tgz
ENV PATH=/usr/local/go/bin:$PATH GOMAXPROCS=2 GOFLAGS=-p=2
RUN useradd --create-home --shell /bin/bash mlpack
COPY --chown=mlpack:mlpack . /home/mlpack/projet
USER mlpack
WORKDIR /home/mlpack/projet
RUN go version && g++ --version && dpkg-query -W libarmadillo-dev libensmallen-dev libcereal-dev
RUN go mod download && go mod verify && go test -v ./...
RUN GOEXPERIMENT=cgocheck2 go test ./...
RUN go test -tags mlpack_memtest -run TestNative -v .
RUN go build -o /home/mlpack/kmeans ./cmd/kmeans && go run ./examples/kmeans
CMD ["/home/mlpack/kmeans", "-n", "100000", "-k", "1000"]
