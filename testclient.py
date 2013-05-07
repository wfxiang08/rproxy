import redis

r = redis.Redis(port=7000)
r.set("hello","world")