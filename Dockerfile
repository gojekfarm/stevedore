FROM golang:1.14 as builder
WORKDIR /go/src/github.com/gojekfarm/stevedore
COPY ./ ./
RUN make compile

FROM alpine:3.12.4
ENV HELM_VERSION="v3.3.1"
RUN wget -q https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz -O - | tar -xzO linux-amd64/helm > /usr/local/bin/helm
RUN chmod +x /usr/local/bin/helm
RUN helm repo add stable https://charts.helm.sh/stable
RUN helm repo update

ARG UID=stevedore
ARG GID=stevedore
RUN addgroup -S $GID && adduser -S $UID -G $GID
COPY --from=builder /go/src/github.com/gojekfarm/stevedore/out/stevedore /usr/local/bin/
USER $UID
WORKDIR /workdir
ENTRYPOINT [ "stevedore" ]
CMD ["version"]
