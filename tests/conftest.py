from typing import Iterator

import pytest
from fastapi.testclient import TestClient

from postfacta.service import app


@pytest.fixture
def client() -> Iterator[TestClient]:
    """Instantiate the test client for testing."""
    yield TestClient(app)
