import json


def get_responses(session, endpoints):
    # Save responses
    for category, category_endpoints in endpoints.items():
        for endpoint, data in category_endpoints.items():
            if data["format"] == "csv":
                continue

            if any(p["required"] for p in data["parameters"]):
                continue

            url = data["link"]
            response = session.get(url)

            if response.status_code != 200:
                print(f"Failed to fetch {url}: {response.text}")
                continue

            if not data["skip_s3"]:
                s3_url = response.json()["link"]
                response = session.get(s3_url)

            with open(f"output/responses/{category}__{endpoint}.json", "w") as f:
                json.dump(response.json(), f, indent=2)
