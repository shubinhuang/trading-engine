ghz -n 100000 -c 30 \
    --load-schedule=step \
    --load-start=50 --load-end=1000 \
    --load-step=10 --load-step-duration=5s \
    --insecure --proto ../protos/processOrder/processOrder.proto \
    --call processOrder.OrderService.CreateOrder \
    -D buyAndSell.json \
    -o report/buy_sell_n500k_c100_rps.html --format html \
    -m '{"symbol": "s1"}' \
    localhost:32799

# 固定并发数为100，rps从50开始，每5s递增10，最高到1000