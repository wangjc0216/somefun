cd ../../cmd/gpt/
GOOS=linux GOARCH=amd64 go build -o gpt-api
mv gpt-api ../../build/gpt/
cp config.json ../../build/gpt/
cd ../../build/gpt/
docker build -t gpt-api ./
rm config.json gpt-api

