import pika

credentials = pika.PlainCredentials('rabbitmq', 'rabbithole')
connectParam = pika.ConnectionParameters(host = 'localhost', port = '5672', virtual_host = '/', credentials = credentials)
connection = pika.BlockingConnection(connectParam)
channel = connection.channel()

channel.exchange_declare(exchange='auth_bus',
                         exchange_type='direct')

channel.queue_bind(exchange='auth_bus',
                   queue='credstoreIn',
                   routing_key='auth_req')

message = '{"jobid": "1234", "username": "todd", "password": "1234abcd", "handwriting": "1a2b3c4d", "action": "create"}'

channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=message)