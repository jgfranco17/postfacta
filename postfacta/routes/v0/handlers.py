import logging
import datetime as dt

from fastapi import APIRouter, status
from pydantic import ValidationError

from postfacta.core.database import get_database_client
from postfacta.core.models import NewIncidentRequest, Note, create_new_incident
from postfacta.core.errors import ErrorResponse, IncidentNotFoundError


logger = logging.getLogger(__name__)
router_v0 = APIRouter(prefix="/v0", tags=["V0"])


incidents_router = APIRouter(prefix="/incidents", tags=["INCIDENTS"])


@incidents_router.get("", status_code=status.HTTP_200_OK)
async def get_all_incidents():
    """Retrieve all incidents."""
    client = get_database_client()
    all_incidents = client.get_all()
    return {"incidents": all_incidents}


@incidents_router.post("/start", status_code=status.HTTP_201_CREATED)
async def create_incident(incident_request: NewIncidentRequest):
    """Create a new incident."""
    client = get_database_client()
    new_incident = create_new_incident(incident_request)
    client.register(new_incident)
    logger.info(f"New incident reported by {new_incident.reporter}: {new_incident.id}")
    return {"message": "Incident created", "id": new_incident.id}


@incidents_router.get("/{incident_id}", status_code=status.HTTP_200_OK)
async def get_incident_by_id(incident_id: str):
    """Retrieve a single incident by its ID."""
    client = get_database_client()
    try:
        incident = client.get_by_id(incident_id)
        logger.info(f"Incident ID {incident_id} retrieved successfully")
        return incident

    except IncidentNotFoundError as infe:
        logger.error(str(infe))
        error_response = ErrorResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            text=str(infe)
        )
        raise error_response.as_http_exception()


@incidents_router.delete("/{incident_id}", status_code=status.HTTP_200_OK)
async def delete_incident_by_id(incident_id: str):
    """Delete a single incident by its ID."""
    client = get_database_client()
    try:
        client.remove_by_id(incident_id)
        logger.info(f"Incident ID {incident_id} deleted successfully")
        return {"message": "Incident deleted", "id": incident_id}

    except IncidentNotFoundError as infe:
        logger.error(str(infe))
        error_response = ErrorResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            text=str(infe)
        )
        raise error_response.as_http_exception()


@incidents_router.get("/{incident_id}/reports", status_code=status.HTTP_200_OK)
async def get_incident_reports(incident_id: str) -> list[Note]:
    """Retrieve the report of a single incident by its ID."""
    client = get_database_client()
    try:
        incident = client.get_by_id(incident_id)
        logger.info(f"Incident ID {incident_id} reports retrieved successfully")
        return incident.get_notes()

    except IncidentNotFoundError as infe:
        error_response = ErrorResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            text=str(infe)
        )
        raise error_response.as_http_exception()


@incidents_router.patch("/{incident_id}/reports", status_code=status.HTTP_204_NO_CONTENT)
async def add_incident_notes(incident_id: str, notes: list[str]) -> None:
    """Retrieve the report of a single incident by its ID."""
    client = get_database_client()
    try:
        incident = client.get_by_id(incident_id)
        for note_message in notes:
            timestamp = dt.datetime.now(dt.timezone.utc).isoformat()
            note = Note(timestamp=timestamp, message=note_message)
            incident.add_note(note)
        client.update_entry(incident)
        logger.info(f"Note added to Incident ID {incident_id} successfully")

    except IncidentNotFoundError as infe:
        logger.error(str(infe))
        error_response = ErrorResponse(
            status_code=status.HTTP_404_NOT_FOUND,
            text=str(infe)
        )
        raise error_response.as_http_exception()


# Include the incidents router in the main v0 router
router_v0.include_router(incidents_router)
