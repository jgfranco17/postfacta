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


class Note(BaseModel):
    timestamp: str
    message: str


class Incident(BaseModel):
    id: str
    title: str
    description: str
    severity: Severity
    status: Status
    reporter: str
    start_time: str
    initial_notes: list[Note] = []
    additional_notes: list[Note] = []
    owner: str | None = None
    end_time: str | None = None

    def add_note(self, note: Note) -> None:
        """Add an additional note to the incident.

        Args:
            note (Note): The note to add
        """
        self.additional_notes.append(note)

    def get_notes(self) -> list[Note]:
        """Generate a summary report of the incident.

        Returns:
            list[Note]: Formatted incident report
        """
        return self.initial_notes + self.additional_notes

    def close_incident(self) -> None:
        """Close the incident by updating its status and end time."""
        self.status = Status.CLOSED
        self.end_time = dt.datetime.now(dt.timezone.utc).isoformat()

    def resolve_incident(self) -> None:
        """Resolve the incident by updating its status."""
        self.status = Status.RESOLVED
        self.end_time = dt.datetime.now(dt.timezone.utc).isoformat()


class NewIncidentRequest(BaseModel):
    title: str
    description: str
    severity: Severity
    reporter: str
    initial_notes: list[str] = []
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
    initial_notes = [
        Note(timestamp=timestamp, message=note)
        for note in incident_request.initial_notes
    ]
    return Incident(
        id=new_id,
        title=incident_request.title,
        description=incident_request.description,
        severity=incident_request.severity,
        status=Status.OPEN,
        reporter=incident_request.reporter,
        start_time=timestamp,
        owner=incident_request.owner,
        initial_notes=initial_notes,
        additional_notes=[],
    )
