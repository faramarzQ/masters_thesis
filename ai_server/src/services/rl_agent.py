
import numpy as np
import os.path
from src.services.q_table import QTable

def Test(body):
    q_table = QTable(
        body['NodesCount'],
        body['NodesDispersion'],
        body['PreviousEpsilonValue'],
        body['SuccessfulRequests'],
        body['TotalRequests']
    )

    # Pass failed requests, success requests, cpu and memory utilization
    if body["ExecutedPreviously"] == True:
        q_table.updateQValueForTheActionInPreviousState(body["PreviousState"], body["PreviousActionTaken"])

    state, action = q_table.chooseActionForState(body['NodesDispersion'])

    q_table.persistQTable()

    return {"state": state, "action": action, "epsilon": q_table.getEpsilon()}