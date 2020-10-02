import random

from cloudfun_pb2 import ServerMetricHistory, RepairServer

from statefun import StatefulFunctions
from statefun import RequestReplyHandler
from statefun import kafka_egress_record

functions = StatefulFunctions()


@functions.bind("model/predictive-maintenance")
def score(context, msg: ServerMetricHistory):
    print(msg)
    if not repair(msg):
        return

    request = RepairServer()
    request.server_id = context.address.identity
    message = kafka_egress_record(topic="repair", key=request.server_id, value=request)
    print("fixing " + request.server_id)
    context.pack_and_send_egress("io/repairs", message)


def repair(history):
    status = is_failing(history)
    correct = random.uniform(0, 1) > 0.6
    if correct:
        return status
    else:
        return not status


def is_failing(history):
    for report in history.reports:
        if report.cpu_percent_nanoseconds_idle > 75:
            return True
    return False

handler = RequestReplyHandler(functions)

from flask import request
from flask import make_response
from flask import Flask

app = Flask(__name__)

@app.route('/functions', methods=['POST'])
def handle():
    response_data = handler(request.data)
    response = make_response(response_data)
    response.headers.set('Content-Type', 'application/octet-stream')
    return response


if __name__ == "__main__":
    app.run()
