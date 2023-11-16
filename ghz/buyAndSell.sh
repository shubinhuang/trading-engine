ghz -n 500000 --rps  1000 \
    --concurrency-schedule=step --concurrency-start=5 \
    --concurrency-step=5 --concurrency-end=500 \
    --concurrency-step-duration=10s \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D buyAndSell.json \
    -o report/buy_sell_n500k_rps1000_500.html --format html \
    localhost:5432

# 买卖