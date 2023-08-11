
import numpy as np
import os.path
from src.services.q_table import QTable

def runReinforcementLearning(body):
    q_table = QTable(
        body['NodesCount'],
        body['NodesDispersion'],
        body["PreviousState"],
        body["PreviousAction"],
        body['EnergyConsumptionWeight'],
        body['SuccessRateWeight'],
        body['Alfa'],
        body['Gamma'],
    )

    q_table.generateNewEpsilonValue(body['Step'], body['MinimumEpsilonValue'], body['MaximumEpsilonValue'], body['EDR'])

    # if executed before, update the q value for the previous state-action pair
    if body["Step"] != 1:
        q_table.updateQValueForTheActionInPreviousRun(body["SuccessRequestRate"], body['ClusterEnergyConsumption'])

    currentState, action = q_table.chooseActionForState(body['NodesDispersion'])

    q_table.persistQTable()

    return {"state": currentState, "action": action}