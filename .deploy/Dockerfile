FROM dhub.kubesre.xyz/library/alpine:latest

RUN apk update && apk add tzdata
COPY go_app /runtime/go_app
WORKDIR /runtime
ENTRYPOINT /runtime/go_app