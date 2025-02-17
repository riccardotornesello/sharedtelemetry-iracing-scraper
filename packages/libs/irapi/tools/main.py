import json
import os
import shutil

from session import authentiate_iracing
from endpoints import generate_endpoints
from responses import get_responses
from schema import gen_schemas
from structs import convert_schema_to_struct


def main():
    # Create the output folder
    shutil.rmtree("output", ignore_errors=True)
    os.makedirs("output")
    os.makedirs("output/schemas")
    os.makedirs("output/responses")
    os.makedirs("output/sdk")

    # Authenticate
    session = authentiate_iracing("asdasd", "asd")  # TODO: add credentials

    # Get the API definition
    res = session.get("https://members-ng.iracing.com/data/doc")
    if res.status_code != 200:
        raise Exception(f"Failed to fetch API definition: {res.text}")

    api_definition = res.json()
    with open("output/api_definition.json", "w") as f:
        json.dump(api_definition, f, indent=2)

    # Format the API definition
    endpoints = generate_endpoints(api_definition)
    with open("output/endpoints.json", "w") as f:
        json.dump(endpoints, f, indent=2)

    # Get some sample responses
    # TODO: get responses that require parameters and csv
    get_responses(session, endpoints)

    # Generate the schemas
    # TODO: handle maps
    gen_schemas()

    # Remove the schemas that require maps
    for schema in ["series__seasons", "car__assets"]:
        try:
            os.remove(f"output/schemas/{schema}.json")
        except FileNotFoundError:
            pass

    # Generate the structs
    for category, category_endpoints in endpoints.items():
        camel_category = "".join([x.title() for x in category.split("_")])

        os.makedirs(f"output/sdk/{category}")

        with open(f"output/sdk/{category}/main.go", "w") as f:
            f.write(
                f"""
                    package {category}

                    type {camel_category}Client struct {{
                        client *IRacingApiClient
                    }}
                """
            )

        for endpoint, data in category_endpoints.items():
            camel_endpoint = "".join([x.title() for x in endpoint.split("_")])

            schema_file = f"output/schemas/{category}__{endpoint}.json"
            endpoint_file = f"output/sdk/{category}/{endpoint}.go"

            struct_name = camel_category + camel_endpoint

            if not os.path.exists(schema_file):
                print(f"Schema file not found: {schema_file}")
                continue

            with open(schema_file) as f:
                data = (
                    f"package {category}\n\ntype {struct_name} "
                    + convert_schema_to_struct(json.load(f))
                )
                with open(endpoint_file, "w") as f2:
                    f2.write(data)

    # Format the output
    os.system("go fmt ./output/sdk/...")


main()
