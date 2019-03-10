import redis

r = redis.Redis(host='localhost', port=6379)
p = r.pubsub()

p.subscribe('result')

while True:
    message = p.get_message()   
    if message:
        print("the message received is {0}".format(message['data']))