ghz -c 10 -n 1000000 --connections 5 \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D trade0.json  \
    -m '{"symbol": "s0"}' \
    localhost:32774

# 买卖
