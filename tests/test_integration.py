from typing import Iterator

import pytest
from requests import Session

from tests.testutils.conditions import integration_test


@pytest.fixture
def live_client() -> Iterator[Session]:
    session = Session()
    session.headers.update({"Content-Type": "application/json"})
    yield session


@pytest.fixture
def server_url() -> str:
    return "http://localhost:8000"


@integration_test("INT-HEALTH")
@pytest.mark.integration
def test_integration_health(live_client: Session, server_url: str) -> None:
    """Test that the health endpoint is reachable and returns a successful status."""
    response = live_client.get(f"{server_url}/healthz")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


@integration_test("INT-HEALTH")
@pytest.mark.integration
def test_integration_create_incident(live_client: Session, server_url: str) -> None:
    """Test that an incident can be created and that the returned ID has the expected format."""
    request_payload = {
        "title": "Integration Test Incident",
        "description": "This is a test incident created during integration testing.",
        "reporter": "John Doe",
        "severity": "LOW",
        "owner": "Jane Smith",
    }
    response = live_client.post(f"{server_url}/api/v0/incidents/start", json=request_payload)
    assert response.status_code == 201
    incident_id = response.json().get("id")
    assert incident_id.startswith("postfacta-inc-")
