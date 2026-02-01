import logging

from fastapi import APIRouter, status
from pydantic import ValidationError

from postfacta.core.database import get_database_client
from postfacta.core.models import NewIncidentRequest, create_new_incident
from postfacta.core.errors import ErrorResponse


logger = logging.getLogger(__name__)
router_v0 = APIRouter(prefix="/v0", tags=["V0"])


incidents_router = APIRouter(prefix="/incidents", tags=["INCIDENTS"])


@incidents_router.get("", status_code=status.HTTP_200_OK)
async def get_all_incidents():
    """Retrieve all incidents."""
    client = get_database_client()
    try:
        all_incidents = client.get_all()
        return {"incidents": all_incidents}

    except ValidationError as ve:
        logger.error(f"Error creating incident: {ve}")
        error_response = ErrorResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            text="Invalid incident data provided.",
            additional_info=str(ve)
        )
        raise error_response.as_http_exception()


@incidents_router.post("/start", status_code=status.HTTP_201_CREATED)
async def create_incident(incident_request: NewIncidentRequest):
    """Create a new incident."""
    client = get_database_client()
    try:
        new_incident = create_new_incident(incident_request)
        client.register(new_incident)
        logger.info(f"New incident reported by {new_incident.reporter}: {new_incident.id}")
        return {"message": "Incident created", "id": new_incident.id}

    except ValidationError as ve:
        logger.error(f"Error creating incident: {ve}")
        error_response = ErrorResponse(
            status_code=status.HTTP_400_BAD_REQUEST,
            text="Invalid incident data provided.",
            additional_info=str(ve)
        )
        raise error_response.as_http_exception()


router_v0.include_router(incidents_router)
