import pytest
from fastapi.testclient import TestClient


@pytest.mark.core
def test_health_endpoint(client: TestClient):
    response = client.get("/healthz")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}


@pytest.mark.core
def test_unsupported_routes(client: TestClient):
    response = client.get("/definitely-invalid")
    assert response.status_code == 404, "Endpoint should not exist in API"
