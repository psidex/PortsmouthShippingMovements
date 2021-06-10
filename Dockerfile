FROM golang:latest AS go-builder
WORKDIR /build
COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./psmserver ./cmd/server/main.go

FROM node:14.17 AS frontend-builder
ENV PATH /app/node_modules/.bin:$PATH
WORKDIR /app
COPY frontend/. .
COPY .git/refs/heads/master ./public/version
RUN yarn install
RUN yarn prodbuild

FROM alpine:latest
WORKDIR /app
COPY --from=go-builder /build/psmserver .
COPY --from=frontend-builder /app/build ./frontend_build
EXPOSE 8080/tcp
ENTRYPOINT ["./psmserver"]

# docker run -d --name psmserver \
#     --network proxynet \
#     -v $(pwd)/config.json:/psm/config.json:ro \
#     -v psmdata:(config.json->storage_path, usually /psm/storage) \
#     psidex/portsmouthshippingmovements
