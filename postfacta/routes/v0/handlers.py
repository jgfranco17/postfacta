import logging

from fastapi import APIRouter, status
from pydantic import ValidationError

from postfacta.core.models import NewIncidentRequest, create_new_incident
from postfacta.core.errors import ErrorResponse


logger = logging.getLogger(__name__)
router_v0 = APIRouter(prefix="/v0", tags=["V0"])


incidents_router = APIRouter(prefix="/incidents", tags=["INCIDENTS"])
router_v0.include_router(incidents_router)


@incidents_router.post("/start", status_code=status.HTTP_201_CREATED)
async def create_incident(incident_request: NewIncidentRequest):
    """Create a new incident."""
    try:
        new_incident = create_new_incident(incident_request)
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
