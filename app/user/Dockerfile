FROM docker.sharkgulf.cn/public/alpine-ca:3.10

RUN mkdir -p /app /app/log

COPY ./user ./app.yml /app/

WORKDIR /app

VOLUME [ "/app/log" ]

RUN chmod +x ./user

EXPOSE 8080

ENTRYPOINT ["./user"]
