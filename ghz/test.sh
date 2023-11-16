ghz -c 1 -n 1 \
    --insecure --proto ../protos/openTrade/openTrade.proto \
    --call openTrade.OpenService.OpenTrade \
    -d '{"symbol": "{{.RequestNumber}}","price": "{{randomInt 100 200}}"}' \
    -m '{"symbol": "s1"}' \
    localhost:6543

# 买卖