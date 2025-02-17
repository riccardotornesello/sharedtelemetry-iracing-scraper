def generate_key(key):
    if key.isdigit():
        key = f"Field{key}"

    return "".join(x.capitalize() for x in key.lower().replace("-", "_").split("_"))


def convert_schema_to_struct(schema, json_name=None):
    data = ""

    if "type" not in schema:
        raise Exception(f"Missing type in key {json_name}")

    t = schema["type"]

    if isinstance(t, list):
        if "null" in t:
            t.remove("null")
            data += "*"
        else:
            data += ""

        if len(t) != 1:
            raise Exception(f"Unknown type {t} in key {json_name}")

        t = t[0]

    if t == "array":
        if not "items" in schema:
            print("Warning: array without items")
            data += "[]interface{}"
        else:
            data += "[]"
            data += convert_schema_to_struct(schema["items"])

    elif t == "object":
        if not "properties" in schema:
            print("Warning: object without properties")
            data += "map[string]interface{}"
        else:
            data += "struct {\n"
            for k, v in schema["properties"].items():
                data += f"{generate_key(k)} {convert_schema_to_struct(v,k)}"
            data += "}"

    elif t == "string":
        data += "string"
    elif t == "integer":
        data += "int"
    elif t == "number":
        data += "float64"
    elif t == "boolean":
        data += "bool"
    elif t == "null":
        print("Warning: null type")
        data += "interface{}"
    else:
        raise Exception(f"Unknown type {t} in key {json_name}")

    if json_name:
        data += f' `json:"{json_name}"`\n'

    return data
