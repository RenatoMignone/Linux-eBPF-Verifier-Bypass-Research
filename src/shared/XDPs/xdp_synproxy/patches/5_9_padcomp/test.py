import sys
import subprocess
import time

SUCCESS_MSG = "succeeded"

if len(sys.argv) != 2:
    print(f"Usage: {sys.argv[0]} <address>")
    sys.exit(1)

address = sys.argv[1]
cnt = 0 

while True:
    try:
        result = subprocess.run(
            ["nc", "-zv", address, "80"],
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True
        )

        output = result.stdout.strip()

        if SUCCESS_MSG in output:
            print(f"[+] SUCCESS: {output}")
        else:
            print(f"[-] FAIL: {output}")
            cnt = cnt + 1

        time.sleep(0.2)

    except KeyboardInterrupt:
        print(f"\nInterrupted by user. Exiting. Fails: {cnt}")
        break
    except Exception as e:
        print(f"[!] Error: {e}")
        time.sleep(1)
