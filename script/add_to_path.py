import os
import platform

root = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
plat = platform.system()
print("Please enter the following in your respective terminal:")
if plat == 'Linux' or plat == 'Darwin':
    print(f"export MKYROOT={root}")
elif plat == 'Windows':
    print(f"set MKYROOT={root}")
