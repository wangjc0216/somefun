# Stage 2
FROM alpine
RUN mkdir /gpt-api
WORKDIR /gpt-api
COPY gpt-api ./
COPY config.json ./
CMD ["./gpt-api"]

