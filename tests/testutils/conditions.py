import os

import pytest


def integration_test(testcase_id: str) -> pytest.MarkDecorator:
    """Mark a test as an integration test."""
    is_integration = os.getenv("RUN_INTEGRATION") is None
    reason = f"[{testcase_id}] Needs live instance of server"
    return pytest.mark.skipif(condition=is_integration, reason=reason)
