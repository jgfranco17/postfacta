import logging
from typing import Final

from faker import Faker
from requests import Session


SERVER_URL: Final[str] = "http://localhost:8000"
MOCK_DATA_COUNT: Final[int] = 10


logger = logging.getLogger(__name__)
logging.basicConfig(
    format="[%(asctime)s][%(levelname)s] %(name)s: %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
    level=logging.DEBUG,
)


def create_mock_incident(session: Session, base_api_url: str, count: int) -> str:
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
    logger.debug(f"Created mock incident ({response.status_code}): {response.text}")
    return response.json()["id"]

def populate_mock_data() -> None:
    """Populate the mock server with test data."""
    with Session() as session:
        session.headers.update({"Content-Type": "application/json"})
        for idx in range(MOCK_DATA_COUNT):
            create_mock_incident(session, SERVER_URL, idx + 1)
        logger.info(f"Successfully created {MOCK_DATA_COUNT} mock incidents")


if __name__ == "__main__":
    populate_mock_data()
