
import numpy as np
import os.path
from src.services.q_table import QTable

def Test(body):
    q_table = QTable(body['NodesCount'], body['NodesDispersion'])

    action = q_table.chooseActionForState(body['NodesDispersion'])

    # return selected action + epsilon value
    return ""

# # Number of bandits
# k = 3

# # Our action values
# Q = [0 for _ in range(k)]

# # This is to keep track of the number of times we take each action
# N = [0 for _ in range(k)]

# # Epsilon value for exploration
# eps = 0.1

# # True probability of winning for each bandit
# p_bandits = [0.45, 0.40, 0.80]

# def pull(a):
#     """Pull arm of bandit with index `i` and return 1 if win, 
#     else return 0."""
#     if np.random.rand() < p_bandits[a]:
#         return 1
#     else:
#         return 0

# while False:
#     if np.random.rand() > eps:
#         # Take greedy action most of the time
#         a = np.argmax(Q)
#     else:
#         # Take random action with probability eps
#         a = np.random.randint(0, k)
    
#     # Collect reward
#     reward = pull(a)
    
#     # Incremental average
#     N[a] += 1
#     Q[a] += 1/N[a] * (reward - Q[a])


