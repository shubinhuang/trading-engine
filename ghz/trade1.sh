ghz -c 1000 -n 1000000 --connections 5 \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D trade1.json  \
    -m '{"symbol": "s1"}' \
    localhost:32774

# 买卖
