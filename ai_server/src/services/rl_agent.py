
import numpy as np
import os.path
from src.services.q_table import QTable

def Test(body):
    q_table = QTable(
        body['NodesCount'],
        body['NodesDispersion'],
        body["PreviousState"],
        body["State"],
        body["PreviousActionTaken"],
        body["ActionTaken"],
        body['EpsilonValue'],
        body['EnergyConsumptionWeight'],
        body['SuccessRateWeight']
    )

    # if executed before, update the q value for the state-action pair
    if body["PreviousState"] != "":
        q_table.updateQValueForTheActionInPreviousRun(body["SuccessRequestRate"], body['ClusterEnergyConsumption'])

    currentState, action = q_table.chooseActionForState(body['NodesDispersion'])

    q_table.persistQTable()

    return {"state": currentState, "action": action, "epsilon": q_table.getEpsilon()}