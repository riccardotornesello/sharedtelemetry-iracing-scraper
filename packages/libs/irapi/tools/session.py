import requests
import hashlib
import base64


def authentiate_iracing(email, password):
    # Login
    token = hashlib.sha256(f"{password}{email}".encode()).digest()
    base64_token = base64.b64encode(token).decode()

    s = requests.Session()

    response = s.post(
        "https://members-ng.iracing.com/auth",
        json={"email": email, "password": base64_token},
    )

    if response.status_code != 200:
        raise Exception(f"Failed to authenticate: {response.text}")

    return s
