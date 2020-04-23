FROM  golang:1.13.10 AS build-env
ADD . /src
RUN cd /src && make build


FROM scratch
COPY --from=build-env /src/main /
COPY --from=build-env /src/configs/config.yml /configs/config.yml
CMD ["/main"]