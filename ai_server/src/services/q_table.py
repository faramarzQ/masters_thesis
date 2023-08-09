
from src.configs.configs import *
import os.path
import numpy as np
import random
from tempfile import TemporaryFile
import pickle 

class QTable:
    def __init__(self, nodesCount, nodesDispersion, previousState, state, previousActionTaken, actionTaken, epsilon, energyConsumptionWeight, successRateWeight):
        self.nodesCount = nodesCount
        self.nodesDispersion = nodesDispersion
        self.epsilon = epsilon
        self.previousState = previousState
        self.state = state
        self.previousActionTaken = previousActionTaken
        self.actionTaken = actionTaken
        self.energyConsumptionWeight = energyConsumptionWeight
        self.successRateWeight = successRateWeight
        self.QTable = self.loadOrCreateTable()
        print("\n ----- LOADED TABLE ------", self.QTable)

        # Generate a string from nodes dispersion in classes 
        # {'active': 0, 'idle': 1, 'off': 2} => "2-1-0"
        self.currentState = ""
        for className, count in reversed(nodesDispersion.items()):
            state += str(count) + "-"
        self.currentState = state[:-1]

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
        print("\n ------ REWARD ------ \n", reward)

        # self.QTable[self.state][self.actionTaken] += self.alpha * (reward + self.gamma * np.max(Q[self.state]) - Q[self.previousState][previousActionTaken])
        # TODO: correct belslman equation
        self.QTable[self.state][self.actionTaken] = (1-self.alfa)*(self.QTable[self.state][self.actionTaken]) + self.alfa * (rewards + self.gama )

    def chooseActionForState(self, nodesDispersion):
        """
            Chooses an action for the given state
        """

        actionsListAlreadyExists = False
        if self.currentState in self.QTable.keys():
            # If the given self.currentState is in the Q-Table
            actions = self.QTable[self.currentState]
            actionsListAlreadyExists = True
        else:
            # If not, generate it
            actions = self.GenerateActionListForState()
            self.QTable[self.currentState] = actions

        if random.randrange(1,100) > self.epsilon and actionsListAlreadyExists:
            print("------ GREEDY ACTION ------")
            # Take greedy action most of the time
            selectedAction = 0
            bestValue = 0
            for key, value in actions.items():
                # if (value > bestValue or # To choose the action with highest value
                if value > bestValue: # To choose the action with highest value
                    # abs(key) < abs(selectedAction)): # If values are the same or less, choose one with less transition number
                    selectedAction = key
                    bestValue = value
        else:
            print("------ RANDOM ACTION ------")
            # Take random action with probability epsilon
            selectedAction = random.choice(list(actions.keys()))
            
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

    def getEpsilon(self):
        return self.epsilon