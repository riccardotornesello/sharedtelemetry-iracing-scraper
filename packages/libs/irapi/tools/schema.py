import os
import json

from genson import SchemaBuilder


def gen_schemas():
    for file in os.listdir("output/responses"):
        builder = SchemaBuilder()
        builder.add_object(json.load(open(f"output/responses/{file}")))
        schema = builder.to_schema()

        with open(f"output/schemas/{file}", "w") as f:
            f.write(json.dumps(schema, indent=2))
