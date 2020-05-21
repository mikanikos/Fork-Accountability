import argparse

# parsing argument
parser = argparse.ArgumentParser()
parser.add_argument('-N', type=int, help='Number of validators')
parser.add_argument('-M', type=int, help='Number of rounds')

args = parser.parse_args()

# deafult decision values
value1 = 10
value2 = 20

m = int(args.M)
startRound = 1
endRound = startRound + m-1

# base port
port = 8080

# parameters
n = int(args.N)
faulty = int((n-1)/3)

# default config blocks

prevote = """
          - type: PREVOTE
            sender: {0}
            round: {1}
            value:
              data: {2}
"""

precommit = """
          - type: PRECOMMIT
            sender: {0}
            round: {1}
            value:
              data: {2}
"""

monitor_block = """--- # auto-generated config file
height: 1
firstDecisionRound: {0}
secondDecisionRound: {1}
timeout: 60
# each validator must have a unique address (it is used as id)
validators:
"""

validator_block = """--- # auto-generated config file
id: {0}
address: 127.0.0.1:{1}
messages:
  # height
  1:
    heightvoteset:
"""

round_block = """
      # round
      {0}:
"""

# writing monitor config file
with open("config.yaml", 'w') as f:
    f.write(monitor_block.format(startRound, endRound))

    for j in range(1, n+1):
        f.write("  - 127.0.0.1:{0}\n".format(port+j))


# writing validator configs
for i in range(1, n+1):
    filename = "config_{}.yaml".format(i)
    with open(filename, 'w') as f:
        f.write(validator_block.format(i, port+i))


# configuration: at most 2f faulty processes, the rest is correct
# correct processes split in voting two different values, faulty processes first vote for one and then vote for the other

# writing messages for each round
for r in range(startRound, endRound+1):
    # writing messages for each validator
    for i in range(1, n+1):

        filename = "config_{}.yaml".format(i)
        with open(filename, 'a') as f:
            f.write(round_block.format(r))

            thresh = int((n - (2*faulty))/2)

            # wiriting sent/received prevotes/precommits for simple configuration
            f.write("""
        received_prevote:""")

            if r == startRound:
                for j in range(1, (2*faulty)+thresh+1):
                    f.write(prevote.format(j, r, value1))

            if r == endRound:
                for j in range(1, (2*faulty)+1):
                    f.write(prevote.format(j, r, value2))

                for j in range((2*faulty)+thresh+1, n+1):
                    f.write(prevote.format(j, r, value2))

            f.write("""
        sent_prevote:""")

            f.write("""
        received_precommit:""")

            if r == startRound:
                for j in range(1, (2*faulty)+thresh+1):
                    f.write(precommit.format(j, r, value1))

            if r == endRound:
                for j in range(1, (2*faulty)+1):
                    f.write(precommit.format(j, r, value2))

                for j in range((2*faulty)+thresh+1, n+1):
                    f.write(precommit.format(j, r, value2))

            f.write("""
        sent_precommit:""")
