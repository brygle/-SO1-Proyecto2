import pika
from pymongo import MongoClient
import redis

connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()

channel.queue_declare(queue='hello')
#---------------configuraciones de mongo-----------------
mongoClient =  MongoClient('35.225.245.55', 27017)
db = mongoClient.proyecto
colleccion = db.casos
#---------------configuraciones de redis-----------------
redisClient = redis.Redis(host = '35.224.140.76', port = 6379)

def callback(ch, method, properties, body):
	cadena = str(body)[2:]
	cadena = cadena[:len(cadena)-1]
	parametros = cadena.split(";")
	nombre = parametros[0].split(",")[1]
	location = parametros[1].split(",")[1]
	age = parametros[2].split(",")[1]
	it = parametros[3].split(",")[1]
	state = parametros[4].split(",")[1]
	print("Nombre: " + nombre )
	print("Location: " + location )
	print("Age: " + age )
	print("InfectedType: " + it )
	print("State: " + state )
	print("")
	colleccion.insert({
		"nombre" : nombre,
		"location" : location,
		"age" : age,
		"infectedType" : it,
		"state" : state
	})
	redisClient.sadd('casos', nombre+','+location+','+str(age)+',' + it + ',' + state)
	#redisClient.sadd('casos', 'caso1')

channel.basic_consume(queue = 'hello', on_message_callback = callback, auto_ack = True )
print('Esperando por mensajes')
channel.start_consuming()