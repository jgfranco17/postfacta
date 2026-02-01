import datetime as dt
from enum import StrEnum
from uuid import uuid4

from pydantic import BaseModel


class Severity(StrEnum):
    LOW = "LOW"
    MEDIUM = "MEDIUM"
    HIGH = "HIGH"
    CRITICAL = "CRITICAL"


class Status(StrEnum):
    OPEN = "OPEN"
    IN_PROGRESS = "IN_PROGRESS"
    RESOLVED = "RESOLVED"
    CLOSED = "CLOSED"


class Incident(BaseModel):
    id: str
    title: str
    description: str
    severity: Severity
    status: Status
    reporter: str
    owner: str | None = None
    start_time: str
    end_time: str | None = None


class NewIncidentRequest(BaseModel):
    title: str
    description: str
    severity: Severity
    reporter: str
    owner: str | None = None


def create_new_incident(incident_request: NewIncidentRequest) -> Incident:
    """Create a new incident from a request object.

    Args:
        incident_request (NewIncidentRequest): Request body details

    Returns:
        Incident: The newly created incident instance
    """
    new_id = f"postfacta-inc-{uuid4()}"
    timestamp = dt.datetime.now(dt.timezone.utc).isoformat()
    return Incident(
        id=new_id,
        title=incident_request.title,
        description=incident_request.description,
        severity=incident_request.severity,
        status=Status.OPEN,
        reporter=incident_request.reporter,
        owner=incident_request.owner,
        start_time=timestamp
    )
