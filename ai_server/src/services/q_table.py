
from src.configs.configs import *
import os.path
import numpy as np
import random
from tempfile import TemporaryFile
import pickle 


class QTable:
    def __init__(self, nodesCount, nodesDispersion, epsilon):
        self.nodesCount = nodesCount
        self.nodesDispersion = nodesDispersion
        self.epsilon = epsilon - 1
        self.QTable = self.loadOrCreateTable()
        print(self.QTable)

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

    def updateQValueForTheActionInPreviousState(self, previousState, previousActionTaken):
        """
            Updates the Q-Value for the action taken in the previous state
        """
        
        # use formula to calculate Q-value
        q_value = 12

        self.QTable[previousState][previousActionTaken] = q_value


    def chooseActionForState(self, nodesDispersion):
        """
            Chooses an action for the given state
        """

        # Generate a string from nodes dispersion in classes 
        # {'active': 0, 'idle': 1, 'off': 2} => "2-1-0"
        state = ""
        for className, count in reversed(nodesDispersion.items()):
            state += str(count) + "-"
        state = state[:-1]

        actionsListAlreadyExists = False
        if state in self.QTable.keys():
            # If the given state is in the Q-Table
            actions = self.QTable[state]
            actionsListAlreadyExists = True
        else:
            # If not, generate it
            actions = self.GenerateActionListForState()
            self.QTable[state] = actions

        if random.randrange(1,1000) > self.epsilon and actionsListAlreadyExists:
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
            # Take random action with probability epsilon
            selectedAction = random.choice(list(actions.keys()))

        return state, selectedAction

    def GenerateActionListForState(self):
        """
            Generates a dict of actions
        """
        action = {}
        print(self.nodesDispersion['idle'], self.nodesDispersion['off'])
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