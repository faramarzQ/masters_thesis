
from src.configs.configs import *
import os.path
import numpy as np
import random
from tempfile import TemporaryFile
import pickle 
import math
import logging

class QTable:
    def __init__(self, nodesCount, nodesDispersion, previousState, previousAction, energyConsumptionWeight, successRateWeight, alfa, gamma):
        self.nodesCount = nodesCount
        self.nodesDispersion = nodesDispersion
        self.previousState = previousState
        self.previousAction = previousAction
        self.energyConsumptionWeight = energyConsumptionWeight
        self.successRateWeight = successRateWeight
        self.alfa = alfa
        self.gamma = gamma
        self.QTable = self.loadOrCreateTable()

        # Generate a string from nodes dispersion in classes 
        # {'active': 0, 'idle': 1, 'off': 2} => "2-1-0"
        self.currentState = ""
        for className, count in reversed(nodesDispersion.items()):
            self.currentState += str(count) + "-"
        self.currentState = self.currentState[:-1]
        logging.info("Current state: %s", self.currentState)

        if self.currentState not in self.QTable.keys():
            # If the given self.currentState is not in the Q-Table, generate it
            actions = self.GenerateActionListForState()
            self.QTable[self.currentState] = actions

        logging.info("Table: %s", self.QTable)

    def loadOrCreateTable(self):
        """
            Loads the Q-Table from pkl file if it exists
            If not, create a new one
        """
        exists = os.path.exists(QTablePath)
        if exists:
            with open(QTablePath, 'rb') as f:
                QTable = pickle.load(f)
        else:
            QTable = {}

        return QTable

    def persistQTable(self):
        """
            Persists the Q-Table as pkl file
        """
        with open(QTablePath, 'wb') as f:
            pickle.dump(self.QTable, f)

    def updateQValueForTheActionInPreviousRun(self, SuccessRequestRate, clusterEnergyConsumption):
        """
            Updates the Q-Value for the action taken in the previous state
        """
        # Calculate reward
        reward = ( self.successRateWeight * SuccessRequestRate ) - ( self.energyConsumptionWeight * clusterEnergyConsumption )
        logging.info("Reward: %s", reward)

        # self.QTable[self.state][self.actionTaken] += self.alpha * (reward + self.gamma * np.max(Q[self.state]) - Q[self.previousState][previousActionTaken])
        bestQValue = 0
        for action in self.QTable[self.currentState]:
            if bestQValue <= self.QTable[self.currentState][action]:
                bestQValue = self.QTable[self.currentState][action]

        self.QTable[self.previousState][self.previousAction] = (1 - self.alfa) * (self.QTable[self.previousState][self.previousAction]) + self.alfa * (reward + self.gamma * bestQValue)

        logging.info("Table %s", self.QTable)

    def chooseActionForState(self, nodesDispersion):
        """
            Chooses an action for the given state
        """
        actions = self.QTable[self.currentState]

        if random.randrange(1,100) > self.epsilon:
            # Take greedy action most of the time
            selectedAction = 0
            bestValue = 0
            for key, value in self.QTable[self.currentState].items():
                # if (value > bestValue or # To choose the action with highest value
                if value > bestValue: # To choose the action with highest value
                    # abs(key) < abs(selectedAction)): # If values are the same or less, choose one with less transition number
                    selectedAction = key
                    bestValue = value
            logging.info("Action: greedy | Selected action: %s, QValue: %s", selectedAction, bestValue)

        else:
            # Take random action with probability epsilon
            selectedAction = random.choice(list(self.QTable[self.currentState].keys()))
            logging.info("Action: random | Selected action:  %s", selectedAction)

        # TODO: decrease epsilon value on each run

        return self.currentState, selectedAction

    def GenerateActionListForState(self):
        """
            Generates a dict of actions
        """
        action = {}
        index = -self.nodesCount
        for i in range((self.nodesCount*2)+1):
            # Pass inexecutable actions
            if ((index < 0 and -self.nodesDispersion['idle'] > index) or
                (index > 0 and self.nodesDispersion['off'] < index)):
                index = index +1
                continue

            action[index] = 0
            index = index +1

        return action

    def generateNewEpsilonValue(self, step, minimumEpsilonValue, maximumEpsilonValue, EDR):
        self.epsilon = minimumEpsilonValue + (maximumEpsilonValue - minimumEpsilonValue) * math.pow(math.e, (-EDR * step))
        logging.info("New epsilon:  %s", self.epsilon)
