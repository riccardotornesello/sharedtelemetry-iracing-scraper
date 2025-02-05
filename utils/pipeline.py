# NOTE: NOT YET USED

import os
import subprocess


wd = os.getcwd()


def find_go_mod():
    go_mods = {}

    for x in ["apps", "libs"]:
        root = os.path.join("packages", x)
        dirs = os.listdir(root)
        for d in dirs:
            dir_path = os.path.join(root, d)
            if os.path.exists(os.path.join(dir_path, "go.mod")):
                with open(os.path.join(dir_path, "go.mod")) as f:
                    module_name = f.readline().split(" ")[1].strip()
                    go_mods[module_name] = {"path": dir_path}

    return go_mods


def get_mod_dependencies():
    os.chdir("packages/apps/api")
    output = subprocess.check_output(["go", "mod", "graph"])
    lines = output.split(b"\n")
    words = []
    for line in lines:
        for word in line.split():
            words.append(word)
    os.chdir(wd)


go_mods = find_go_mod()
print(go_mods)

get_mod_dependencies()
