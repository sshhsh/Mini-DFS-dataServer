FROM alpine
COPY dataServer /dataServer
RUN mkdir /data
CMD /dataServer -addr 172.18.0.10