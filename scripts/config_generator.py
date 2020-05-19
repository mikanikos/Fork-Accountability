import argparse

# parsing argument
parser = argparse.ArgumentParser()
parser.add_argument('-N', type=int, help='Number of validators')
parser.add_argument('-M', type=int, help='Number of rounds')

args = parser.parse_args()

# random decision values
value1 = 10
value2 = 20

m = int(args.M)
startRound = 3
endRound = startRound + m-1

# base port
port = 8080

# parameters
n = int(args.N)
faulty = int((n-1)/3)

# default config block
prevote = """
          - type: PREVOTE
            sender: {0}
            round: 3
            value:
              data: {1}
"""

precommit = """
          - type: PRECOMMIT
            sender: {0}
            round: 3
            value:
              data: {1}
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
      # round
      3:
"""

# writing monitor config file 
with open("config.yaml", 'w') as f:
  f.write(monitor_block.format(startRound, endRound)

  for j in range(1,n+1):
    f.write("  - 127.0.0.1:{0}\n".format(port+j))


# writing validator configs

for r in range(startRound, endRound+1):



    for i in range(1,n+1):

        filename = "config_{}.yaml".format(i)
        with open(filename, 'a') as f:
            f.write(validator_block.format(i, port+i))

            # wiriting sent/received prevotes/precommits for simple configuration
            f.write("""
            received_prevote:""")

            for j in range(1,n+1):
                f.write(prevote.format(j, value1))

            f.write("""
            sent_prevote:""")

            f.write(prevote.format(i, value1))

            f.write("""
            received_precommit:""")

            for j in range(1,faulty+2):
                f.write(precommit.format(j, value1))
                f.write(precommit.format(j, value2))

            thresh = int((n - (faulty+1))/2)
            for j in range(faulty+2, faulty+thresh+2):
                f.write(precommit.format(j, value1))

            for j in range(faulty+2+thresh, n+1):
                f.write(precommit.format(j, value2))

            f.write("""
            sent_precommit:""")

            f.write(precommit.format(i, value1))