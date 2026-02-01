import argparse
import logging
from typing import Final

from requests import Session


logger = logging.getLogger(__name__)
logging.basicConfig(
    format="[%(asctime)s][%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
    level=logging.DEBUG,
)

SERVER_URL: Final[str] = "http://localhost:8000"


def get_incident_by_id(incident_id: str) -> None:
    """Retrieve a single incident by its ID from the server."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        response = session.get(f"{SERVER_URL}/api/v0/incidents/{incident_id}")
        assert response.status_code == 200, f"Expected 200 OK, got {response.status_code}"
        incident_data = response.json()
        logger.info(f"Successfully retrieved incident with ID {incident_id}")
        reporter = incident_data.get("reporter", "Unknown")
        print(f"Incident ID {incident_id}: Reported by {reporter}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Get a single incident by its ID from the server.")
    parser.add_argument("incident", type=str, help="The ID of the incident to retrieve.")
    args = parser.parse_args()
    get_incident_by_id(args.incident)
