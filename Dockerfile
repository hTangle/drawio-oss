FROM alpine:3.15
RUN mkdir -p /home/agent
WORKDIR /home/agent
COPY ./output/editor/ /home/agent
RUN chmod +x /home/agent/editor-web
CMD ["./editor-web"]

