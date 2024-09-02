FROM golang:1.22-alpine

WORKDIR /app

COPY . .

# 必要なビルドツールをインストール
RUN apk add --no-cache build-base

# 依存関係を整理
RUN go mod tidy

RUN go build -o /portal_app

EXPOSE 8080

CMD [ "/portal_app" ]