# syntax=docker/dockerfile:1

FROM ubuntu:latest
WORKDIR /sync_server
COPY /bin/ /sync_server/

RUN cd sync_server
RUN chmod +x sync_server
RUN mkdir photos       

CMD ["sync_server", "-p 3000 -d /photos"]
EXPOSE 3000