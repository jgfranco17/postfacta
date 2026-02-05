import argparse
import os
from typing import Final

from faker import Faker
from requests import Session

DEFAULT_SERVER_URL: Final[str] = "http://localhost:8000"
MOCK_DATA_COUNT: Final[int] = 10


def _get_url_from_env() -> str:
    """Retrieve the server URL from environment variables or use the default."""
    env_var = "POSTFACTA_SERVER_URL"
    url_from_env = os.getenv(env_var)
    if url_from_env is None:
        print(f"Environment variable {env_var} not set, using default URL: {DEFAULT_SERVER_URL}")
        return DEFAULT_SERVER_URL

    print(f"Using server URL from environment: {url_from_env}")
    return url_from_env


def _create_single_mock_incident(session: Session, base_api_url: str, count: int) -> str:
    """Create a mock incident via the API and return its ID."""
    fake = Faker()
    request_payload = {
        "title": f"Integration Test Incident {count}",
        "description": f"This is a test incident #{count} created during integration testing.",
        "reporter": fake.name(),
        "severity": "LOW",
        "owner": fake.name(),
    }
    api_url = f"{base_api_url}/api/v0/incidents/start"
    response = session.post(api_url, json=request_payload)
    print(f"Created mock incident ({response.status_code}): {response.text}")
    return response.json()["id"]


def create_mock_incidents(base_api_url: str) -> None:
    """Populate the mock server with test data."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        for idx in range(MOCK_DATA_COUNT):
            _create_single_mock_incident(session, base_api_url, idx + 1)
        print(f"Successfully created {MOCK_DATA_COUNT} mock incidents")


def get_all_incidents_from_server(base_api_url: str) -> None:
    """Retrieve the list of incidents from the server."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        response = session.get(f"{base_api_url}/api/v0/incidents")
        assert response.status_code == 200
        incidents = response.json().get("incidents")
        print(f"Successfully retrieved {len(incidents)} incidents")
        for incident_id, incident_data in incidents.items():
            reporter = incident_data.get("reporter", "Unknown")
            print(f"Incident ID {incident_id}: Reported by {reporter}")


def get_incident_by_id(base_api_url: str, incident_id: str) -> None:
    """Retrieve a single incident by its ID from the server."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        response = session.get(f"{base_api_url}/api/v0/incidents/{incident_id}")
        assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}"
        incident_data = response.json()
        print(f"Successfully retrieved incident with ID {incident_id}")
        reporter = incident_data.get("reporter", "Unknown")
        print(f"Incident ID {incident_id}: Reported by {reporter}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Get a single incident by its ID from the server.")
    parser.add_argument("action", type=str, help="Action to perform against target service.")
    parser.add_argument("--incident", "-i", type=str, help="The ID of the incident to interact with.")
    args = parser.parse_args()

    base_api_url = _get_url_from_env()

    match args.action:
        case "fill":
            create_mock_incidents(base_api_url)
        case "get":
            if not args.incident:
                get_all_incidents_from_server(base_api_url)
            else:
                get_incident_by_id(base_api_url, args.incident)
        case _:
            raise ValueError(f"Unknown action specified: {args.action}")
