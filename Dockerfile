FROM golang AS build-stage
RUN go install github.com/mannemsolutions/pgtester/cmd/pgtester@main

FROM alpine AS export-stage
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=build-stage /go/bin/pgtester /
COPY testdata /
CMD /pgtester -v
