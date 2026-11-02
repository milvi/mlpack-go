FROM ubuntu:26.04
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
      build-essential cmake pkg-config curl ca-certificates patch \
      libarmadillo-dev libensmallen-dev libcereal-dev && \
    rm -rf /var/lib/apt/lists/*
RUN useradd --create-home --shell /bin/bash mlpack
COPY --chown=mlpack:mlpack . /home/mlpack/projet/
USER mlpack
WORKDIR /home/mlpack/projet
ARG JOBS=2
ENV JOBS=${JOBS}
RUN bash installer.sh
RUN bash -c '. "$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/env.sh"; go test -buildvcs=false ./...'
RUN bash -c '. "$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/env.sh"; GOEXPERIMENT=cgocheck2 go test -buildvcs=false ./...'
RUN bash -c '. "$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/env.sh"; GOMAXPROCS=1 OMP_NUM_THREADS=1 OPENBLAS_NUM_THREADS=1 go test -buildvcs=false -tags mlpack_memtest -run TestNative -v .'
CMD ["/home/mlpack/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/bin/kmeans-demo"]
