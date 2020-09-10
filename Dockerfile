FROM golang:latest AS builder
WORKDIR /psmbuild
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./psmserver ./cmd/server/main.go

FROM alpine:latest
WORKDIR /psm
COPY --from=builder /psmbuild/psmserver .
COPY static static
ENV INSIDE_DOCKER "True"
EXPOSE 8080/tcp
CMD ["./psmserver"]

# e.g.
# docker build -t latest .
# docker run -d --name psmserver \
#     -p 80:8080
#     -v /home/user/psmdockervolume:/data
#     latest
