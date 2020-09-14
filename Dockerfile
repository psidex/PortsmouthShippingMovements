FROM golang:latest AS builder
WORKDIR /psmbuild
COPY . .
# Write short hash of current commit to version file
RUN head -c 7 ./.git/refs/heads/master > ./version
RUN CGO_ENABLED=0 GOOS=linux go build -o ./psmserver ./cmd/server/main.go

FROM alpine:latest
WORKDIR /psm
COPY static static
COPY --from=builder /psmbuild/psmserver .
COPY --from=builder /psmbuild/version ./static
ENV INSIDE_DOCKER "True"
EXPOSE 8080/tcp
CMD ["./psmserver"]

# docker run -d --name psmserver \
#     --network proxynet \
#     -v /home/user/psmdockervolume:/data \
#     psidex/portsmouthshippingmovements
