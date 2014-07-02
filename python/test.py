import sys
sys.path.append("./")

from mtconfig import mtconfigserver
import time

#--------------------------------------
mtcs = mtconfigserver.mtconfigserver()
while 1:
    print mtcs.get("inftest", "common", "test")
    time.sleep(2)

