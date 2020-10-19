import time

def handler(event, context):
    print(event, context)
    print("handling ...")
    time.sleep(2)
    print("done")

