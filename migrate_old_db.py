import pymongo
from pymongo import MongoClient, UpdateOne
from pymongo.errors import BulkWriteError
from pprint import pprint



client= MongoClient("hydra")
DBT = client.truth

COL_KILL = DBT.killmails
COL_HASH = DBT.hashes

cnt = 0

updates = []

for hash in COL_HASH.find():
    # COL_KILL.update({'_id':hash['_id']}, {'$set': {'hash': hash['hash']}}, upsert=True)
    update = UpdateOne({'_id': hash['_id']}, {'$set': {'hash': hash['hash']}}, upsert=True)
    updates.append(update)
    cnt += 1

print("Queeud {} Updates!! This may take a while...".format(cnt))

print("Executing updates")
try:
    result = COL_KILL.bulk_write(updates, ordered=False)
except BulkWriteError as bwe:
    pprint(bwe.details)

print("Completed!!")
pprint(result)