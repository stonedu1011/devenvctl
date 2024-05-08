# Note: CockroachDB v21.1.x cannot be built with go 1.17+
FROM golang:1.16 as builder
ENV DEBIAN_FRONTEND=noninteractive
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
    && echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list \
    && apt-get -y update && apt-get -y upgrade \
    && apt-get -y install libncurses-dev cmake autoconf nodejs yarn bison ccache python \
    && go install github.com/bazelbuild/bazelisk@latest \
    && ln -s $GOPATH/bin/bazelisk $GOPATH/bin/bazel

FROM builder as build
ARG COCKROACHDB_VERSION
ENV USE_BAZEL_VERSION=cockroachdb/latest
RUN wget -qO- https://binaries.cockroachdb.com/cockroach-v$COCKROACHDB_VERSION.src.tgz | tar xvz --strip-components=1
RUN make install

FROM ubuntu:latest
ENV PATH=/cockroach:$PATH
RUN apt-get -y clean && apt-get -y update && apt-get -y upgrade \
    && mkdir -p /usr/local/lib/cockroach /cockroach

COPY --from=build /usr/local/bin/cockroach /cockroach/cockroach
COPY cockroach.sh /cockroach
RUN chmod 755 /cockroach/cockroach.sh

WORKDIR /cockroach/
EXPOSE 26257 26258
ENTRYPOINT ["/cockroach/cockroach.sh"]