import pika
from pymongo import MongoClient
import redis
from json import loads, dump
connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()

channel.queue_declare(queue='hello')
#---------------configuraciones de mongo-----------------
#mongoClient =  MongoClient('35.225.245.55', 27017)
mongoClient =  MongoClient('mongodb://35.225.245.55:27017/')
db = mongoClient['proyecto']
colleccion = db.casos
#---------------configuraciones de redis-----------------
redisClient = redis.Redis(host = '34.72.157.182', port = 6379)

def callback(ch, method, properties, body):

	colleccion.insert_one(loads(body.decode()))
	parsed = loads(body.decode())
	string_json=str(parsed)
	string_json = string_json.replace("'", "")
	redisClient.rpush('mylist', string_json)

channel.basic_consume(queue = 'hello', on_message_callback = callback, auto_ack = True )
print('Esperando por mensajes')
channel.start_consuming()

