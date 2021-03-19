FROM golang:1.15 as builder
WORKDIR /go/src/github.com/gojekfarm/stevedore
COPY ./ ./
RUN make compile

FROM alpine:3.12.4
ENV HELM_VERSION="v3.4.0"
RUN wget -q https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz -O - | tar -xzO linux-amd64/helm > /usr/local/bin/helm
RUN chmod +x /usr/local/bin/helm
RUN helm repo add stable https://charts.helm.sh/stable
RUN helm repo update

COPY --from=builder /go/src/github.com/gojekfarm/stevedore/out/stevedore /usr/local/bin/
WORKDIR /workdir
ENTRYPOINT [ "stevedore" ]
CMD ["version"]
