import pika
import time

credentials = pika.PlainCredentials('rabbitmq', 'rabbithole')
connectParam = pika.ConnectionParameters(host = 'localhost', port = '5672', virtual_host = '/', credentials = credentials)
connection = pika.BlockingConnection(connectParam)
channel = connection.channel()

channel.exchange_declare(exchange='auth_bus',
                         exchange_type='direct')

channel.queue_bind(exchange='auth_bus',
                   queue='credstoreIn',
                   routing_key='auth_req')

create_message = '{"jobid": "1234", "username": "todd", "password": "1234abcd", "handwriting": "1a2b3c4d", "race": "1234", "action": "create"}'
auth_message = '{"jobid": "1235", "username": "todd", "password": "1234abcd", "action": "authenticate"}'
failed_auth_message = '{"jobid": "1236", "username": "todd", "password": "1a2b3c4d", "action": "authenticate"}'
collect_message = '{"jobid": "1237", "username": "todd", "handwriting": "1a2b3c4d", "action": "collect"}'
failed_collect_message = '{"jobid": "1238", "username": "ddot", "handwriting": "1a2b3c4d", "action": "collect"}'


# this is for creating a record
channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=create_message)

time.sleep(3)

# this is for a success authenticate
channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=auth_message)

time.sleep(3)

# this is for a fail authenticate
channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=failed_auth_message)

time.sleep(3)

# this is for collect handwriting
channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=collect_message)

time.sleep(3)

# this is for a failure collect handwriting
channel.basic_publish(exchange='auth_bus',
                      routing_key='auth_req',
                      body=failed_collect_message)

channel.close()
connection.close()